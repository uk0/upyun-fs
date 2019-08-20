package fs

import (
	"github.com/upyun/go-sdk/upyun"
)

func GetINFO(dirPath string) (*upyun.FileInfo) {
	FileInfo, _ := FSsysTemp.Client.GetInfo(dirPath)
	return FileInfo;
}
