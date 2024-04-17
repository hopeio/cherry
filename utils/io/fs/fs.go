package fs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/hopeio/cherry/utils/io/fs/path"
	"github.com/hopeio/cherry/utils/log"
	runtimei "github.com/hopeio/cherry/utils/scheduler/monitor"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Dir string

func (d Dir) Open(name string) (*os.File, error) {
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	fullName := filepath.Join(dir, filepath.FromSlash(filepath.Clean(string(os.PathSeparator)+name)))
	f, err := os.Open(fullName)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// path和filepath两个包，filepath文件专用
func FindFile(path string) (string, error) {
	files, err := FindFiles(path, 8, 1)
	if err != nil {
		return "", err
	}
	return files[0], nil
}

func FindFiles(path string, deep int8, num int) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var files []string
	filepath1 := filepath.Join(wd, path)
	if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
		files = append(files, filepath1)
		if len(files) == num {
			return files, nil
		}
	}

	subDirFiles(wd, path, "", &files, deep, 0, num)
	supDirFiles(wd+string(os.PathSeparator), path, &files, deep, 0, num)
	if len(files) == 0 {
		return nil, errors.New("找不到文件")
	}
	return files, nil
}

func subDirFiles(dir, path, exclude string, files *[]string, deep, step int8, num int) {
	step += 1
	if step-1 == deep {
		return
	}
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Error(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			if exclude != "" && fileInfos[i].Name() == exclude {
				continue
			}
			filepath1 := filepath.Join(dir, fileInfos[i].Name(), path)
			if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
				*files = append(*files, filepath1)
				if len(*files) == num {
					return
				}
			}
			subDirFiles(filepath.Join(dir, fileInfos[i].Name()), path, "", files, deep, step, num)
		}
	}
}

func supDirFiles(dir, path string, files *[]string, deep, step int8, num int) {
	step += 1
	if step-1 == deep {
		return
	}
	dir, dirName := filepath.Split(dir[:len(dir)-1])
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return
	}
	filepath1 := filepath.Join(dir, path)
	if _, err := os.Stat(filepath1); !os.IsNotExist(err) {
		*files = append(*files, filepath1)
		if len(*files) == num {
			return
		}
	}
	subDirFiles(dir, path, dirName, files, deep, 0, num)
	supDirFiles(dir, path, files, deep, step, num)
}

// path和filepath两个包，filepath文件专用
func FindFiles2(path string, deep int8, num int) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var file = make(chan string, 1)
	//属于回调而不是通知
	ctx := runtimei.New(context.Background(), func() {
		close(file)
	})
	defer ctx.Cancel()
	// 当前目录下先找
	filepath1 := filepath.Join(wd, path)
	if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
		file <- filepath1
	}

	ctx.Run(func() {
		err := subDirFiles2(wd, path, "", file, deep, 0, ctx)
		if err != nil {
			log.Error(err)
		}
	})

	ctx.Run(func() {
		err := supDirFiles2(wd+string(os.PathSeparator), path, file, deep, 0, ctx)
		if err != nil {
			log.Error(err)
		}
	})
	var files []string
	for filepath1 := range file {
		if files = append(files, filepath1); len(files) == num {
			//close(file) 这里无需做关闭操作，会关的
			return files, nil
		}
	}
	return files, nil
}

func subDirFiles2(dir, path, exclude string, file chan string, deep, step int8, ctx *runtimei.Monitor) error {

	step += 1
	if step-1 == deep {
		return nil
	}
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			if exclude != "" && fileInfos[i].Name() == exclude {
				continue
			}
			subDir := filepath.Join(dir, fileInfos[i].Name())
			filepath1 := filepath.Join(subDir, path)
			if _, err = os.Stat(filepath1); !os.IsNotExist(err) {
				//①如果给出了default语句，那么就会执行default的流程，同时程序的执行会从select语句后的语句中恢复。
				//②如果没有default语句，那么select语句将被阻塞，直到至少有一个case可以进行下去。
				select {
				case <-ctx.Done():
					return ctx.Err()
				case file <- filepath1:
				}
			}
			ctx.Run(func() {
				//TODO: subDirFiles2(filepath.Join(dir, fileInfos[i].Name()), path, "", file, deep, step, ctx) 这种语句有问题,不清楚原因,i会错乱?
				// 还是go1.22之前的问题，升级go mod go 和toolchain 标签到1.22
				err := subDirFiles2(filepath.Join(dir, fileInfos[i].Name()), path, "", file, deep, step, ctx)
				if err != nil {
					log.Error(err)
				}
			})

		}
	}
	return nil
}

func supDirFiles2(dir, path string, file chan string, deep, step int8, ctx *runtimei.Monitor) error {
	step += 1
	if step-1 == deep {
		return nil
	}
	dir, dirName := filepath.Split(dir[:len(dir)-1])
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return err
	}
	filepath1 := filepath.Join(dir, path)
	if _, err := os.Stat(filepath1); !os.IsNotExist(err) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case file <- filepath1:
		}
	}

	ctx.Run(func() {
		err := subDirFiles2(dir, path, dirName, file, deep, 0, ctx)
		if err != nil {
			log.Error(err)
		}
	})
	ctx.Run(func() {
		err := supDirFiles2(dir, path, file, deep, step, ctx)
		if err != nil {
			log.Error(err)
		}
	})
	return nil
}

func Mkdir(src string) error {
	_, err := os.Stat(src)
	if os.IsNotExist(err) {
		err = os.Mkdir(src, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return err
}

func MkdirAll(src string) error {
	return os.MkdirAll(src, os.ModePerm)
}

func CheckExist(src string) bool {
	_, err := os.Stat(src)
	return !os.IsNotExist(err)
}

func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + PathSeparator + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = Mkdir(src)
	if err != nil {
		return nil, fmt.Errorf("mkdir src: %s, err: %v", src, err)
	}

	f, err := os.OpenFile(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

func GetLogFilePath(RuntimeRootPath, LogSavePath string) string {
	return RuntimeRootPath + PathSeparator + LogSavePath + PathSeparator
}

func OpenLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("权限不足 src: %s", src)
	}

	err = Mkdir(src)
	if err != nil {
		return nil, fmt.Errorf("文件不存在 src: %s, err: %v", src, err)
	}

	f, err := os.OpenFile(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开失败 :%v", err)
	}

	return f, nil
}

func Create(filepath string) (*os.File, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		dir := path.GetDirName(filepath)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return os.Create(filepath)
}

func Open(filepath string) (*os.File, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		dir := path.GetDirName(filepath)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
		return os.Create(filepath)
	}
	return os.OpenFile(filepath, os.O_RDWR, 0666)
}

// LastFile 当前目录最后一个创建的文件
func LastFile(dir string) (os.FileInfo, map[string]os.FileInfo, error) {
	entries, err := os.ReadDir(dir)
	if len(entries) == 0 {
		return nil, nil, err
	}
	sort.Sort(DirEntries(entries))
	lastFile, err := entries[0].Info()
	if err != nil {
		return nil, nil, err
	}
	m := make(map[string]os.FileInfo)
	for _, entity := range entries {
		m[entity.Name()], _ = entity.Info()
	}
	return lastFile, m, nil
}

// CopyDir 递归复制目录
func CopyDir(src, dst string) error {
	if src[len(src)-1] == os.PathSeparator {
		src = src[:len(src)-1]
	}
	if dst[len(dst)-1] == os.PathSeparator {
		dst = dst[:len(dst)-1]
	}
	_, err := os.Stat(dst)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dst, os.ModePerm)
		if err != nil {
			return err
		}
	}
	entries, err := os.ReadDir(src)
	if len(entries) == 0 {
		return nil
	}
	for _, entry := range entries {
		entityName := entry.Name()
		if entry.IsDir() {
			err = CopyDir(src+PathSeparator+entityName, dst+PathSeparator+entityName)
			if err != nil {
				return err
			}
		} else {
			_, err = os.Stat(dst + PathSeparator + entityName)
			if os.IsNotExist(err) {
				err = CopyFile(src+PathSeparator+entityName, dst+PathSeparator+entityName)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Check(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = Mkdir(dir + PathSeparator + src)
	if err != nil {
		return fmt.Errorf("mkdir err: %v", err)
	}

	perm := CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}

func Move(src, dst string) error {
	dir := path.GetDirName(dst)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return os.Rename(src, dst)
}

type FileSize int64

// MarshalText
func (f FileSize) MarshalText() ([]byte, error) {
	buffer := bytes.NewBufferString("")
	if f/FileSize(1024*1024*1024*8) > 0 {
		buffer.WriteString(fmt.Sprintf("%.2f", float64(f)/float64(1024*1024*1024*8)))
		buffer.WriteString("GB")
	} else if f/FileSize(1024*1024*8) > 0 {
		buffer.WriteString(fmt.Sprintf("%.2f", float64(f)/float64(1024*1024*8)))
		buffer.WriteString("MB")
	} else if f/FileSize(1024*8) > 0 {
		buffer.WriteString(fmt.Sprintf("%.2f", float64(f)/float64(1024*8)))
		buffer.WriteString("KB")
	} else {
		buffer.WriteString(fmt.Sprintf("%d", f/8))
		buffer.WriteString("B")
	}
	return buffer.Bytes(), nil
}

// UnMarshalText
func (f *FileSize) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	unitLen := 1
	unit := text[len(text)-1]
	if unit == 'B' {
		unit = text[len(text)-2]
	}
	switch unit {
	case 'G', 'g':
		unitLen = 2
		size, err := strconv.Atoi(string(text[:len(text)-unitLen]))
		if err != nil {
			return err
		}
		*f = FileSize(size * 1024 * 1024 * 1024 * 8)
	case 'M', 'm':
		unitLen = 2
		size, err := strconv.Atoi(string(text[:len(text)-unitLen]))
		if err != nil {
			return err
		}
		*f = FileSize(size * 1024 * 1024 * 8)
	case 'K', 'k':
		unitLen = 2
		size, err := strconv.Atoi(string(text[:len(text)-unitLen]))
		if err != nil {
			return err
		}
		*f = FileSize(size * 1024)
	default:
		unitLen = 1
		size, err := strconv.Atoi(string(text[:len(text)-unitLen]))
		if err != nil {
			return err
		}
		*f = FileSize(size * 8)
	}
	return nil
}
