package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

type Option struct {
	FilerMountRootPath string
	Collection         string
	Replication        string
	TtlSec             int32
	ChunkSizeLimit     int64
	DataCenter         string
	DirListingLimit    int
	EntryCacheTtl      time.Duration
	Umask              os.FileMode

	MountUid   uint32
	MountGid   uint32
	MountMode  os.FileMode
	MountCtime time.Time
	MountMtime time.Time
}

type FuseAttributes struct {
	FileSize      uint64
	Mtime         int64
	FileMode      uint32
	Uid           uint32
	Gid           uint32
	Crtime        int64
	Mime          string
	Replication   string
	Collection    string
	TtlSec        int32
	UserName      string
	GroupName     []string
	SymlinkTarget string
}

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	Path       string
	option     Option
	ufs        *UFS
	attributes FuseAttributes
}

func (dir *Dir) setRootDirAttributes(attr *fuse.Attr) {
	attr.Uid = dir.option.MountUid
	attr.Gid = dir.option.MountGid
	attr.Mode = dir.option.MountMode
	attr.Crtime = dir.option.MountCtime
	attr.Ctime = dir.option.MountCtime
	attr.Mtime = dir.option.MountMtime
	attr.Atime = dir.option.MountMtime
}

// 文件夹的子参数等。
var _ fs.Node = (*Dir)(nil)
var _ = fs.Node(&Dir{})
var _ = fs.NodeRequestLookuper(&Dir{})
var _ = fs.HandleReadDirAller(&Dir{})

func (d *Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	fmt.Println("3 Node Attr  === " + d.Path)
	// 如果和主目录一致
	if d.Path == ROOTPATH {
		attr.Mode = os.ModeDir
		attr.Valid = time.Second
		return nil
	}

	info ,err:= GetINFO(d.Path)

	if err!=nil{
		return nil
	}
	// 断点
	data, _ := json.Marshal(info)
	fmt.Println(fmt.Sprintf("3 Node Attr INFO  === %s ", string(data)))
	if info.IsDir {
		attr.Atime = time.Time{}
		attr.Ctime = time.Time{}
		attr.Size = uint64(info.Size)
		attr.Mode = os.ModeDir
		return nil
	} else {
		attr.Atime = time.Time{}
		attr.Ctime = time.Time{}
		attr.Mode = os.FileMode(os.ModePerm) //中间没有存储 只能0777了
		attr.Size = uint64(info.Size)
		return nil
	}

	return nil
}

func (dir *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	fmt.Println(fmt.Sprintf("-------------------- Mkdir --------------------"))
	return nil, nil
}


func (dir *Dir) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (node fs.Node, err error) {

	if dir.Path == "" {
		dir.Path = ROOTPATH
		// 默认参数
		fmt.Println(fmt.Sprintf("1 NodeStringLookuper  Path %s", dir.Path))
	}

	fmt.Println(fmt.Sprintf("1-2 layer NodeStringLookuper  Path %s", dir.Path))

	entry,err := GetINFO(dir.Path)

	if err!=nil{
		return nil, fuse.EIO
	}

	if entry != nil {
		jsonData, _ := json.Marshal(entry)
		fmt.Println(fmt.Sprintf("data Json %s", string(jsonData)))
		if entry.IsDir {
			node = &Dir{Path: path.Join(dir.Path, req.Name), ufs: dir.ufs, option: Option{}}
		} else {
			node = dir.newFile(req.Name, &Entry{
				Name:        entry.Name,
				IsDirectory: entry.IsDir,
				Attributes: &FuseAttributes{
					FileSize: uint64(entry.Size),
					FileMode: uint32(os.ModePerm),
					UserName: "test",
				},
			})
		}

		resp.EntryValid = time.Duration(0)
		resp.Attr.Mtime = time.Unix(0, 0)
		resp.Attr.Ctime = time.Unix(0, 0)
		resp.Attr.Mode = os.FileMode(os.ModePerm)
		resp.Attr.Gid = 0
		resp.Attr.Uid = 0

		return node, nil
	}
	return nil, fuse.ENOENT
}

func (dir *Dir) newFile(name string, entry *Entry) *File {
	return &File{
		Name: name,
		dir:  dir,
		ufs:  dir.ufs,
	}
}

var _ = fs.HandleReadDirAller(&Dir{})

func (d *Dir) ReadDirAll(ctx context.Context) (ret []fuse.Dirent, err error) {
	//TODO 真的懒的写了。。。
	fmt.Println("2 HandleReadDirAller ")

	fInfoChan := GetInDirListFiles(vmConf, d.Path)

	for entry := range fInfoChan {
		if entry.IsDir {
			fmt.Println(fmt.Sprintf("fInfoChan IsDir Info %s",entry.Name))
			dirent := fuse.Dirent{Name: entry.Name, Type: fuse.DT_Dir}
			ret = append(ret, dirent)
		} else {
			fmt.Println(fmt.Sprintf("fInfoChan NotDir Info %s",entry.Name))
			dirent := fuse.Dirent{Name: entry.Name, Type: fuse.DT_File}
			ret = append(ret, dirent)
		}
	}
	return ret, nil
}
