package crawler

import (
	"testing"
)

func TestRunTestProxy(t *testing.T) {
	if runTestProxy("http://www.baidu.com", "http://127.0.0.1:7890") {
		t.Log("success")
	} else {
		t.Error("fail")
	}
}
