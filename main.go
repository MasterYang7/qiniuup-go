package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

// go run main.go --ak --sk --bucket --localpath --destpath
type Options struct {
	Accesskey string `long:"ak" description:"your access key" required:""`
	Secretkey string `long:"sk" description:"your secret key" default:""`
	Bucket    string `long:"bucket" description:"your bucket name" required:""`
	LocalPath string `long:"localpath" description:"源文件路径" required:""`
	DestPath  string `long:"destpath" description:"目标文件路径" required:""`
}

func main() {
	var opts Options
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}
	accessKey := opts.Accesskey
	secretKey := opts.Secretkey
	mac := credentials.NewCredentials(accessKey, secretKey)

	bucket := opts.Bucket
	localFile := opts.LocalPath
	keypath := opts.DestPath
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	// 判断路径是否为目录
	fileInfo, err := os.Stat(localFile)
	if err != nil {
		log.Fatalf("无法获取文件信息: %v", err)
	}
	if !fileInfo.IsDir() {
		err = uploadManager.UploadFile(context.Background(), localFile, &uploader.ObjectOptions{
			BucketName: bucket,
			ObjectName: &keypath,
		}, nil)
	} else {
		err = uploadManager.UploadDirectory(context.Background(), localFile, &uploader.DirectoryOptions{
			BucketName: bucket,
			UpdateObjectName: func(key string) string {
				return keypath + key
			},
			ObjectConcurrency: 16, // 对象上传并发度
		})
	}

	if err != nil {
		fmt.Println(err)
		return
	}
}
