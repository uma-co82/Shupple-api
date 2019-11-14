# Shupple-api

<img src="https://images.unsplash.com/photo-1541278107931-e006523892df?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2851&q=80" align="right"
     title="Size Limit logo by Anton Lovchikov" width="" height="178">


## What Run

* Go 1.12.3
* gin
* gorm
* Docker
* mysql

## Set Up

### 1. Install Docker
`https://docs.docker.com/`

### 2. Pull Repository
`git clone git@github.com:uma-co82/Shupple-api.git`

### 3. Docker-compose up
`docker-compose up -d`

## AWS EC2 Settings
[reference](http://kitakitabauer.hatenablog.com/entry/2017/10/17/215316)

### 1. Install go

```
sudo yum update -y
```
```
mkdir ~/tmp; cd ~/tmp
```
```
wget https://storage.googleapis.com/golang/go1.12.3.linux-amd64.tar.gz
```
```
tar zxvf go1.12.3.linux-amd64.tar.gz
```
```
sudo mv go /usr/local/
```
```
sudo ln -s /usr/local/go/bin/go /usr/bin/go
```
```
sudo ln -s /usr/local/go/bin/go /usr/local/bin/go
```
```
sudo ln -s /usr/local/go/bin/godoc /usr/local/bin/godoc
```
```
sudo ln -s /usr/local/go/bin/gofmt /usr/local/bin/gofmt
```

### 2. SetUp GoEnv

.bash_profile
```
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
```
export S3AK="s3accesskey"
export S3SK="s3secretkey"
```

### 3. go get

```
sudo yum install -y git
```
```
go get -u github.com/uma-co82/Shupple-api/src/api
```
```
nohup go run main.go  >> nohup.out 2>&1 < /dev/null &
```

## License

#### Members

 ##### [Yuta Isozaki](https://github.com/uma-co82)
