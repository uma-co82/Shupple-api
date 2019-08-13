# !/bin/bash

echo 'waiting for mysqld to be connectable...'
# MySQLサーバーが起動するまで待機する
# TODO: mysqlサーバーの起動を検知するように書き換えるuntilとかで
sleep 4

echo "app is starting...!"
exec go run main.go