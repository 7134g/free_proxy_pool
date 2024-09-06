package crawler

import (
	"fmt"
	"sort"
	"testing"
)

func TestRunTestProxy(t *testing.T) {
	if runTestProxy("http://www.baidu.com", "http://127.0.0.1:7890") {
		t.Log("success")
	} else {
		t.Error("fail")
	}
}

func TestName(t *testing.T) {
	x := []int{1, 4, 6, 2, 7, 3}

	sort.Slice(x, func(i, j int) bool {
		return x[i] > x[j]
	})
	fmt.Println(x)
}
