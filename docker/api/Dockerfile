FROM golang:latest

# ディレクトリ構成は今後変更の可能性あり
COPY src/api /go/src/api

WORKDIR /go/src/api/

RUN go get -u github.com/gin-gonic/gin \
  && go get github.com/jinzhu/gorm

CMD ["go", "run", "main.go"]