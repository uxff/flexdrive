package utils

import (
	"encoding/json"
	"testing"
	"time"
)

type TTT struct {
	CreateAt JsonTime
	StdTime time.Time
}

func TestJsonTime_IsEmptyTime(t *testing.T) {
	a := &TTT{CreateAt: JsonTime{}}

	b, err := json.Marshal(a)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%s", b)
}
