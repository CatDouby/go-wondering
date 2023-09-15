package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// 静态资源根目录
var STATIC_PATH string

func main() {
	// request()
	// go build http_server.go

	serve()
}

// http 服务端
func serve() {
	host := flag.String("h", "0.0.0.0", "the serve host")
	port := flag.Int("p", 8081, "the listen port")
	STATIC_PATH := flag.String("fd", "public", "the static file dir, relative/absolute path\nvisit http://host/static/a.jpg => public/a.jgp")
	flag.Parse()

	// 访问到未注册的路由时进入此
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("client " + r.RemoteAddr)
		w.Write([]byte("client " + r.RemoteAddr + "\nvisit /req-header show request header.\n"))
		w.Write([]byte(currentTime()))
	})

	// 输出请求头
	http.HandleFunc("/req-header", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Servedby", "foo")
		// h := map[string]string{}
		for k, v := range r.Header {
			// h[k] = v[0]
			w.Write([]byte(k + ":" + v[0]))
		}
		// rs, _ := json.Marshal(h)
		// w.Write(rs)
	})

	handleStatic(*STATIC_PATH)
	log.Println(fmt.Sprintf("Server is listening on %s:%d...", *host, *port))

	// ListenAndServe 相比 Serve 使用上更简单，而 Serve 可以更自由的控制 listener
	// err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.Serve(l, nil))
}

// handle static resource
func handleStatic(staticPath string) {
	fmt.Println(os.Getwd())
	fmt.Println(staticPath)

	// 设置静态文件目录
	fs := http.FileServer(http.Dir(staticPath))
	// 设置静态文件路由 http://host/static/image/a.jpg => public/image/a.jpg
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 渲染 html
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// 解析静态模板文件
		tpl, err := template.ParseFiles(STATIC_PATH + "/login.html")
		if err != nil {
			log.Fatal(err)
		}

		d := struct {
			Name string
		}{
			Name: r.URL.Query().Get("name"),
		}
		// 数据替换到模板中变量，然后写到输出区域
		tpl.Execute(w, d)
	})
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
