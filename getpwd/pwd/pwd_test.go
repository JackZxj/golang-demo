package pwd

import "testing"

func TestRun(t *testing.T) {
	str := Run()
	if str == "" {
		t.Fatalf("empty")
	} else {
		t.Log(str)
	}
}
