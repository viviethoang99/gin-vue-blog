package middleware

import (
	"bytes"
	g "gin-blog/internal/global"
	"gin-blog/internal/handle"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"io"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TODO: Optimize API path format
var optMap = map[string]string{
	"Article":      "Article",
	"BlogInfo":     "Blog Info",
	"Category":     "Category",
	"Comment":      "Comment",
	"FriendLink":   "Friend Link",
	"Menu":         "Menu",
	"Message":      "Message",
	"OperationLog": "Operation Log",
	"Resource":     "Resource Permission",
	"Role":         "Role",
	"Tag":          "Tag",
	"User":         "User",
	"Page":         "Page",
	// "Login":        "Login",

	"POST":   "Create or Update",
	"PUT":    "Update",
	"DELETE": "Delete",
}

func GetOptString(key string) string {
	return optMap[key]
}

// Get Response Body content in gin: wrap gin's ResponseWriter so that each response also writes to a buffer
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer // Response body cache
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // Write response data to cache
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s) // Write response data to cache
	return w.ResponseWriter.WriteString(s)
}

// Operation log middleware
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Record file uploads
		// Do not record GET requests (too many) and file upload operations (request body too long)
		if c.Request.Method != "GET" && !strings.Contains(c.Request.RequestURI, "upload") {
			blw := &CustomResponseWriter{
				body:           bytes.NewBufferString(""),
				ResponseWriter: c.Writer,
			}
			c.Writer = blw

			auth, _ := handle.CurrentUserAuth(c)

			body, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			ipAddress := utils.IP.GetIpAddress(c)
			ipSource := utils.IP.GetIpSource(ipAddress)

			moduleName := getOptResource(c.HandlerName())
			operationLog := model.OperationLog{
				OptModule:     moduleName, // TODO: optimize
				OptType:       GetOptString(c.Request.Method),
				OptUrl:        c.Request.RequestURI,
				OptMethod:     c.HandlerName(),
				OptDesc:       GetOptString(c.Request.Method) + " " + moduleName, // TODO: optimize
				RequestParam:  string(body),
				RequestMethod: c.Request.Method,
				UserId:        auth.UserInfoId,
				Nickname:      auth.UserInfo.Nickname,
				IpAddress:     ipAddress,
				IpSource:      ipSource,
			}
			c.Next()
			operationLog.ResponseData = blw.body.String() // Get response body content from cache

			db := c.MustGet(g.CTX_DB).(*gorm.DB)
			if err := db.Create(&operationLog).Error; err != nil {
				slog.Error("Failed to record operation log: ", err)
				handle.ReturnError(c, g.ErrDbOp, err)
				return
			}
		} else {
			c.Next()
		}
	}
}

// "gin-blog/api/v1.(*Resource).Delete-fm" => "Resource"
func getOptResource(handlerName string) string {
	s := strings.Split(handlerName, ".")[1]
	return s[2 : len(s)-1]
}
