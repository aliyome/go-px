package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Body struct {
	Url    string `json:"url"`
	Body   string `json:"body"`
	Method string `json:"method"`
}

func main() {
	http.HandleFunc("/", handleProxy)
	http.ListenAndServe(":54312", nil)
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse body
	defer r.Body.Close()
	var b Body
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bb, err := url.QueryUnescape(b.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b.Body = bb

	fmt.Println(b)

	// forward to target
	if strings.ToLower(b.Method) == "post" {
		resp, err := http.Post(b.Url, "application/json", strings.NewReader(b.Body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)
		w.Write(byteArray)
	}
	if strings.ToLower(b.Method) == "get" {
		resp, err := http.Get(b.Url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)
		w.Write(byteArray)
	}

}
