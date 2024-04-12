package m3u8

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/hopeio/cherry/utils/console"
	fs2 "github.com/hopeio/cherry/utils/io/fs"
	"github.com/hopeio/cherry/utils/io/fs/path"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
)

const (
	tsExt         = ".ts"
	progressWidth = 40
)

type Downloader struct {
	lock     sync.Mutex
	queue    []int
	FilePath string
	tsFolder string
	finish   int32
	SegLen   int
	url      string
}

// NewTask returns a Task instance
func NewTask(filePath, tsFolder string, url string) (*Downloader, error) {
	result, err := FromURL(url)
	if err != nil {
		return nil, err
	}
	// If no output folder specified, use current directory
	if filePath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		filePath = pwd + fs2.PathSeparator + filePath
	} else {
		if err := os.MkdirAll(path.GetDirName(filePath), os.ModePerm); err != nil {
			return nil, fmt.Errorf("create storage folder failed: %s", err.Error())
		}
	}

	if err := os.MkdirAll(tsFolder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create ts folder '[%s]' failed: %s", tsFolder, err.Error())
	}
	d := &Downloader{
		FilePath: filePath,
		tsFolder: tsFolder,
		url:      url,
	}
	d.SegLen = len(result.M3u8.Segments)
	d.queue = genSlice(d.SegLen)
	return d, nil
}

// Start runs downloader
func (d *Downloader) Start(concurrency int) error {
	var wg sync.WaitGroup
	// struct{} zero size
	limitChan := make(chan struct{}, concurrency)
	for {
		tsIdx, end, err := d.next()
		if err != nil {
			if end {
				break
			}
			continue
		}
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if err := d.Downloadts(idx); err != nil {
				// Back into the queue, retry request
				fmt.Printf("[failed] %s\n", err.Error())
				if err := d.back(idx); err != nil {
					fmt.Printf(err.Error())
				}
			}
			<-limitChan
		}(tsIdx)
		limitChan <- struct{}{}
	}
	wg.Wait()
	if err := d.Merge(); err != nil {
		return err
	}
	return nil
}

// single thread downloader
func (d *Downloader) Download() error {
	mFile, err := os.Create(d.FilePath)
	if err != nil {
		return fmt.Errorf("create main TS file failed：%s", err.Error())
	}
	//noinspection GoUnhandledErrorResult
	defer mFile.Close()

	w := bufio.NewWriter(mFile)
	for segIndex := 0; segIndex < d.SegLen; segIndex++ {
		result, err := FromURL(d.url)
		if err != nil {
			return err
		}
		data, err := result.Download(segIndex)
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return fmt.Errorf("write to %s: %s", d.FilePath, err.Error())
		}
		w.Flush()
	}
	return nil
}

func (d *Downloader) Downloadts(segIndex int) error {
	tsFilename := tsFilename(segIndex)

	fPath := filepath.Join(d.tsFolder, tsFilename)

	if fs2.NotExist(fPath) {
		result, err := FromURL(d.url)
		if err != nil {
			return err
		}
		data, err := result.Download(segIndex)
		if err != nil {
			return err
		}

		err = fs2.CreatFileFromReader(fPath, bytes.NewReader(data))
		if err != nil {
			return err
		}
	}
	// Maybe it will be safer in this way...
	atomic.AddInt32(&d.finish, 1)
	//tool.DrawProgressBar("Downloading", float32(d.finish)/float32(d.SegLen), progressWidth)
	fmt.Printf("[download %6.2f%%] %s\n", float32(d.finish)/float32(d.SegLen)*100, d.url)
	return nil
}

func (d *Downloader) next() (segIndex int, end bool, err error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if len(d.queue) == 0 {
		err = fmt.Errorf("queue empty")
		if d.finish == int32(d.SegLen) {
			end = true
			return
		}
		// Some segment indexes are still running.
		end = false
		return
	}
	segIndex = d.queue[0]
	d.queue = d.queue[1:]
	return
}

func (d *Downloader) back(segIndex int) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.queue = append(d.queue, segIndex)
	return nil
}

func (d *Downloader) Merge() error {
	// In fact, the number of downloaded segments should be equal to number of m3u8 segments
	missingCount := 0
	for idx := 0; idx < d.SegLen; idx++ {
		tsFilename := tsFilename(idx)
		f := filepath.Join(d.tsFolder, tsFilename)
		if _, err := os.Stat(f); err != nil {
			missingCount++
		}
	}
	if missingCount > 0 {
		fmt.Printf("[warning] %d files missing\n", missingCount)
	}

	// Create a TS file for merging, all segment files will be written to this file.
	mFile, err := os.Create(d.FilePath)
	if err != nil {
		return fmt.Errorf("create main TS file failed：%s", err.Error())
	}
	//noinspection GoUnhandledErrorResult
	defer mFile.Close()

	writer := bufio.NewWriter(mFile)
	mergedCount := 0
	for segIndex := 0; segIndex < d.SegLen; segIndex++ {
		tsFilename := tsFilename(segIndex)
		bytes, err := os.ReadFile(filepath.Join(d.tsFolder, tsFilename))
		_, err = writer.Write(bytes)
		if err != nil {
			continue
		}
		mergedCount++
		console.DrawProgressBar("merge",
			float32(mergedCount)/float32(d.SegLen), progressWidth)
	}
	_ = writer.Flush()
	// Remove `ts` folder
	_ = os.RemoveAll(d.tsFolder)

	if mergedCount != d.SegLen {
		fmt.Printf("[warning] \n%d files merge failed", d.SegLen-mergedCount)
	}

	fmt.Printf("\n[output] %s\n", d.FilePath)

	return nil
}

func tsFilename(ts int) string {
	return strconv.Itoa(ts) + tsExt
}

func genSlice(len int) []int {
	s := make([]int, 0)
	for i := 0; i < len; i++ {
		s = append(s, i)
	}
	return s
}

func (d *Downloader) FfmpegFiles() (string, error) {
	var data bytes.Buffer
	for i := 0; i < d.SegLen; i++ {
		data.WriteString(`file '` + d.tsFolder + "/" + strconv.Itoa(i) + `.ts'
`)
	}
	ffmpegFilePath := d.tsFolder + fs2.PathSeparator + "file.txt"

	file, err := os.Create(ffmpegFilePath)
	if err != nil {
		return "", fmt.Errorf("create ffmpeg file failed：%s", err.Error())
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()
	_, err = file.Write(data.Bytes())
	if err != nil {
		return "", fmt.Errorf("write to %s: %s", ffmpegFilePath, err.Error())
	}
	return ffmpegFilePath, nil
}

func (d *Downloader) Finished() bool {
	return d.finish == int32(d.SegLen)
}

func (d *Downloader) RemoveTmp() error {
	return os.RemoveAll(d.tsFolder)
}
