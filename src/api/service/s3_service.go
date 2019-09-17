package service

import (
	"bytes"
	"encoding/base64"
)

type S3Service struct {
}

func UploadImage(image string) {
	data, _ := base64.StdEncoding.DecodeString(image)
	wb := new(bytes.Buffer)
	wb.Write(data)
}
