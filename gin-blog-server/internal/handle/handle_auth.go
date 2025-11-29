package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"gin-blog/internal/utils/jwt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserAuth struct{}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=4,max=20"`
	
}

type LoginVO struct {
	model.UserInfo

	// Like sets: record which articles/comments the user liked
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
	Token          string   `json:"token"`
}

// @Summary Login
// @Description Login
// @Tags UserAuth
// @Param form body LoginReq true "Login"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.LoginVO]
// @Router /login [post]
func (*UserAuth) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	userAuth, err := model.GetUserAuthInfoByName(db, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrUserNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// Check whether password is correct
	if !utils.BcryptCheck(req.Password, userAuth.Password) {
		ReturnError(c, g.ErrPassword, nil)
		return
	}

	// Get IP related info
	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSourceSimpleIdle(ipAddress)

	// browser, os := "unknown", "unknown"
	// if userAgent := utils.IP.GetUserAgent(c); userAgent != nil {
	// 	browser = userAgent.Name + " " + userAgent.Version.String()
	// 	os = userAgent.OS + " " + userAgent.OSVersion.String()
	// }

	userInfo, err := model.GetUserInfoById(db, userAuth.UserInfoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrUserNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	roleIds, err := model.GetRoleIdsByUserId(db, userAuth.ID)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	articleLikeSet, err := rdb.SMembers(rctx, g.ARTICLE_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}
	commentLikeSet, err := rdb.SMembers(rctx, g.COMMENT_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	// Credentials valid, generate token

	// UUID idea: ip + browser info + OS info
	// uuid := utils.MD5(ipAddress + browser + os)
	conf := g.Conf.JWT
	token, err := jwt.GenToken(conf.Secret, conf.Issuer, int(conf.Expire), userAuth.ID, roleIds)
	if err != nil {
		ReturnError(c, g.ErrTokenCreate, err)
		return
	}

	// Update user login info: ip + last login time
	err = model.UpdateUserLoginInfo(db, userAuth.ID, ipAddress, ipSource)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	slog.Info("User login success: " + userAuth.Username)

	session := sessions.Default(c)
	session.Set(g.CTX_USER_AUTH, userAuth.ID)
	session.Save()

	// Remove offline status from Redis
	offlineKey := g.OFFLINE_USER + strconv.Itoa(userAuth.ID)
	rdb.Del(rctx, offlineKey).Result()

	ReturnSuccess(c, LoginVO{
		UserInfo: *userInfo,

		ArticleLikeSet: articleLikeSet,
		CommentLikeSet: commentLikeSet,
		Token:          token,
	})
}


// @Summary Logout
// @Description Logout
// @Tags UserAuth
// @Accept json
// @Produce json
// @Success 0 {object} string
// @Router /logout [post]
func (*UserAuth) Logout(c *gin.Context) {
	c.Set(g.CTX_USER_AUTH, nil)
	
	// Already logged out
	auth, _ := CurrentUserAuth(c)
	if auth == nil {
		ReturnSuccess(c, nil)
		return
	}

	session := sessions.Default(c)
	session.Delete(g.CTX_USER_AUTH)
	session.Save()
	
	// Remove online status from Redis
	rdb := GetRDB(c)
	onlineKey := g.ONLINE_USER + strconv.Itoa(auth.ID)
	rdb.Del(rctx, onlineKey)

	ReturnSuccess(c, nil)
}

// Complete registration flow
// First check whether the username exists to avoid duplicate registration; then store encrypted info in Redis waiting for verification
// Errors: 1) Email already registered 2) Invalid email causing send failure
func (*UserAuth) Register(c *gin.Context) {
	var regreq RegisterReq
	if err := c.ShouldBindJSON(&regreq); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}
	// Normalize username
	regreq.Username = utils.Format(regreq.Username)

	// Check whether username exists to avoid duplicate registration
	auth,err := model.GetUserAuthInfoByName(GetDB(c),regreq.Username)
	if err != nil {
		var flag bool = false
		if errors.Is(err,gorm.ErrRecordNotFound) {
			flag = true
		}
		if !flag{
			ReturnError(c,g.ErrDbOp,err)
			return
		}
	}

	if auth != nil {
		ReturnError(c,g.ErrUserExist,err)
		return
	}
	

	// Verify via email
	info := utils.GenEmailVerificationInfo(regreq.Username,regreq.Password)
	SetMailInfo(GetRDB(c),info,15*time.Minute) // expires in 15 minutes
	EmailData := utils.GetEmailData(regreq.Username,info)
	err = utils.SendEmail(regreq.Username,EmailData)
	if err != nil {
		ReturnError(c,g.ErrSendEmail,err)
		return
	}

	ReturnSuccess(c,nil)
}

// Email verification
// When the user clicks the link in the email, it sends info (encrypted username/password) to this endpoint.
// Verify checks whether info exists in Redis; if present, verification succeeds and registration completes.
// Errors: 1) Missing info in request 2) Info not in Redis (expired) 3) Failed to create user
func (*UserAuth) VerifyCode(c *gin.Context) {
    var code string
    if code = c.Query("info"); code == "" {
        returnErrorPage(c)
        return
    }

	// Verify code exists in Redis
    ifExist, err := GetMailInfo(GetRDB(c), code)
    if err != nil {
        returnErrorPage(c)
        return
    }
    if !ifExist {
        returnErrorPage(c)
        return
    }

    DeleteMailInfo(GetRDB(c), code)

    username, password, err := utils.ParseEmailVerificationInfo(code)
    if err != nil {
        returnErrorPage(c)
        return
    }

	// Register user
      _,_,_,err = model.CreateNewUser(GetDB(c), username, password)
    if err != nil {
        returnErrorPage(c)
        return
    }

	// Registration success: return success page
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
        <!DOCTYPE html>
		<html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Registration Successful</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #5cb85c;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
				<h1>Registration Successful</h1>
				<p>Congratulations, registration succeeded!</p>
            </div>
        </body>
        </html>
    `))
}

func returnErrorPage(c *gin.Context) {
    c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`
        <!DOCTYPE html>
		<html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Registration Failed</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #d9534f;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
				<h1>Registration Failed</h1>
				<p>Please try again.</p>
            </div>
        </body>
        </html>
    `))
}