package libjxl

import (
	"fmt"
	osi "github.com/hopeio/cherry/utils/os"
	"strings"
)

// https://github.com/libjxl/libjxl/releases
// windows support: https://github.com/saschanaz/jxl-winthumb/releases administrator regsvr32 jxl_winthumb.dll
const ImgToJxlCmd = `cjxl %s %s.jxl -q %d --lossless_jpeg=0`
const JxlImgToOtherCmd = `djxl %s %s`

func ImgToJxl(filePath, dst string, quality int) error {
	if strings.HasSuffix(dst, ".jxl") {
		dst = dst[:len(dst)-4]
	}
	_, err := osi.Cmd(fmt.Sprintf(ImgToJxlCmd, filePath, dst, quality))
	return err
}
