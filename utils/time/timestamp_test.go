package time

import (
	"log"
	"testing"
	"time"
)

func TestTime2(t *testing.T) {
	var tm Timestamp = 1572838282583
	log.Println(tm.Time())
}

func TestTimeAdd(t *testing.T) {
	log.Println(time.Now().AddDate(0, 0, -16))
}
