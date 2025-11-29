package upload

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/utils"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

// Local file upload
type Local struct{}

// Upload file to local storage
func (*Local) UploadFile(file *multipart.FileHeader) (filePath, fileName string, err error) {
	ext := path.Ext(file.Filename)                                     // Read file extension
	name := strings.TrimSuffix(file.Filename, ext)                     // Read file name
	name = utils.MD5(name)                                             // Hash file name
	filename := name + "_" + time.Now().Format("20060102150405") + ext // Compose new file name

	conf := g.Conf.Upload
	mkdirErr := os.MkdirAll(conf.StorePath, os.ModePerm) // Try to create storage path
	if mkdirErr != nil {
		slog.Error("function os.MkdirAll() Filed", slog.Any("err", mkdirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}

	storePath := conf.StorePath + "/" + filename // File storage path
	filepath := conf.Path + "/" + filename       // File access path

	f, openError := file.Open() // Read file
	if openError != nil {
		slog.Error("function file.Open() Filed", slog.String("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // Defer close after file open

	out, createErr := os.Create(storePath)
	if createErr != nil {
		slog.Error("function os.Create() Filed", slog.String("err", createErr.Error()))
		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close() // Defer close after file create

	_, copyErr := io.Copy(out, f) // Copy file
	if copyErr != nil {
		slog.Error("function io.Copy() Filed", slog.String("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return filepath, filename, nil
}

// Delete file from local storage
func (*Local) DeleteFile(key string) error {
	p := g.GetConfig().Upload.StorePath + "/" + key
	if strings.Contains(p, g.GetConfig().Upload.StorePath) {
		if err := os.Remove(p); err != nil {
			return errors.New("Failed to delete local file, err:" + err.Error())
		}
	}
	return nil
}
