package filehash

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"

	"github.com/uxff/flexdrive/pkg/log"
)

// 可计算大文件的hash
func CalcSha1(fileName string) (string, error) {
	fileHandle, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer fileHandle.Close()

	sha1Enc := sha1.New()
	if _, err := io.Copy(sha1Enc, fileHandle); err != nil {
		log.Errorf("err of doing this:%v", err)
		return "", err
	}

	hashStr := hex.EncodeToString(sha1Enc.Sum(nil))
	return hashStr, nil
}

// 可计算大文件的hash
func CalcStrSha1(fileName string) (string, error) {

	sha1Enc := sha1.New()
	// if _, err := io.Copy(sha1Enc, fileHandle); err != nil {
	// 	log.Errorf("err of doing this:%v", err)
	// 	return "", err
	// }
	sha1Enc.Write([]byte(fileName))

	hashStr := hex.EncodeToString(sha1Enc.Sum(nil))
	return hashStr, nil
}

// 调用前必须Seek到文件头
func CalcFileSha1(fileHandle io.Reader) (string, error) {
	sha1Enc := sha1.New()
	if _, err := io.Copy(sha1Enc, fileHandle); err != nil {
		log.Errorf("err of doing this:%v", err)
		return "", err
	}

	hashStr := hex.EncodeToString(sha1Enc.Sum(nil))
	return hashStr, nil

}
