package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

// SHA1 对字符串进行SHA1加密
func SHA1(s string) string {

	o := sha1.New()

	o.Write([]byte(s))

	return hex.EncodeToString(o.Sum(nil))
}
