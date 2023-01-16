package util

import "time"

//Seconds timestamp to "mm-dd" date string
func Timestamp2Date(seconds int64) string {
	return time.Unix(seconds, 0).Format("2006-01-02 15:04:05")[5:10]
}
