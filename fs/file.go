package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"syscall"
	"time"
)

// 对于文件操作读取所有内容

var _ fs.HandleReadAller = (*File)(nil)
func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	fmt.Println("ReadAll File")
	return nil, nil
}

type File struct {
	fuse     *fs.Server
	content  atomic.Value
	count    uint64
	fileSize uint64
	Name     string
	dir      *Dir
	ufs      *UFS
	isOpen   bool
}

var _ fs.Node = (*File)(nil)
func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0444
	a.Size = f.fileSize
	return nil
}

// 文件详情接口实现

var _ = fs.FSStatfser(&File{})

func (f *File) Statfs(ctx context.Context, req *fuse.StatfsRequest, resp *fuse.StatfsResponse) error {

	return nil
}

var _ fs.NodeOpener = (*File)(nil)

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	if !req.Flags.IsReadOnly() {
		return nil, fuse.Errno(syscall.EACCES)
	}
	resp.Flags |= fuse.OpenKeepCache
	return f, nil
}

var _ fs.Handle = (*File)(nil)

var _ fs.HandleReader = (*File)(nil)

func (f *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	t := f.content.Load().(string)
	fuseutil.HandleRead(req, resp, []byte(t))
	fmt.Println("4 FileRead--------------------------------------")
	return nil
}

func (f *File) tick() {
	// Intentionally a variable-length format, to demonstrate size changes.
	f.count++
	s := fmt.Sprintf("%d\t%s\n", f.count, time.Now())
	f.content.Store(s)

	// For simplicity, this example tries to send invalidate
	// notifications even when the kernel does not hold a reference to
	// the node, so be extra sure to ignore ErrNotCached.
	if err := f.fuse.InvalidateNodeData(f); err != nil && err != fuse.ErrNotCached {
		log.Printf("invalidate error: %v", err)
	}
}

func (f *File) update() {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for range tick.C {
		f.tick()
	}
}
