package util

import (
	"fmt"
	"strconv"
	"strings"
)

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func Contains(arrStr string, str string) bool {
	arr := strings.Split(arrStr, ",")
	b := false
	for _, s := range arr {
		if s == str {
			b = true
			break
		}
	}
	return b
}

func ConvertFileSize(size uint64) string {
	var kb uint64 = 1024
	var mb = kb * 1024
	var gb = mb * 1024
	if size >= gb {
		var g = strconv.FormatUint(size/gb, 10)
		return g + "GB"
	} else if size >= mb {
		var m = strconv.FormatUint(size/mb, 10)
		return m + "MB"
	} else if size >= kb {
		var k = strconv.FormatUint(size/kb, 10)
		return k + "KB"
	} else {
		var g = strconv.FormatUint(size, 10)
		return g + "B"
	}
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func DecimalPercent(value float64) string {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	str := strconv.FormatFloat(value, 'f', 2, 64)
	return str + "%"
}
