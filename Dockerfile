FROM golang:latest

# ディレクトリ構成は今後変更の可能性あり
COPY src/ /go/src/

WORKDIR /go/src/

RUN go get -u github.com/gin-gonic/gin 

CMD ["go", "run", "main.go"]