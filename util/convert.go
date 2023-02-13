package util

func StatusConvert(status int) int64 {
	if status != 200 {
		return 1
	}
	return 0
}
