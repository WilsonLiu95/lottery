#!/bin/sh
go get -u -v github.com/gin-gonic/gin
go get -u -v github.com/gorilla/websocket
go get -u -v github.com/Luxurioust/excelize
go get -u -v github.com/alecthomas/log4go

rm -rf log
mkdir log 

cd wwwroot/html/
tnpm i