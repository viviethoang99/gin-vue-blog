package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/utils/upload"

	"github.com/gin-gonic/gin"
)

type Upload struct{}

// @Summary Upload file
// @Description Upload file
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 0 {object} Response[string]
// @Router /upload/file [post]
func (*Upload) UploadFile(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ReturnError(c, g.ErrFileReceive, err)
		return
	}

	oss := upload.NewOSS()
	filePath, _, err := oss.UploadFile(fileHeader)
	if err != nil {
		ReturnError(c, g.ErrFileUpload, err)
		return
	}

	ReturnSuccess(c, filePath)
}
