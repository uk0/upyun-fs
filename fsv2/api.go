package fsv2

import (
	"bytes"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/upyun/go-sdk/upyun"
	"time"
)

func INIT(conf Config) {
	INITConfig(conf)
}

var vmConf = Config{}

type FSSys struct {
	Client  *upyun.UpYun;
	DBInfo  *leveldb.DB;
	DBTable *leveldb.DB;
}
type UPFSFiles struct {
	InfoID string
	Path   string
	Name   string
	Size   int64
	IsDir  bool
	CTime  time.Time
	MTime  time.Time
	Files  []UPFSFiles
}

var FSsysTemp = FSSys{};

func GetPathInfo(Path string) (restlt *UPFSFiles, err error) {
	FileInfo, err := FSsysTemp.Client.GetInfo(Path)
	return &UPFSFiles{
		Size:  FileInfo.Size,
		IsDir: FileInfo.IsDir,
		CTime: FileInfo.Time,
		MTime: FileInfo.Time,
	}, err;
}

func GetFile(FilePath string) (info *upyun.FileInfo, err error, data bytes.Buffer) {
	info, err = FSsysTemp.Client.Get(&upyun.GetObjectConfig{
		Path:   FilePath,
		Writer: &data,
	})

	return info, err, data
}

// 方法区域，
func GetContext(path string) (b bytes.Buffer, FInfo *upyun.FileInfo) {
	FInfo, _ = FSsysTemp.Client.Get(&upyun.GetObjectConfig{
		Path:   path,
		Writer: &b,
	})
	return b, FInfo
}

func GetPathSize(path string) uint64 {
	FInfo, _ := FSsysTemp.Client.GetInfo(path)
	return uint64(FInfo.Size)
}

// 获得文件夹下的文件列表以及文件夹列表
func INITConfig(conf Config) {
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		//some args
		Bucket:   conf.Upx.Bucket,
		Operator: conf.Upx.Operator,
		Password: conf.Upx.Password,
	})
	FSsysTemp.Client = up;
}

// level db

func TableInsert(key string, value []byte) {
	_ = FSsysTemp.DBTable.Put([]byte(hash(key)), value, nil)
}

func TableHasIn(key string) bool {
	status, _ := FSsysTemp.DBTable.Has([]byte(hash(key)), nil)
	return status;
}

func TableSelectsFileArray(key string) (upfs []UPFSFiles) {
	tempUPFSFiles  :=UPFSFiles{}
	data, _ := FSsysTemp.DBTable.Get([]byte(hash(key)), nil)
	_ = json.Unmarshal(data, &tempUPFSFiles)
	return tempUPFSFiles.Files
}
