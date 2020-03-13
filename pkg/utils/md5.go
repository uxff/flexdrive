package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string) string {
	enc := md5.New()
	enc.Write([]byte(str))
	return hex.EncodeToString(enc.Sum(nil))
}
