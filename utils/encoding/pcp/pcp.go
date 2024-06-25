package pcp

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type PCP struct {
	Header Header
	Data   []PointXYZGray // [][][]int
}

type Header struct {
	Type             string
	Size             int
	Width            int
	Height           int
	Channel          int
	IntervalValid    int
	Minx, Miny, Minz int
	Maxx, Maxy, Maxz int
	Dx, Dy           int
	Expose           int
	Timestamp        int
	DeviceMessage    string
	decimal          int
}

type PointXYZGray struct {
	X, Y, Z int
	Gray    int
}

func Parse(filePath string) error {
	var pcp PCP
	m := toMap(&pcp.Header)
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file error:%w", err)
	}
	var isHeader bool
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		str := scanner.Text()
		if str == "header" {
			isHeader = true
			continue
		}
		if str == "data" {
			isHeader = false
			pcp.Data = make([]PointXYZGray, pcp.Header.Size)
			continue
		}
		if isHeader {
			kv := strings.Split(str, ":")
			if len(kv) != 2 {
				return fmt.Errorf("parse header error:%s", str)
			}
			if kv[1] == "(null)" {
				continue
			}
			if v, ok := m[strings.ToUpper(kv[0])]; ok {
				switch v.Kind() {
				case reflect.Int:
					num := kv[1]
					if strings.Contains(kv[1], ".") {
						num = trimZero(strings.ReplaceAll(kv[1], ".", ""))
						if pcp.Header.decimal == 0 {
							pcp.Header.decimal = len(kv[1][strings.Index(kv[1], "."):])
						}
					}

					istr, err := strconv.Atoi(num)
					if err != nil {
						return err
					}
					v.SetInt(int64(istr))

				case reflect.String:
					v.SetString(kv[1])
				}
			}
		} else {
			values := strings.Split(str, " ")
			if len(values) != 4 {
				return fmt.Errorf("parse data error:%s", str)
			}

			x, err := strconv.Atoi(trimZero(strings.ReplaceAll(values[0], ".", "")))
			if err != nil {
				return err
			}
			y, err := strconv.Atoi(trimZero(strings.ReplaceAll(values[1], ".", "")))
			if err != nil {
				return err
			}
			z, err := strconv.Atoi(trimZero(strings.ReplaceAll(values[2], ".", "")))
			if err != nil {
				return err
			}
			gary, err := strconv.Atoi(trimZero(strings.ReplaceAll(values[3], ".", "")))
			if err != nil {
				return err
			}

			pcp.Data = append(pcp.Data, PointXYZGray{X: x, Y: y, Z: z, Gray: gary})

		}
	}
	return nil
}

func toMap(v any) map[string]reflect.Value {
	m := make(map[string]reflect.Value)
	value := reflect.ValueOf(v).Elem()
	for i := 0; i < value.NumField(); i++ {
		m[strings.ToUpper(value.Type().Field(i).Name)] = value.Field(i)
	}
	return m
}

func trimZero(str string) string {
	for i := 0; i < len(str); i++ {
		if str[i] != '0' {
			return str[i:]
		}
	}
	return "0"
}
