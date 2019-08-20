package main_test

import (
	"fmt"
	"github.com/upyun/go-sdk/upyun"
	"testing"
)

const VERSION = "v0.2.4"
func GetUPClient() *upyun.UpYun {
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		//some args
		Bucket:   "1",
		Operator: "1",
		Password: "1",
		UserAgent: fmt.Sprintf("upx/%s", VERSION),
	})
	fInfoChan := make(chan *upyun.FileInfo, 50)
	_ = up.List(&upyun.GetObjectsConfig{
		Path:        "/demos",
		ObjectsChan: fInfoChan,
	})
	for fileList := range fInfoChan {
		fmt.Print(fileList)
	}
	return up;
}

func Test_getfileList(T *testing.T) {
	GetUPClient()
}
