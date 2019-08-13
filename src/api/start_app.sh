# !/bin/bash

# MySQLサーバーが起動するまでmain.goを実行せずにループで待機する
echo 'waiting for mysqld to be connectable...'
sleep 4

echo "app is starting...!"
exec go run main.go