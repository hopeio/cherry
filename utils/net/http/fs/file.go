package fs

import (
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type File struct {
	File http.File
	Name string
}

type FileInterface interface {
	io.Reader
	Name() string
}

type FileInfo struct {
	name    string
	modTime time.Time
	size    int64
	mode    fs.FileMode
	Binary  []byte
}

func (f *FileInfo) Name() string {
	return f.name
}

func (f *FileInfo) Size() int64 {
	return f.size
}

func (f *FileInfo) Mode() fs.FileMode {
	return f.mode
}

func (f *FileInfo) ModTime() time.Time {
	return f.modTime
}

func (f *FileInfo) IsDir() bool {
	return false
}

func (f *FileInfo) Sys() any {
	return nil
}

type UploadFile struct {
	ID           uint64 `gorm:"primary_key" json:"id"`
	FileName     string `gorm:"type:varchar(100);not null" json:"file_name"`
	OriginalName string `gorm:"type:varchar(100);not null" json:"original_name"`
	URL          string `json:"url"`
	MD5          string `gorm:"type:varchar(32)" json:"md5"`
	Mime         string `json:"mime"`
	Size         uint64 `json:"size"`
}

func GetExt(file *multipart.FileHeader) (string, error) {
	var ext string
	var index = strings.LastIndex(file.Filename, ".")
	if index == -1 {
		return "", nil
	} else {
		ext = file.Filename[index:]
	}
	if len(ext) == 1 {
		return "", errors.New("无效的扩展名")
	}
	return ext, nil
}

func CheckSize(f multipart.File, uploadMaxSize int) bool {
	size := GetSize(f)
	if size == 0 {
		return false
	}

	return size <= uploadMaxSize
}

func GetSize(f multipart.File) int {
	content, err := io.ReadAll(f)
	if err != nil {
		return 0
	}
	return len(content)
}
