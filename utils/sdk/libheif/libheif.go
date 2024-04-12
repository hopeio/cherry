package libheif

import (
	"fmt"
	osi "github.com/hopeio/cherry/utils/os"
	"strings"
)

// https://github.com/pphh77/libheif-Windowsbinary/releases
// https://github.com/strukturag/libheif 纯库,需要自己编译
const ImgToHeifCmd = `heif-enc -q 50 %s -o %s.heif`
const ImgToHeifCmd1 = `heif-enc -p x265:crf=20.5 -p x265:colorprim=smpte170m -p x265:rdoq-level=1 -p x265:aq-strength=1.2 -p x265:deblock=-2:-2 %s -o %s.heif
`

func ImgToHeif(filePath, dst string) error {
	if strings.HasSuffix(dst, ".heif") {
		dst = dst[:len(dst)-5]
	}
	_, err := osi.ContainQuotedCMD(fmt.Sprintf(ImgToHeifCmd, filePath, dst))

	return err
}
