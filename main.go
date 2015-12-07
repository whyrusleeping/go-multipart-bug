package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"
)

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header["Content-Type"][0]
		parts := strings.Split(ct, "\"")
		fmt.Println(parts[1])
		mpr := multipart.NewReader(r.Body, parts[1])
		part, err := mpr.NextPart()
		if err != nil {
			fmt.Println(part)
			fmt.Println("ERROR: ", err)
			return
		}

		fmt.Println(part.FileName())
		out, err := ioutil.ReadAll(part)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d bytes\n", len(out))
	})

	panic(http.ListenAndServe(":9999", nil))
}

func main() {
	// 4031 triggers the bug, any other number does not.
	data := make([]byte, 4031)

	filename := "badstuff"
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	// write the boundary and headers
	header := make(textproto.MIMEHeader)
	filename = url.QueryEscape(filename)
	contentDisposition := fmt.Sprintf("form-data; name=\"file\"; filename=\"%s\"", filename)
	header.Set("Content-Disposition", contentDisposition)
	header.Set("Content-Type", "application/octet-stream")

	pw, err := w.CreatePart(header)
	if err != nil {
		panic(err)
	}

	pw.Write(data)
	w.Close()

	req, err := http.NewRequest("POST", "http://localhost:9999/asd", buf)
	if err != nil {
		panic(err)
	}
	req.Header["Content-Type"] = []string{fmt.Sprintf("multipart/form-data; boundary=\"%s\"", w.Boundary())}

	go server()

	time.Sleep(time.Millisecond * 100)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
