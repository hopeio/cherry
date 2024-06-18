package path

import (
	"github.com/hopeio/cherry/utils/io/fs"
	stringsi "github.com/hopeio/cherry/utils/strings"
	timei "github.com/hopeio/cherry/utils/time"
	"strconv"
	"strings"
	"time"
)

type Path interface {
	Path() string
}

// userId/year/date_userId_key_filename
type ByIdUId struct {
	UserId    int       `json:"userId"`
	UserIdStr string    `json:"userIdStr"`
	Id        int       `json:"id"`
	IdStr     string    `json:"idStr"`
	FileName  string    `json:"fileName"`
	PubAt     time.Time `json:"pubAt" gorm:"type:timestamptz(0);default:0001-01-01 00:00:00"`
	PubAtStr  string    `json:"pubAtStr" comment:"20230321"`
}

func (d *ByIdUId) PreHandle() {
	if d.IdStr == "" {
		d.IdStr = strconv.Itoa(d.Id)
	}
	if d.PubAtStr == "" {
		d.PubAtStr = d.PubAt.Format(timei.LayoutCompactTime)
	}
	if d.UserIdStr == "" {
		d.UserIdStr = strconv.Itoa(d.UserId)
	}
	d.PubAtStr = stringsi.ReplaceRunesEmpty(d.PubAtStr, '-', ' ', ':')
}

func (d *ByIdUId) Path() string {
	d.PreHandle()
	filepath := strings.Join([]string{d.UserIdStr, d.PubAtStr[:4], strings.Join([]string{d.PubAtStr, d.UserIdStr, d.IdStr, d.FileName}, "_")}, fs.PathSeparator)
	return filepath
}

// year/key_date_title/filename
type ById struct {
	Id    int    `json:"id"`
	IdStr string `json:"idStr"`
	Title string
	//PrePath   string    `json:"prePath" comment:""`
	DateStr  string
	FileName string `json:"fileName"`
}

func (d *ById) Path() string {
	if d.IdStr == "" {
		d.IdStr = strconv.Itoa(d.Id)
	}
	filepath := strings.Join([]string{d.DateStr[:4], d.IdStr + "_" + d.DateStr + "_" + d.Title, d.FileName}, "/")
	return filepath
}

// filename
type ByPath string

func (d ByPath) Path() string {
	return string(d)
}

// userId/date_key_title/date_userId_key_filename
type ByIdUIdTitle struct {
	ByIdUId
	Title string
	//PrePath   string    `json:"prePath" comment:""`
}

func (d *ByIdUIdTitle) Path() string {
	d.PreHandle()
	filepath := strings.Join([]string{d.UserIdStr, d.PubAtStr + "_" + d.Title + "_" + d.IdStr, strings.Join([]string{d.PubAtStr, d.UserIdStr, d.IdStr, d.FileName}, "_")}, "/")
	return filepath
}
