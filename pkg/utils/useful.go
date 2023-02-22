package utils

import (
	"douyin/consts"
	"strconv"
	"strings"
)

// I64toa int64转为string
func I64toa(k int64) string {
	return strconv.FormatInt(k, 10)
}

// AtoI64 string转为int64
func AtoI64(k string) int64 {
	res, _ := strconv.ParseInt(k, 10, 64)
	return res
}

// ISlice64toa intSlice64转为string
func ISlice64toa(ks []int64) string {
	var b strings.Builder
	b.Grow(len(ks)*2 - 1)
	for i, k := range ks {
		b.WriteString(strconv.FormatInt(k, 10))
		if i != len(ks)-1 {
			b.WriteString(",")
		}
	}
	return b.String()
}

// SearchZero 搜索数组中0的位置
func SearchZero(ary []int64) int {
	l, r := 0, len(ary)
	for l < r {
		mid := (l + r) >> 1
		if ary[mid] == 0 {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

// AddCacheKey 拼接缓存key
func AddCacheKey(key ...string) string {
	var b strings.Builder
	b.Grow(len(key) + 1)
	b.WriteString(consts.CacheDouyin)
	for _, k := range key {
		b.WriteString(k)
	}
	return b.String()
}
