package filehash

import (
	"testing"
)

func TestSha1Enc(t *testing.T) {

	file := "/usr/local/gopath/src/github.com/uxff/flexdrive/main"
	hash, err := CalcSha1(file)

	if err != nil {
		t.Errorf("hash, err := %v, %v", hash, err)
	}
}
