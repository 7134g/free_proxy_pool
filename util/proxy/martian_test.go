package proxy

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestMartian(t *testing.T) {
	SetServeProxyAddress("http://127.0.0.1:7890", "", "")
	OpenCert()
	if err := Martian(); err != nil {
		t.Fatal(err)
	}

}

func TestName(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://www.baidu.com", nil)
	req.WriteProxy(os.Stdout)
	client := http.Client{}
	res, _ := client.Do(req)
	b, _ := io.ReadAll(res.Body)
	fmt.Println(string(b))
}
