package middleware

import (
	"douyin/consts"
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var videoIndexMap = map[string]struct{}{
	".mp4":  {},
	".avi":  {},
	".wmv":  {},
	".flv":  {},
	".mpeg": {},
	".mov":  {},
}

// Ffmpeg 使用Ffmpeg工具进行封面截取和保存
func Ffmpeg(isDebug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64("user_id")
		file, err := c.FormFile("data")
		if err != nil {
			zap.L().Error("middleware ffmpeg FormFile get data fail!", zap.Error(err))
			common.FailWithCode(c, e.FailParamInvalid)
			c.Abort()
			return
		}
		// 获取视频文件后缀进行格式校验
		suffix := filepath.Ext(file.Filename)
		if _, ok := videoIndexMap[suffix]; !ok {
			zap.L().Error("middleware ffmpeg video illegal!", zap.String("suffix", suffix))
			common.FailWithCode(c, e.FailVideoIllegal)
			c.Abort()
			return
		}
		// 获取拼接名称
		name := fmt.Sprintf("%d-%d", time.Now().Unix(), userId)
		newName := fmt.Sprintf("%s%s", name, suffix)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			newVideoPath := filepath.Join("./public", newName)
			// 保存视频到本地
			if err = c.SaveUploadedFile(file, newVideoPath); err != nil {
				zap.L().Error("middleware ffmpeg video Save fail!", zap.Error(err))
			}
		}()
		// 保存图片到本地
		if err = utils.SaveImageFromVideo(name, isDebug); err != nil {
			zap.L().Error("middleware ffmpeg video cover Save fail!", zap.Error(err))
		}
		wg.Wait()
		if err != nil {
			common.FailWithCode(c, e.FailServerBusy)
			c.Abort()
			return
		}
		c.Set("play_url", utils.GetFileUrl(newName))
		c.Set("cover_url", utils.GetPicUrl(fmt.Sprintf("%s%s", name, consts.DefaultImageSuffix)))
		c.Next()
	}
}
