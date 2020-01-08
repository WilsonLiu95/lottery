# 环境配置
直接执行初始化脚本: source init.sh 或者按下列步骤手动执行

## 安装 go 依赖
- go get -u -v github.com/gin-gonic/gin
- go get -u -v github.com/gorilla/websocket
- go get -u -v github.com/Luxurioust/excelize
- go get -u -v github.com/alecthomas/log4go

## 安装 node 依赖
- cd wwwroot/html
- tnpm i

## 创建日志文件夹
mkdir log

# 启动
go run lottery.go

# 页面预览地址
- 管理页: http://localhost:12345/view/behind
- 展示页: http://localhost:12345/view/index

# 头像转 base64 工具
[avatar_tool](avatar_tool/ReadMe.md)