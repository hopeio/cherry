package decimal

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

type Decimal3 struct {
	mant uint64
	exp  int
	neg  bool
}

func New3(mant uint64, exp int, neg bool) *Decimal3 {
	if mant == 0 {
		return &Decimal3{}
	}
	mantStr := strconv.FormatUint(mant, 10)
	for i := len(mantStr) - 1; i >= 0; i-- {
		if mantStr[i] == '0' {
			mant /= 10
			exp += 1
		} else {
			break
		}
	}
	return &Decimal3{
		mant: mant,
		exp:  exp,
		neg:  neg,
	}
}

func New3FromStr(str string) (Decimal3, error) {
	var dec Decimal3
	if str != "" && str[0] == '-' {
		dec.neg = true
		str = str[1:]
	}

	nums := strings.Split(str, ".")
	if len(nums) == 2 {
		for i := len(nums[1]) - 1; i >= 0; i-- {
			if nums[1][i] == '0' {
				nums[1] = nums[1][:i]
			} else {
				break
			}
		}
		dec.exp = -len(nums[1])
		str = nums[0] + nums[1]
	} else {
		str = nums[0]
	}

	dec.mant, _ = strconv.ParseUint(str, 10, 64)

	return dec, nil
}

func (x *Decimal3) Add(v Decimal3) *Decimal3 {
	var dec = *x

	if x.exp > v.exp {
		dec.exp = v.exp
		dec.mant = dec.mant * uint64(math.Pow10(x.exp-v.exp))
	} else if x.exp < v.exp {
		v.mant = v.mant * uint64(math.Pow10(v.exp-x.exp))
		v.exp = x.exp
	}

	if x.neg == v.neg {
		dec.mant += v.mant
	} else {
		if x.mant >= v.mant {
			dec.mant -= v.mant
		} else {
			dec.mant = v.mant - dec.mant
			dec.neg = v.neg
		}
	}

	return &dec
}

func (x *Decimal3) Sub(v Decimal3) *Decimal3 {
	v.neg = !v.neg
	return x.Add(v)
}

func (x *Decimal3) Mul(v Decimal3) *Decimal3 {
	if x.mant == 0 || v.mant == 0 {
		return &Decimal3{}
	}
	v.mant *= x.mant
	v.exp += x.exp
	if x.neg != v.neg {
		v.neg = true
	}
	return &v
}

type RoundingMode uint8

const (
	UP RoundingMode = iota
	DOWN
	HALFUP
)

func (x *Decimal3) Div(v Decimal3, mode RoundingMode) *Decimal3 {
	if x.mant == 0 {
		return &Decimal3{}
	}
	if v.mant == 0 {
		panic("除数不能为0")
	}
	if v.exp == 0 {
		return &*x
	}
	d1 := x.mant
	d2 := v.mant

	if v.exp < 0 {
		d1 = x.mant * uint64(math.Pow10(0-v.exp))
	} else {
		v.exp = x.exp - v.exp
	}

	v.mant = d1 / d2
	if x.neg != v.neg {
		v.neg = true
	}
	return &v
}

func (x *Decimal3) DivInt(v int, mode RoundingMode) *Decimal3 {
	dec := *x
	dec.mant = dec.mant / uint64(v)
	return &dec
}

func (x *Decimal3) Float() float64 {
	if x.neg {
		return -float64(x.mant) / math.Pow10(abs(x.exp))
	}
	return float64(x.mant) / math.Pow10(abs(x.exp))
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

/*
	func (x *Decimal3) Div2(v Decimal3) *Decimal3 {
		f1:=big.NewFloat(x.Float())
		f2:=big.NewFloat(v.Float())
		f3,_:=f1.Quo(f1,f2).Float64()
	}
*/
func (x Decimal3) String() string {
	if x.mant == 0 {
		return "0"
	}
	d := x.mant
	if x.exp > 0 {
		d = x.mant * uint64(math.Pow10(x.exp))
	}

	str := strconv.FormatUint(d, 10)
	var in, de string
	if x.neg == true {
		in = "-"
	}
	if len(str)+x.exp < 0 {
		in += "0"
		de = "." + strings.Repeat("0", -x.exp-len(str)) + str
	} else {
		if x.exp >= 0 {
			in += str
		} else {
			in += str[:len(str)+x.exp]
			de = "." + str[len(str)+x.exp:]
		}

	}

	return fmt.Sprintf("%s%s", in, de)
}

func (x *Decimal3) Decompose(buf []byte) (form byte, negative bool, mant []byte, exponent int32) {
	return 0, x.neg, (*(*[8]byte)(unsafe.Pointer(&x.mant)))[:], int32(x.exp)
}

func (x *Decimal3) Compose(form byte, negative bool, mant []byte, exponent int32) error {
	x.neg = negative
	x.exp = int(exponent)
	var array [8]byte
	copy(array[:], mant)
	x.mant = *(*uint64)(unsafe.Pointer(&array))
	return nil
}

func (x *Decimal3) Scan(raw interface{}) error {
	var err error
	switch v := raw.(type) {
	case []byte:
		*x, err = New3FromStr(string(v))
	case string:
		*x, err = New3FromStr(v)
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Base64 from: %#v", v)
	}
	if err != nil {
		return err
	}
	return nil
}

func (x Decimal3) Value() (driver.Value, error) {
	return driver.Value(x.String()), nil
}
