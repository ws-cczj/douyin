package utils

import (
	"douyin/pkg/e"
	"fmt"
	"github.com/importcjj/sensitive"
)

var filter *sensitive.Filter

const WordDictPath = "./pkg/document/sensitiveDict.txt"

func InitFilter() {
	filter = sensitive.New()
	err := filter.LoadWordDict(WordDictPath)
	if err != nil {
		panic(fmt.Sprintf("%s, err: %v", e.FailInitFilter.Msg(), err))
	}
}

// Replace 替换字符!
func Replace(content string) string {
	return filter.Replace(content, '*')
}
