# !/bin/bash

echo 'waiting for mysqld to be connectable...'
# MySQLサーバーが起動するまで待機する
sleep 4

echo "app is starting...!"
exec go run main.go