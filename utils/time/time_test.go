package time

import (
	"encoding/json"
	"testing"
	"time"
)

type Foo struct {
	T1 UnixTime[MilliTime]
	T2 UnixTime[NanoTime]
}

func TestType(t *testing.T) {
	foo := Foo{T1: UnixTime[MilliTime](time.Now()), T2: UnixTime[NanoTime](time.Now())}
	data, _ := json.Marshal(&foo)
	t.Log(string(data))
}

func TestTimestamp(t *testing.T) {
	t.Log(time.Unix(-62135596800, 0)) // 0001-01-01 08:00:00 +0800 CST
	t.Log(time.Unix(-62135596899, 0)) // 0001-01-01 07:58:21 +0800 CST
}

type Foo1 struct {
	T1 DateTime
	T2 Date
	T3 UnixSecondTime
	T4 UnixNanoTime
}

func TestUnionTime(t *testing.T) {

	foo := Foo1{T1: DateTime(time.Now()),
		T2: Date(time.Now()),
		T3: NewUnixTime[SecondTime](time.Now()),
		T4: NewUnixTime[NanoTime](time.Now()),
	}
	data, _ := json.Marshal(&foo)
	t.Log(string(data)) // {"T1":"2023-02-09 15:00:49","T2":"2023-02-09","T3":1675926049,"T4":1675926049057035300}
	data = []byte(`{"T1":"2023-02-09 15:00:49","T2":"2023-02-09","T3":1675926049,"T4":1675926049057035300}`)
	foo1 := Foo1{
		T1: DateTime{},
		T2: Date{},
		T3: UnixSecondTime{},
		T4: UnixNanoTime{},
	}
	json.Unmarshal(data, &foo1)
	t.Log(foo1)
}

func TestTimeScan(t *testing.T) {
	var d interface{}
	d = Date(time.Now())
	switch s := d.(type) {
	case time.Time:
		t.Log(s)
	case Date:
		t.Log(s)
	}
	d = SecondTimestamp(1)
	switch s := d.(type) {
	case int64:
		t.Log(s)
	case SecondTimestamp:
		t.Log(s)
	}
}

func TestTimeStampScan(t *testing.T) {
	tm := time.Now()
	t.Log(tm)
	ts := tm.Unix()
	t.Log(ts)
	t.Log(time.Unix(ts, 0))
	t.Log(time.Unix(ts, 0).Local())
}
