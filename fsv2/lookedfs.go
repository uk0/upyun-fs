package fsv2

import (
	"encoding/json"
	"fmt"
	"github.com/billziss-gh/cgofuse/fuse"
	"github.com/upyun/go-sdk/upyun"
	"hash/fnv"
	"strings"
)

const (
	filename = "hello"
	ROOTPATH = "/"
)

type Lookedfs struct {
	fuse.FileSystemBase
}

func (self *Lookedfs) Open(path string, flags int) (errc int, fh uint64) {
	switch path {
	default:
		return -fuse.ENOENT, ^uint64(0)
	}
}

func GetLastIndex(path string) string {
	if path == "/" {
		return path
	}
	var array = strings.Split(path, "/")

	return array[len(array)-1]
}

func (self *Lookedfs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	infoFIle, _ := GetPathInfo(path)

	switch infoFIle.IsDir {
	case true:
		stat.Mode = fuse.S_IFDIR | 0555
		return 0
	case false:
		stat.Mode = fuse.S_IFREG | 0777
		stat.Size = infoFIle.Size
		return 0
	default:
		return -fuse.ENOENT
	}
}
func Str(n uint32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return Str(h.Sum32())

}

func (self *Lookedfs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	fmt.Println(path)
	fmt.Println(buff)
	fmt.Println(ofst)
	fmt.Println(fh)

	return
}

func (self *Lookedfs) Write(path string, buff []byte, ofst int64, fh uint64) (n int) {
	fmt.Println(path)
	fmt.Println(string(buff))
	fmt.Println(ofst)
	fmt.Println(fh)
	return
}

func (self *Lookedfs) Readdir(path string, fill func(name string, stat *fuse.Stat_t, ofst int64) bool, ofst int64, fh uint64) (errc int) {
	fill(".", nil, 0)
	fill("..", nil, 0)
	fill(filename, nil, 0)
	// 判断缓存是否成立

	var Dirs = []UPFSFiles{}
	if TableHasIn(path){
		Dirs = TableSelectsFileArray(path)
	}else{
		var DirsChan = make(chan *upyun.FileInfo)
		go func() {
			_ = FSsysTemp.Client.List(&upyun.GetObjectsConfig{
				Path:        path,
				ObjectsChan: DirsChan,
			})
		}();
		for fileObj := range DirsChan {
			Dirs = append(Dirs, UPFSFiles{
				Path:fileObj.Name,
				Size:fileObj.Size,
				Name:fileObj.Name,
				CTime:fileObj.Time,
				MTime:fileObj.Time,
			})
		}
		// 存储
		value,_:=json.Marshal(&UPFSFiles{
			Path:path,
			Size:0,
			Name:path,
			Files:Dirs,
		})
		// 数据缓存起来
		TableInsert(path,value)

	}


	for _,fileObj := range Dirs {
		fmt.Println("----------------------------------------------")
		fmt.Println(fileObj)
		fmt.Println("----------------------------------------------")
		if fileObj.IsDir {
			fill(fileObj.Name, &fuse.Stat_t{Mode: fuse.S_IFDIR}, 0)
		} else {
			fill(fileObj.Name, &fuse.Stat_t{Mode: fuse.S_IFREG, Size: fileObj.Size}, 0)
		}
	}
	return 0
}

func Run(conf Config) {
	INIT(conf)
	lookedfs := &Lookedfs{}
	host := fuse.NewFileSystemHost(lookedfs)

	var mountPath = []string{conf.Mountpoint}
	host.Mount("", mountPath)
}
