package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"fmt"
	"github.com/upyun/go-sdk/upyun"
	"upyun-fs/config"
)

var vmConf = config.Config{}

type FSSys struct {
	Client *upyun.UpYun;
}

const ROOTPATH = "/"

var FSsysTemp = FSSys{};

func Run(conf config.Config) error {
	c, err := fuse.Mount(
		conf.Mountpoint,
		fuse.FSName("upyun-fs"),
		fuse.Subtype("upyun-fs"),
		fuse.LocalVolume(),
		fuse.VolumeName("Test upyun  filesystem"),
	)

	vmConf = conf

	if err != nil {
		return err
	}
	defer c.Close()

	if p := c.Protocol(); !p.HasInvalidate() {
		return fmt.Errorf("kernel FUSE support is too old to have invalidations: version %v", p)
	}

	srv := fs.New(c, nil)
	filesys := &UFS{
		testFile: &File{
			fuse: srv,
		},
	}
	filesys.testFile.tick()
	// This goroutine never exits. That's fine for this example.
	go filesys.testFile.update()
	if err := srv.Serve(filesys); err != nil {
		return err
	}

	// Check if the mount process has an error to report.
	<-c.Ready
	if err := c.MountError; err != nil {
		return err
	}
	return nil
}

var _ fs.FS = (*UFS)(nil)

func (f *UFS) Root() (fs.Node, error) {
	return &Dir{ufs: f, Path: f.option.FilerMountRootPath}, nil
}

type UFS struct {
	testFile *File
	option   Option
}

type Entry struct {
	Name        string
	IsDirectory bool
	Attributes  *FuseAttributes
}


var _ = fs.Node(&Dir{})

// 方法区域，
func GetDirDataSize(path string) uint64 {
	FInfo, _ := FSsysTemp.Client.GetInfo(path)
	return uint64(FInfo.Size)
}

// 获得文件夹下的文件列表以及文件夹列表
func GetInDirListFiles(conf config.Config,dirPath string) []fuse.Dirent {
	var dirDirs = []fuse.Dirent{}
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		//some args
		Bucket:   conf.Upx.Bucket,
		Operator: conf.Upx.Operator,
		Password: conf.Upx.Password,
	})
	FSsysTemp.Client = up;
	usage, _ := up.Usage()
	fmt.Println(fmt.Sprintf(" %d MB", usage/1024/1024))
	fInfoChan := make(chan *upyun.FileInfo, 50)
	go func() {
		_ = up.List(&upyun.GetObjectsConfig{
			Path:        dirPath,
			ObjectsChan: fInfoChan,
		})
	}();
	for fileList := range fInfoChan {
		dirDirs = append(dirDirs, fuse.Dirent{Name: fileList.Name, Type: GetFileType(fileList)})
	}
	return dirDirs;
}

func GetFileType(info *upyun.FileInfo) fuse.DirentType {
	switch info.IsDir {
	case true:
		return fuse.DT_Dir;
	case false:
		return fuse.DT_File;
	}
	return 0;
}
