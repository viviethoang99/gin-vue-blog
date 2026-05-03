package upload

import (
	g "gin-blog/internal/global"
	"mime/multipart"
)

// OSS object storage interface
type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

// Select file upload implementation based on configuration
func NewOSS() OSS {
	switch g.GetConfig().Upload.OssType {
	case "local":
		return &Local{}
	case "qiniu":
		return &Qiniu{}
	case "aws":
		return &Aws{}
	default:
		return &Local{}
	}
}
