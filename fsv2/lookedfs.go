package fsv2

import (
	"fmt"
	"github.com/billziss-gh/cgofuse/fuse"
	"hash/fnv"
	"strings"
)

const (
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
	var fileNmae = GetLastIndex(path)
	infoFIle := GetKVStore(fileNmae)
	fmt.Println("Path = " + fmt.Sprint(HasKVStore(path)))

	fmt.Println(" GET KV " + fmt.Sprint(infoFIle))
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

func (self *Lookedfs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	fill(".", nil, 0)
	fill("..", nil, 0)
	for _, fileObj := range GetAllFileList() {
		if fileObj.IsDir {
			fill(fileObj.Name, &fuse.Stat_t{Mode: fuse.S_IFDIR}, 0)
		} else {
			fill(fileObj.Name, &fuse.Stat_t{Mode: fuse.S_IFREG, Size: fileObj.Size}, 0)
		}

	}
	return 0
}

//管他呢先存起来
func Format(files *UPFSFiles) {
	SaveKVStore(files.Path, files)
}

func Run(conf Config) {
	INIT(conf)
	lookedfs := &Lookedfs{}
	host := fuse.NewFileSystemHost(lookedfs)

	GetFileObjList(ROOTPATH) // 获得第一批文件 以及文件夹



	var mountPath = []string{conf.Mountpoint}
	host.Mount("", mountPath)
}
