package util

import "strconv"

func Str2Int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}
