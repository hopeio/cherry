package fs

import (
	"github.com/hopeio/cherry/utils/errors/multierr"
	"io/fs"
	"os"
	"path/filepath"
)

type FileRangeCallback = func(dir string, entry os.DirEntry) error

// 遍历根目录中的每个文件，为每个文件调用callback,包括文件夹,与filepath.WalkDir不同的是回调函数的参数不同,filepath.WalkDir的第一个参数是文件完整路径,RangeFile是文件所在目录的路径
func Range(dir string, callback FileRangeCallback) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	errs := &multierr.MultiError{}
	for _, entry := range entries {
		if entry.IsDir() {
			err = RangeFile(dir+PathSeparator+entry.Name(), callback)
			if err != nil {
				errs.Append(err)
			}
		}
		err = callback(dir, entry)
		if err != nil {
			errs.Append(err)
		}
	}
	if errs.HasErrors() {
		return errs
	}
	return nil
}

// 遍历根目录中的每个文件，为每个文件调用callback,不包括文件夹,与filepath.WalkDir不同的是回调函数的参数不同,filepath.WalkDir的第一个参数是文件完整路径,RangeFile是文件所在目录的路径
func RangeFile(dir string, callback FileRangeCallback) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	errs := &multierr.MultiError{}
	for _, entry := range entries {
		if entry.IsDir() {
			err = RangeFile(dir+PathSeparator+entry.Name(), callback)
			if err != nil {
				errs.Append(err)
			}
		} else {
			err = callback(dir, entry)
			if err != nil {
				errs.Append(err)
			}
		}
	}
	if errs.HasErrors() {
		return errs
	}
	return nil
}

// RangeDir 遍历根目录中的每个文件夹，为文件夹调用callback
// callback 返回值为需要递归遍历的目录和error
// 几乎每个文件夹下的文件夹都会被循环两次！
func RangeDir(dir string, callback func(dir string, entries []os.DirEntry) ([]os.DirEntry, error)) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	errs := &multierr.MultiError{}

	dirs, err := callback(dir, entries)
	if err != nil {
		errs.Append(err)
	}
	for _, e := range dirs {
		if e.IsDir() {
			err = RangeDir(dir+PathSeparator+e.Name(), callback)
			if err != nil {
				errs.Append(err)
			}
		}
	}
	if errs.HasErrors() {
		return errs
	}
	return nil
}

func WalkDirWithFS(fsys fs.FS, root string, fn fs.WalkDirFunc) error {
	return fs.WalkDir(fsys, root, fn)
}

func Walk(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(root, fn)
}

func WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}
