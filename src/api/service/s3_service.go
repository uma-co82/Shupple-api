package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kelseyhightower/envconfig"
)

type (
	S3Service struct{}
	// NOTE: 環境変数にS3AKとS3SKを設定
	Env struct {
		S3AK string
		S3SK string
	}
)

func (s S3Service) UploadToS3(image string, userName string) error {
	// 環境変数からS3Credential周りの設定を取得
	var env Env
	_ = envconfig.Process("", &env)
	fmt.Println(env)

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(env.S3AK, env.S3SK, ""),
		Region:      aws.String("ap-northeast-1"),
	}))

	uploader := s3manager.NewUploader(sess)

	data, _ := base64.StdEncoding.DecodeString(image)
	wb := new(bytes.Buffer)
	wb.Write(data)

	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("isozaki-images"),
		// TODO: 直す
		Key:  aws.String("sample.png"),
		Body: wb,
	})

	if err != nil {
		fmt.Println(res)
		if err, ok := err.(awserr.Error); ok && err.Code() == request.CanceledErrorCode {
			return RaiseError(400, "Upload TimuOut", nil)
		} else {
			return RaiseError(400, "Upload Failed", nil)
		}
	}

	return nil
}
