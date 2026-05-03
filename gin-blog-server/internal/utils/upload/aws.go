package upload

import (
	"errors"
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/utils"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)
import "context"
import "github.com/aws/aws-sdk-go-v2/service/s3"

type Aws struct{}

func (*Aws) UploadFile(file *multipart.FileHeader) (filePath, fileName string, err error) {
	cfg := awsConfig()
	client := s3.NewFromConfig(cfg)

	ext := path.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	name = utils.MD5(name)
	filename := name + "_" + time.Now().Format("20060102150405") + ext

	// Open the file to upload
	f, err := file.Open()
	if err != nil {
		return "", "", errors.New("function file.Open() Filed, err:" + err.Error())
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:       aws.String(g.GetConfig().Aws.Bucket),
		Key:          aws.String(filename),
		Body:         f,
		ContentType:  aws.String("image/jpeg"), // nên detect tự động
		CacheControl: aws.String("public, max-age=31536000"),
		ACL:          types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		return "", "", errors.New("function client.PutObject() Filed, err:" + err.Error())
	}

	if g.GetConfig().Aws.ImgPath != "" {
		return g.GetConfig().Aws.ImgPath + "/" + filename, filename, nil
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", g.GetConfig().Aws.Bucket, filename), filename, nil
}

func (*Aws) DeleteFile(key string) error {
	return nil
}

func awsConfig() aws.Config {
	return aws.Config{
		Region: g.GetConfig().Aws.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			g.GetConfig().Aws.AccessKey,
			g.GetConfig().Aws.SecretKey,
			"",
		)),
	}
}
