package core

import "time"

func timestampToString(timestamp *int64) string {
	if timestamp == nil {
		return ""
	}
	return time.Unix(*timestamp, 0).UTC().Format(time.RFC3339)
}

func valFromPtr(ptr *uint64) uint64 {
	var v uint64
	if ptr != nil {
		v = *ptr
	}
	return v
}

func emptyStrIfZero(s string) string {
	if s == "0" {
		s = ""
	}
	return s
}
