package middleware

import (
	"bytes"
	"encoding/json"
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

			// 未登录时(JWTAuth 对资源表中不存在的接口会跳过鉴权)拿不到用户,
			// 此处不能直接取字段, 否则会 nil 解引用
			auth, _ := handle.CurrentUserAuth(c)
			var userId int
			var nickname string
			if auth != nil {
				userId = auth.UserInfoId
				if auth.UserInfo != nil {
					nickname = auth.UserInfo.Nickname
				}
			}

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
				OptDesc:       GetOptString(c.Request.Method) + moduleName, // TODO: 优化
				RequestParam:  maskSensitive(string(body)),
				RequestMethod: c.Request.Method,
				UserId:        userId,
				Nickname:      nickname,
				IpAddress:     ipAddress,
				IpSource:      ipSource,
			}
			c.Next()
			operationLog.ResponseData = maskSensitive(blw.body.String()) // 从缓存中获取响应体内容

			db := c.MustGet(g.CTX_DB).(*gorm.DB)
			if err := db.Create(&operationLog).Error; err != nil {
				slog.Error("操作日志记录失败", "err", err)
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

const sensitiveMask = "******"

// 不能原文落库的字段名, 按小写子串匹配:
// 改密码接口 (PUT /user/current/password) 就挂在这个中间件下面,
// 旧密码与新密码原来会明文写进 operation_log.request_param(longtext) 永久留存。
var sensitiveKeys = []string{"password", "token", "secret", "access_key", "accesskey", "captcha"}

func isSensitiveKey(key string) bool {
	k := strings.ToLower(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(k, s) {
			return true
		}
	}
	return false
}

/*
按字段名给 JSON 文本脱敏。三种情况:
  - 没有敏感字段: 原样返回, 不重新序列化 —— 免得字段顺序变了, 日志详情看着和请求对不上
  - 合法 JSON 且命中敏感字段: 只把那几个值换成掩码, 其余保留
  - 不是合法 JSON 却带敏感字样(如表单编码的请求体): 整体丢弃。
    宁可少记一条日志, 也不能把明文密码留在库里
*/
func maskSensitive(raw string) string {
	if raw == "" || !hasSensitiveWord(raw) {
		return raw
	}

	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return sensitiveMask
	}
	masked, changed := maskValue(v)
	if !changed {
		return raw
	}
	b, err := json.Marshal(masked)
	if err != nil {
		return sensitiveMask
	}
	return string(b)
}

// 先做一次廉价筛查, 绝大多数请求体不含敏感字样, 不必解析 JSON
func hasSensitiveWord(raw string) bool {
	lower := strings.ToLower(raw)
	for _, s := range sensitiveKeys {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}

func maskValue(v any) (any, bool) {
	switch t := v.(type) {
	case map[string]any:
		changed := false
		for k, val := range t {
			if isSensitiveKey(k) {
				t[k] = sensitiveMask
				changed = true
				continue
			}
			if nv, c := maskValue(val); c {
				t[k] = nv
				changed = true
			}
		}
		return t, changed
	case []any:
		changed := false
		for i, val := range t {
			if nv, c := maskValue(val); c {
				t[i] = nv
				changed = true
			}
		}
		return t, changed
	}
	return v, false
}
