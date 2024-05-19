package fs

import (
	"github.com/hopeio/cherry/utils/crypto/md5"
	"io"
	"os"
)

type mode int

const (
	Cover mode = iota
	SameNameSkip
	SameNameAndMd5Skip
	// TODO
	sameNameRename
)

func (c mode) handle(dst string, src io.Reader) (newname string, skip bool, err error) {
	switch c {
	case Cover:
		return "", false, nil
	case SameNameSkip:
		if IsExist(dst) {
			return "", true, nil
		}
	case SameNameAndMd5Skip:
		if IsExist(dst) {
			md51, err := Md5(dst)
			if err != nil {
				return "", false, err
			}
			md52, err := md5.EncodeReaderString(src)
			if err != nil {
				return "", false, err
			}

			if md51 == md52 {
				return "", true, nil
			}
		}
	}
	return "", false, nil
}

// CopyFile : General Approach
func CopyFile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	return CreatFileFromReader(dst, r)
}

func CopyFileByMode(src, dst string, c mode) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()
	_, skip, err := c.handle(dst, r)
	if err != nil {
		return err
	}
	if skip {
		return nil
	}
	return CreatFileFromReader(dst, r)
}

const DownloadKey = ".downloading"

func CreatFileFromReader(filepath string, reader io.Reader) error {
	f, err := Create(filepath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(f, reader); err != nil {
		f.Close()
		os.Remove(filepath)
		return err
	}

	if err = f.Close(); err != nil {
		os.Remove(filepath)
		return err
	}
	return nil
}

func CreatFileFromReaderByMode(filepath string, reader io.Reader, c mode) error {
	_, skip, err := c.handle(filepath, reader)
	if err != nil {
		return err
	}
	if skip {
		return nil
	}
	return CreatFileFromReader(filepath, reader)
}

func DownloadFile(filepath string, reader io.Reader) error {
	tmpFilepath := filepath + DownloadKey
	err := CreatFileFromReader(tmpFilepath, reader)
	if err != nil {
		return err
	}
	return os.Rename(tmpFilepath, filepath)
}

func DownloadFileByMode(filepath string, reader io.Reader, c mode) error {
	_, skip, err := c.handle(filepath, reader)
	if err != nil {
		return err
	}
	if skip {
		return nil
	}
	return DownloadFile(filepath, reader)
}

// CopyDirByMode 递归复制目录
func CopyDirByMode(src, dst string, c mode) error {
	if src[len(src)-1] == os.PathSeparator {
		src = src[:len(src)-1]
	}
	if dst[len(dst)-1] == os.PathSeparator {
		dst = dst[:len(dst)-1]
	}
	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if len(entries) == 0 {
		return nil
	}
	for _, entry := range entries {
		entityName := entry.Name()
		if entry.IsDir() {
			err = CopyDirByMode(src+PathSeparator+entityName, dst+PathSeparator+entityName, c)
			if err != nil {
				return err
			}
		} else {
			err = CopyFileByMode(src+PathSeparator+entityName, dst+PathSeparator+entityName, c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyDir 递归复制目录
func CopyDir(src, dst string) error {
	return CopyDirByMode(src, dst, Cover)
}
