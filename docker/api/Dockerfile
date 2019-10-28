FROM golang:latest

COPY src/api /go/src/api

WORKDIR /go/src/api/
#WORKDIR /go/src/github.com/uma-co82/Shupple-api/src/api/

RUN go get -u github.com/gin-gonic/gin \
  && go get -u github.com/jinzhu/gorm \
  && go get -u github.com/go-sql-driver/mysql \
  && go get -u gopkg.in/go-playground/validator.v9 \
  && go get -u github.com/aws/aws-sdk-go \
  && go get -u github.com/kelseyhightower/envconfig \
  && go get -u github.com/carlescere/scheduler \
  && go get -u github.com/bamzi/jobrunner \
#  && go get -u github.com/uma-co82/Shupple-api/src/api \
  && apt-get update -qq \
  && apt-get install -y mariadb-client vim


CMD ["sh", "./start_app.sh"]