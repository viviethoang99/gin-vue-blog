package middleware

import (
	"errors"
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/handle"
	"gin-blog/internal/model"
	"gin-blog/internal/utils/jwt"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// JWT-based authorization
// If a session exists, get user info from session
// If no session, read token from Authorization, parse it to get user info, and set session
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// FIXME: Session mix-up between front/back; cannot cache user in gin context for now
		// auth, _ := handle.CurrentUserAuth(c)
		// if auth != nil {
		// 	slog.Debug("[middleware-JWTAuth] user auth exist, skip jwt auth")
		// 	c.Next()
		// 	return
		// }

		slog.Debug("[middleware-JWTAuth] user auth not exist, do jwt auth")

		db := c.MustGet(g.CTX_DB).(*gorm.DB)

		// Only system-managed resources require verification; others are skipped
		url, method := c.FullPath()[4:], c.Request.Method
		resource, err := model.GetResource(db, url, method)
		if err != nil {
			// Resource not found; skip further checks
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Debug("[middleware-JWTAuth] resource not exist, skip jwt auth")
				c.Set("skip_check", true)
				c.Next()
				c.Set("skip_check", false)
				return
			}
			handle.ReturnError(c, g.ErrDbOp, err)
			return
		}

		// Anonymous resource; skip further checks
		if resource.Anonymous {
			slog.Debug(fmt.Sprintf("[middleware-JWTAuth] resource: %s %s is anonymous, skip jwt auth!", url, method))
			c.Set("skip_check", true)
			c.Next()
			c.Set("skip_check", false)
			return
		}

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			handle.ReturnError(c, g.ErrTokenNotExist, nil)
			return
		}

		// Token format: `Bearer [tokenString]`
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handle.ReturnError(c, g.ErrTokenType, nil)
			return
		}

		claims, err := jwt.ParseToken(g.Conf.JWT.Secret, parts[1])
		if err != nil {
			handle.ReturnError(c, g.ErrTokenWrong, err)
			return
		}

		// Check whether token expired
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			handle.ReturnError(c, g.ErrTokenRuntime, nil)
			return
		}

		user, err := model.GetUserAuthInfoById(db, claims.UserId)
		if err != nil {
			handle.ReturnError(c, g.ErrUserNotExist, err)
			return
		}

		// session
		session := sessions.Default(c)
		session.Set(g.CTX_USER_AUTH, claims.UserId)
		session.Save()

		// gin context
		c.Set(g.CTX_USER_AUTH, user)
	}
}

// Resource access permission check
func PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool("skip_check") {
			c.Next()
			return
		}

		db := c.MustGet(g.CTX_DB).(*gorm.DB)
		auth, err := handle.CurrentUserAuth(c)
		if err != nil {
			handle.ReturnError(c, g.ErrUserNotExist, err)
			return
		}

		if auth.IsSuper {
			slog.Debug("[middleware-PermissionCheck]: super admin no need to check, pass!")
			c.Next()
			return
		}

		url := c.FullPath()[4:]
		method := c.Request.Method

		slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] %v, %v, %v\n", auth.Username, url, method))
		for _, role := range auth.Roles {
			slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] %v\n", role.Name))
			pass, err := model.CheckRoleAuth(db, role.ID, url, method)
			if err != nil {
				handle.ReturnError(c, g.ErrDbOp, err)
				return
			}
			if !pass {
				handle.ReturnError(c, g.ErrPermission, nil)
				return
			}
		}

		slog.Debug("[middleware-PermissionCheck]: pass")
		c.Next()
	}
}
