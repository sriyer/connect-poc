package main

import (
	"os"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
	"bytes"
)

const targetAddress = "TARGET_ADDRESS"
const typeError = "[ERROR] - "

func main() {
	if os.Getenv(targetAddress) == "" {
		log.Fatal("target address not set, please set ", targetAddress, " to the proxy target.")
	}
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":80", nil)
}

func defaultHandler(writer http.ResponseWriter, request *http.Request) {
	response, err := proxy(request)
	if err != nil {
		if response == nil {
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			writer.WriteHeader(response.StatusCode)
		}
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(response.StatusCode)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write(body)
}

func proxy (r *http.Request) (*http.Response, error) {
	target := os.Getenv(targetAddress)
	client := http.Client{}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://%s%s?%s", target, r.URL.Path, r.URL.RawQuery)
	log.Println("calling ", r.Method, url)
	request, err := http.NewRequest(r.Method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println(typeError, err)
		return nil, err
	}
	return client.Do(request)
}


