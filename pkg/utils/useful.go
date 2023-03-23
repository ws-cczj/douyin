package utils

import (
	"douyin/conf"
	"douyin/consts"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
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

// StrI64 拼接字符串和int64
func StrI64(key string, id int64) string {
	// return fmt.Sprintf("%s:%d", key, id)
	return key + I64toa(id)
}

// ISlice64toa intSlice64转为string
func ISlice64toa(ks []int64) string {
	var b strings.Builder
	b.Grow(len(ks)*2 - 1)
	for i, k := range ks {
		b.WriteString(I64toa(k))
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

// GetFileUrl 获取文件URL
func GetFileUrl(fileName string) string {
	return fmt.Sprintf("http://%s:%d/static/%s", conf.Conf.Ip, conf.Conf.Port, fileName)
}

// GetPicUrl 获取封面URL
func GetPicUrl(fileName string) string {
	return fmt.Sprintf("http://%s:%d/static/pic/%s", conf.Conf.Ip, conf.Conf.Port, fileName)
}

// SaveImageFromVideo 将视频切一帧保存到本地
// isDebug用于控制是否打印出执行的ffmpeg命令
func SaveImageFromVideo(name string, isDebug bool) error {
	v2i := NewVideo2Image()
	if isDebug {
		v2i.Debug()
	}
	v2i.InputPath = filepath.Join(conf.Conf.PublicPath, name+consts.DefaultVideoSuffix)
	v2i.OutputPath = filepath.Join(conf.Conf.PublicPicPath, name+consts.DefaultImageSuffix)
	v2i.FrameCount = 1
	queryString, err := v2i.GetQueryString()
	if err != nil {
		zap.L().Error("middleware useful GetQueryString method exec fail!", zap.Error(err))
		return err
	}
	return v2i.ExecCommand(queryString)
}

// FormatTime 格式化评论时间
func FormatTime(time time.Time) string {
	return time.Format("15:04:05")
}
