package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const tpl = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<title>Web Font Viewer</title>

		<style type="text/css">
			{{range .}}
				@font-face {
					font-family: {{.Family}};
					src: url('fonts/{{.Url}}') format('{{.Format}}') /* chrome、firefox */
				}
			{{end}}
		</style>
	</head>
	<body>
	<h1> 字体预览效果 </h1>
	<table border = "1px">
		<tr>
			<th width = "80px">名称</th>
			<th width = "100px">自定义文字</th>
			<th width = "600px">默认文字</th>
			<th width = "100px">格式</th>
		</tr>
		{{range .}}
			<tr>
				<td> {{.Family}} </td>
				<td style="font-family: {{.Family}}; font-size:24px">
					{{.PreviewText}}
				</td>
				<td style="font-family: {{.Family}}; font-size:24px">
					The quick brown fox jumps over a lazy dog. 敏捷的棕色狐狸跳过了一只懒惰的狗
				</td>
				<td> {{.Format}}</td>
			</tr>
		{{end}}
	<table>
	</body>
	</html>
`

type Font struct {
	Family      string
	Url         string
	Format      string
	PreviewText string
}

var port uint = 80
var h bool

func init() {
	flag.BoolVar(&h, "help", false, "help")
	flag.UintVar(&port, "p", 8080, "listen port")
}

func main() {
	flag.Parse()

	fs := http.FileServer(http.Dir("."))
	http.Handle("/fonts/", http.StripPrefix("/fonts", fs))
	http.HandleFunc("/", index)
	fmt.Printf("Navigate to http://localhost:%d  to view font in current directory\r\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")

	var fonts []Font
	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return err
		}
		name := f.Name()
		if strings.HasSuffix(path, ".ttf") || strings.HasSuffix(path, ".TTF") {
			fonts = append(fonts, Font{name[:len(name)-4], name, "truetype", text})
		} else if strings.HasSuffix(path, ".woff2") || strings.HasSuffix(path, ".WOFF2") {
			fonts = append(fonts, Font{name[:len(name)-6], name, "woff2", text})
		} else if strings.HasSuffix(path, ".woff") || strings.HasSuffix(path, ".WOFF") {
			fonts = append(fonts, Font{name[:len(name)-5], name, "woff", text})
		} else if strings.HasSuffix(path, ".otf") || strings.HasSuffix(path, ".OTF") {
			fonts = append(fonts, Font{name[:len(name)-4], name, "opentype", text})
		}

		return err
	})

	t, err := template.New("font").Parse(tpl)

	if err != nil {
		fmt.Println(err)
	}

	t.Execute(w, fonts)
}
