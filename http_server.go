//go:build amd64

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// request()
	// go build http_server.go

	serve()
}

// http 服务端
func serve() {
	host := flag.String("h", "0.0.0.0", "the serve host")
	port := flag.Int("p", 8081, "the listen port")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("client " + r.RemoteAddr)
		w.Write([]byte("client " + r.RemoteAddr + "\nvisit /req-header show request header.\n"))
		w.Write([]byte(currentTime()))
	})

	http.HandleFunc("/req-header", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(currentTime()))
		r.Header.Set("Servedby", "foo")
		for k, v := range r.Header {
			w.Write([]byte(k + ": " + v[0] + "\n"))
		}
	})

	log.Println(fmt.Sprintf("Server is listening on %s:%d...", *host, *port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil))
}

// 发起 http 请求
func request() {
	reader := strings.NewReader("a=1&b=22&code=6677")
	req, err := http.NewRequest("POST", "http://localhost:9506/user/user/loginByPhoneCode.do", reader)
	if nil != err {
		log.Fatal("create request fail")
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	res, _ := io.ReadAll(resp.Body)

	log.Println(string(res))
}

// 当前时间
func currentTime() string {
	return time.Now().Local().Format("2006-01-02 15:04:05")
}
