package utils

//#include <stdlib.h>
//int startCmd(const char* cmd){
//	  return system(cmd);
//}
import "C"

import (
	"douyin/conf"
	"douyin/consts"
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"go.uber.org/zap"
)

type Video2Image struct {
	InputPath  string
	OutputPath string
	StartTime  string
	KeepTime   string
	Filter     string
	FrameCount int64
	debug      bool
}

func NewVideo2Image() *Video2Image {
	return &videoChanger
}

var videoChanger Video2Image

// cmdJoin 执行命令拼接
func cmdJoin(s1, s2 string) string {
	return fmt.Sprintf(" %s %s ", s1, s2)
}

func (v *Video2Image) Debug() {
	v.debug = true
}

func (v *Video2Image) GetQueryString() (string, error) {
	if v.InputPath == "" || v.OutputPath == "" {
		return "", errors.New("输入输出路径未指定")
	}
	var b strings.Builder
	b.WriteString(conf.Conf.Ffmpeg.FfmpegPath)
	b.WriteString(cmdJoin(consts.OptionInputVideoPath, v.InputPath))
	b.WriteString(cmdJoin(consts.OptionFormatToImage, "image2"))
	if v.Filter != "" {
		b.WriteString(cmdJoin(consts.OptionVideoFilter, v.Filter))
	}
	if v.StartTime != "" {
		b.WriteString(cmdJoin(consts.OptionStartTime, v.StartTime))
	}
	if v.KeepTime != "" {
		b.WriteString(cmdJoin(consts.OptionKeepTime, v.KeepTime))
	}
	if v.FrameCount != 0 {
		b.WriteString(fmt.Sprintf(" %s %d", consts.OptionFrames, v.FrameCount))
	}
	b.WriteString(cmdJoin(consts.OptionAutoReWrite, v.OutputPath))
	return b.String(), nil
}

func (v *Video2Image) ExecCommand(cmd string) error {
	if v.debug {
		zap.L().Debug("pkg ffmpeg ExecCommand!", zap.String("command", cmd))
	}
	cCmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(cCmd))
	status := C.startCmd(cCmd)
	if status != 0 {
		return errors.New("视频切截图失败")
	}
	return nil
}
