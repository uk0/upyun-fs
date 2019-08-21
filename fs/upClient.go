package fs

import (
	"github.com/upyun/go-sdk/upyun"
)

func GetINFO(Path string) (*upyun.FileInfo) {
	FileInfo, _ := FSsysTemp.Client.GetInfo(Path)
	return FileInfo;
}
