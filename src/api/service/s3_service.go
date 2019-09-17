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
	"os"
)

type S3Service struct {
}

func UploadToS3(image string) {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("", "", ""),
		Region:      aws.String("ap-northeast-1"),
	}))

	uploader := s3manager.NewUploader(sess)

	data, _ := base64.StdEncoding.DecodeString(image)
	wb := new(bytes.Buffer)
	wb.Write(data)

	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(""),
		Key:    aws.String("sample.png"),
		Body:   wb,
	})

	if err != nil {
		fmt.Println(res)
		if err, ok := err.(awserr.Error); ok && err.Code() == request.CanceledErrorCode {
			fmt.Fprint(os.Stderr, "upload canceled due to timeout, %v\n", err)
		} else {
			fmt.Fprint(os.Stderr, "failed to upload object %v\n", "bucket")
		}
	}
}
