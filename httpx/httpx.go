package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func init() {
	log.SetPrefix("[PanGO] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

//Post http get method
func Get(url string, params map[string]string, headers map[string]string) (*http.Response, error) {
	//new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail ")
	}
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	client := &http.Client{}
	log.Printf("Go GET URL : %s \n", req.URL.String())
	return client.Do(req)
}

//DefaultGet
func DefaultGet(url string, params map[string]string) (*http.Response, error) {
	return Get(url, params, nil)
}

//Post http post method
func Post(url string, body map[string]string, params map[string]string, headers map[string]string) (*http.Response, error) {
	//add post body
	var bodyJson []byte
	var req *http.Request
	if body != nil {
		var err error
		bodyJson, err = json.Marshal(body)
		if err != nil {
			log.Println(err)
			return nil, errors.New("http post body to json failed")
		}
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail: %v \n")
	}
	req.Header.Set("Content-type", "application/json")
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	client := &http.Client{}
	log.Printf("Go POST URL : %s \n", req.URL.String())
	return client.Do(req)
}

func DefaultPost(url string, body map[string]string) (*http.Response, error) {
	return Post(url, body, nil, nil)
}

//Parse parse http response
func Parse(resp *http.Response) (interface{}, error) {
	defer resp.Body.Close()
	log.Printf("HTTP code: %d \n", resp.StatusCode)
	byArr, err := ioutil.ReadAll(resp.Body)
	if bytes.ContainsAny(byArr, ">") {
		return string(byArr), nil
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, byArr, "   ", "\t")
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(prettyJSON.Bytes()), err
}

func PrintResult(resp *http.Response, err error, t *testing.T) {
	if err != nil {
		t.Errorf("err = %s \n", err.Error())
		return
	}
	parse, err := Parse(resp)
	if err != nil {
		t.Errorf("err = %s \n", err.Error())
		return
	}
	fmt.Printf("DATA=%v\n", parse)
}
