package middleware

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

// 40bd001563085fc35165329ea1ff5c5ecbdbbeef
func TestSHA1(t *testing.T) {
	sha1 := SHA1("123")
	assert.Equal(t, sha1, "40bd001563085fc35165329ea1ff5c5ecbdbbeef", "ok!")
	t.Logf("SHA1 exec res : %v", sha1)
}
