package httpExample

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func homeHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		_, _ = w.Write([]byte("Home page"))
	} else if r.URL.Path == "/about" {
		_, _ = w.Write([]byte("About page"))
	} else {
		_, _ = w.Write([]byte("Page not found"))
	}
}

func RegisterURL() {
	http.HandleFunc("/", homeHandle)
	fmt.Println("Server listening on port 3000 ...")
	fmt.Println(http.ListenAndServe(":3000", nil))
}

// HandleGetPostExample Method
func HandleGetPostExample() {
	http.HandleFunc("/", returnHtml)
	fmt.Println("Server listening on port 8080 ...")
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func returnHtml(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/get/demo" && r.Method == "GET" {
		getReturnHtml(w, r)
		return
	}
	if r.URL.Path == "/post/demo" && r.Method == "POST" {
		postReturnHtml(w, r)
		return
	}
	if r.URL.Path == "/delete/demo" && r.Method == "DELETE" {
		deleteReturnHtml(w, r)
		return
	}
	if r.URL.Path == "/get/demo" && r.Method == "PUT" {
		putReturnHtml(w, r)
		return
	}
	http.NotFound(w, r)
}

func getReturnHtml(w http.ResponseWriter, _ *http.Request) {
	result := GetAPIExample()
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(result)
	if err != nil {
		return
	}
}
func GetAPIExample() []byte {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos")
	if err != nil {
		fmt.Println("Failed==>", err)
		return []byte{}
	} else {
		body, _ := io.ReadAll(resp.Body)
		return body
	}
}

// HandlePostExample Method
func postReturnHtml(w http.ResponseWriter, _ *http.Request) {
	result := PostAPIExample()
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(result)
	if err != nil {
		return
	}
}

func PostAPIExample() []byte {
	body := []byte(`{
        "userId": 300,
        "id": 300,
        "title": "Demo",
        "completed": false
    }`)
	url := "https://jsonplaceholder.typicode.com/todos"
	req := bytes.NewReader(body)
	_, err := http.Post(url, "application/json", req)
	if err != nil {
		fmt.Println("Failed==>", err)
		return []byte("Failed")
	} else {
		//return body
		return []byte("Post Success")
	}
}

// HandleDeleteExample Method
func deleteReturnHtml(w http.ResponseWriter, _ *http.Request) {
	//Delete Id = 1
	url := "https://jsonplaceholder.typicode.com/todos/1"
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode == http.StatusOK {
		_, _ = w.Write([]byte("Delete Success"))
		_, _ = w.Write(body)
	} else {
		_, _ = w.Write([]byte("Delete Failed"))
	}
}

// HandlePutExample Method
func putReturnHtml(w http.ResponseWriter, _ *http.Request) {
	//Delete Id = 1
	url := "https://jsonplaceholder.typicode.com/todos/1"
	//Body
	body := []byte(`{
        "userId": 1,
        "id": 1,
        "title": "Test",
        "completed": true
    }`)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	body, _ = io.ReadAll(res.Body)
	if res.StatusCode == http.StatusOK {
		//_, _ = w.Write([]byte("Put Success"))
		_, _ = w.Write(body)
	} else {
		_, _ = w.Write([]byte("Put Failed"))
	}
}
