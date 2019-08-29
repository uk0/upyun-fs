package fsv2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/upyun/go-sdk/upyun"
	"path"
	"time"
)

func INIT(conf Config) {
	INITConfig(conf, "/")
}

var vmConf = Config{}

type FSSys struct {
	Client *upyun.UpYun;
	DB     *leveldb.DB;
}
type UPFSFiles struct {
	Path  string
	Name  string
	Size  int64
	IsDir bool
	CTime time.Time
	MTime time.Time
	Files []UPFSFiles
}


var FSsysTemp = FSSys{};

func GetPathINFO(Path string) (restlt *UPFSFiles, err error) {
	FileInfo, err := FSsysTemp.Client.GetInfo(Path)
	restlt.Size = FileInfo.Size
	restlt.Name = FileInfo.Name
	restlt.IsDir = FileInfo.IsDir
	return restlt, err;
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

func GetFileObjList(dirPath string) {
	usage, _ := FSsysTemp.Client.Usage()
	fmt.Println(fmt.Sprintf(" %d MB", usage/1024/1024))
	var DirsChan = make(chan *upyun.FileInfo)
	go func() {
		_ = FSsysTemp.Client.List(&upyun.GetObjectsConfig{
			Path:        dirPath,
			ObjectsChan: DirsChan,
		})
	}();
	for fileInfo:=range DirsChan{
		var RemotePath = path.Join(dirPath,"/",fileInfo.Name)
		if fileInfo.IsDir {
			GetFileObjList(RemotePath)
			fmt.Println("---------------IS DIR----------------" + RemotePath)
		} else {
			Format(&UPFSFiles{
				Path:  RemotePath,
				Name:  fileInfo.Name,
				Size:  fileInfo.Size,
				IsDir: fileInfo.IsDir,
			})
		}
	}
}

// 获得文件夹下的文件列表以及文件夹列表
func INITConfig(conf Config, dirPath string) {
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		//some args
		Bucket:   conf.Upx.Bucket,
		Operator: conf.Upx.Operator,
		Password: conf.Upx.Password,
	})
	FSsysTemp.Client = up;
}

func SaveKVStore(key string, value *UPFSFiles) {
	valByte, _ := json.Marshal(value)
	err := FSsysTemp.DB.Put([]byte(hash(key)), valByte, nil)
	if err != nil {
		fmt.Println("cache inner")
	}
	fmt.Println("[KV] Insert ERR " + fmt.Sprint(err))
}

func GetAllFileList() (FileList []UPFSFiles) {
	iter := FSsysTemp.DB.NewIterator(nil, nil)
	for iter.Next() {
		//key := iter.Key()
		value := iter.Value()
		//fmt.Printf("key: %s | value: %s\n", key, value)
		temp := UPFSFiles{}
		_ = json.Unmarshal(value, &temp)
		FileList = append(FileList, temp)
	}
	return FileList
}

func HasKVStore(key string) (status bool) {
	status, _ = FSsysTemp.DB.Has([]byte(hash(key)), nil)
	return status
}
func GetTreeList() {
}

func GetKVStore(key string) (resultFileInto *UPFSFiles) {
	data, err := FSsysTemp.DB.Get([]byte(hash(key)), nil)
	if err != nil {
		fmt.Println("Get KV ERROR " + fmt.Sprint(err) + "  Key   " + key)
		fmt.Println("Get KV ERROR " + fmt.Sprint(err) + " Hash Key   " + hash(key))

	}
	//fmt.Println("Get KV " + string(data))
	_ = json.Unmarshal(data, &resultFileInto)

	return resultFileInto;
}
