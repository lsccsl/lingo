package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ServerFromHttp struct {
	SrvID int64
	IP string
	Port int
}

func httpAddDial(sh*ServerFromHttp) {
	jBin, _ := json.Marshal(sh)
	http.Post(Global_testCfg.httpAddr, "application/text", strings.NewReader(string(jBin)))
	//lin_common.LogDebug(err, res)
	//http.DefaultClient.CloseIdleConnections()
}

func testhttp() {
	data := make(url.Values)
	data["key"] = []string{"val"}
	data["key1"] = []string{"val2"}

	tj := &ServerFromHttp{
		SrvID:123,
		IP :"10.0.0.1",
		Port : 123,
	}
	jBin, _ := json.Marshal(tj)
	res, err := http.Post("http://10.0.14.48:8802/addserver", "application/text", strings.NewReader(string(jBin)))
	//res, err := http.PostForm("http://10.0.14.48:8802/post", data)
	fmt.Println(res)
	fmt.Println(err, res.ContentLength, res.Body)
	bin := make([]byte, res.ContentLength, res.ContentLength)
	n, err := res.Body.Read(bin)
	fmt.Println(n, err, string(bin))

	res, err = http.PostForm("http://10.0.14.48:8802/addserver", data)
	fmt.Println(res)
	fmt.Println(err, res.ContentLength, res.Body)
	bin = make([]byte, res.ContentLength, res.ContentLength)
	n, err = res.Body.Read(bin)
	fmt.Println(n, err, string(bin))
}
