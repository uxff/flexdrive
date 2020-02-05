package filehash

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"

	"github.com/uxff/flexdrive/pkg/log"
)

func CalcSha1(fileName string) (string, error) {
	fileHandle, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer fileHandle.Close()

	sha1Enc := sha1.New()
	if _, err := io.Copy(sha1Enc, fileHandle); err != nil {
		log.Errorf("err of doing this")
		return "", err
	}

	hashStr := hex.EncodeToString(sha1Enc.Sum(nil))
	return hashStr, nil
}
