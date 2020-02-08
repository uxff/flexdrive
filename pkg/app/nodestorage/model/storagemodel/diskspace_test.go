package storagemodel

import "testing"

func TestDiskspace(t *testing.T) {
	ds := DiskStatus("/")
	t.Errorf("ds=+v", ds)
}
