package util

import "strings"

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
