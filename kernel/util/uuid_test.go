package util

import "testing"

func TestGetUUID(t *testing.T) {
	uuid := GetUUID()
	t.Log(uuid)
	if len(uuid) != 36 {
		t.Error("uuid length error")
	}
}
