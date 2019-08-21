package fs

import (
	"bytes"
	"github.com/upyun/go-sdk/upyun"
)

func GetINFO(Path string) (*upyun.FileInfo, error) {
	FileInfo, err := FSsysTemp.Client.GetInfo(Path)
	return FileInfo, err;
}

func GetFile(FilePath string) (info *upyun.FileInfo, err error, data bytes.Buffer) {

	info, err = FSsysTemp.Client.Get(&upyun.GetObjectConfig{
		Path: FilePath,
		Writer:&data,
	})

	return info, err,data
}
