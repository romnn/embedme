package fs

// import (
// 	"io"
// 	realfs "io/fs"
// 	"os"
// 	"time"
// )
//
// type FileSystem interface {
// 	Open(path string) (File, error)
// 	Stat(path string) (os.FileInfo, error)
// 	DirFS(path string) realfs.FS
// 	// DirFS(path string) FileSystem
// }
//
// type File interface {
// 	io.Closer
// 	io.Reader
// 	io.ReaderAt
// 	io.Seeker
// 	Stat() (os.FileInfo, error)
// }
//
// // type FileSystem interface {
// // 	Open(path string) (os.File, error)
// // 	ReadFile(path string) ([]byte, error)
// // }
//
// // Name() string       // base name of the file
// // 	Size() int64        // length in bytes for regular files; system-dependent for others
// // 	Mode() FileMode     // file mode bits
// // 	ModTime() time.Time // modification time
// // 	IsDir() bool        // abbreviation for Mode().IsDir()
// // 	Sys() any           // underlying data source (can return nil)
//
// // type RealFileSystem struct {
// // 	FileSystem
// // }
//
// type RealFileSystem struct {
// 	FileSystem
// }
//
// func (*RealFileSystem) Open(path string) (File, error) {
// 	return os.Open(path)
// }
//
// func (*RealFileSystem) Stat(path string) (os.FileInfo, error) {
// 	return os.Stat(path)
// }
//
// func (*RealFileSystem) DirFS(path string) realfs.FS {
// 	return os.DirFS(path)
// }
//
// type MockFileSystem struct {
// 	FileSystem
// 	Files map[string]string
// }
//
// type MockFile struct {
// 	File
// 	Content []byte
// }
//
// type MockFileStat struct {
// 	name    string
// 	size    int64
// 	mode    os.FileMode
// 	modTime time.Time
// 	sys     syscall.Stat_t
// }
//
// func (fs *MockFileStat) Size() int64        { return fs.size }
// func (fs *MockFileStat) Mode() os.FileMode  { return fs.mode }
// func (fs *MockFileStat) ModTime() time.Time { return fs.modTime }
// func (fs *MockFileStat) Sys() any           { return &fs.sys }
//
// func (fs *MockFileSystem) Open(path string) (File, error) {
// 	if entry, ok := fs.Files[path]; ok {
// 		return MockFile{Content: entry}, nil
// 	}
// 	// file does not exist
// 	return nil, &os.PathError{Op: "open", Path: path, Err: nil}
// }
//
// func (fs *MockFileSystem) Stat(path string) (os.FileInfo, error) {
// 	// if entry, ok := fs.Files[path]; ok {
// 	// 	return MockFileStat{
// 	// 		name: path,
// 	// 		size: len(entry),
// 	// 		mode: os.FileMode
// 	// 	}, nil
// 	// }
// 	// file does not exist
// 	return nil, &os.PathError{Op: "stat", Path: path, Err: nil}
// }
//
// func (fs *MockFileSystem) DirFS(path string) (realfs.FS, error) {
// 	return os.Stat(path)
// }
