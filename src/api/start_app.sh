# !/bin/bash

# MySQLサーバーが起動するまで待機する
# TODO: mysqlサーバーの起動を検知するように書き換えるuntilとかで
until mysqladmin ping -h mysql -P 3306 --silent; do
  echo 'waiting for mysqld to be connectable...'
  sleep 2
done

echo "app is starting...!"
exec go run main.go