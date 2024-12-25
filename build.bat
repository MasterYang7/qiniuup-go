SET CGO_ENABLED=0
SET GOOS=linux
SET GOPROXY=https://goproxy.cn,direct
SET GOARCH=amd64
go build -o qiniuup-go  main.go