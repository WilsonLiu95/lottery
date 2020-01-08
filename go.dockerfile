FROM golang:1.13

# MAINTAINER wilsonsliu "wilsonsliu@tencent.com"

RUN go get -u -v github.com/gin-gonic/gin \
    && go get -u -v github.com/gorilla/websocket \
    && go get -u -v github.com/Luxurioust/excelize \
    && go get -u -v github.com/alecthomas/log4go

# WORKDIR $GOPATH/lottery
# COPY . .
# EXPOSE 12345
# CMD [ "go", "run", "lottery.go" ]