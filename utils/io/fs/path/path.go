package path

import (
	stringsi "github.com/hopeio/cherry/utils/strings"
	sdpath "path"
	"strings"
)

// windows需要,由于linux的文件也要放到windows看,统一处理
func FileNameRewrite(filename string) string {
	filename = stringsi.ReplaceRunesEmpty(filename, '/', '\\', '*', '|')
	filename = strings.ReplaceAll(filename, "<", "《")
	filename = strings.ReplaceAll(filename, ">", "》")
	filename = strings.ReplaceAll(filename, "?", "？")
	filename = strings.ReplaceAll(filename, ":", "：")
	return filename
}

// 仅仅针对文件名
func FileNameClean(filename string) string {

	filename = strings.Trim(filename, ".")
	filename = strings.TrimPrefix(filename, "-")
	filename = strings.TrimPrefix(filename, "+")
	// windows
	//filename = stringsi.ReplaceRunesEmpty(filename, '/', '\\', ':', '*', '?', '"', '<', '>', '|')
	// linux
	//filename = stringsi.ReplaceRunesEmpty(filename, '\'', '*','?', '@', '#', '$', '&', '(', ')', '|', ';',  '/', '%', '^', ' ', '\t', '\n')

	filename = stringsi.ReplaceRunesEmpty(filename, '/', '\\', ':', '*', '?', '"', '<', '>', '|', ';', '/', '%', '^', ' ', '\t', '\n', '$', '&')
	// 中文符号
	//filename = stringsi.ReplaceRunesEmpty(filename, '：', '，', '。', '！', '？', '、', '“', '”', '、')
	return filename
}

// 仅仅针对目录名
func DirNameClean(dir string) string { // will be used when save the dir or the part
	// remove special symbol
	// :unix允许存在，windows需要
	// windows path
	if len(dir) > 2 && dir[1] == ':' && ((dir[0] >= 'A' && dir[0] <= 'Z') || (dir[0] >= 'a' && dir[0] <= 'z')) && (dir[2] == '/' || dir[2] == '\\') {
		return dir[:3] + stringsi.ReplaceRunesEmpty(dir[3:], ':', '*', '?', '"', '<', '>', '|', ',', ' ', '\t', '\n')
	}
	return stringsi.ReplaceRunesEmpty(dir, ':', '*', '?', '"', '<', '>', '|', ',', ' ', '\t', '\n')
}

// 针对带目录的完整文件名
func PathClean(path string) string { // will be used when save the dir or the part
	dir, file := Split(path)
	if dir == "" {
		return DirNameClean(dir)
	}
	if file == "" {
		return FileNameClean(file)
	}
	// remove special symbol
	return DirNameClean(dir) + string(path[len(dir)-1-len(file)]) + FileNameClean(file)
}

// 获取url的目录部分
func GetDirName(path string) string {
	dir, _ := Split(path)
	return sdpath.Clean(dir)
}

// 获取url的文件部分
func GetFileName(path string) string {
	_, file := Split(path)
	return file
}

// 返回目录名和文件名
func Split(path string) (dir, file string) {
	i := lastSlash(path)
	return path[:i+1], path[i+1:]
}

// 获取文件名除去扩展名
func FileNameExcludeExt(filepath string) string {
	base := sdpath.Base(filepath)
	return base[:len(base)-len(sdpath.Ext(base))]
}
