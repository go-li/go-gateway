package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const gateport = "8333"

// UrlEncoded encodes a string like Javascript's encodeURIComponent()
func UrlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

type outside struct {
	In inside `json:"0"`
}

type inside struct {
	Id       string `json:"_id"`
	Version  int    `json:"version"`
	Body     string `json:"body"`
	Compiler int    `json:"compiler"`
}

func perr(e error) {
	fmt.Println(e)
}

func hello(w http.ResponseWriter, req *http.Request) {

	// First, decode the json

	decoder := json.NewDecoder(req.Body)
	var t outside
	err1 := decoder.Decode(&t)
	if err1 != nil {
		print("ERR0\n")
		return
	}

	// Next, security check
	if (t.In.Compiler != 8080) && osbanned(t.In.Body) {
		print("ERR1 security\n")
		return
	}

	// Now, build url encoded request for the actual playground docker

	data := url.Values{}
	data.Set("version", "2")
	data.Add("body", t.In.Body)

	if (t.In.Compiler > 8080) || (t.In.Compiler < 8000) {
		print("ERR1x\n")
		return
	}

	var port = fmt.Sprintf("%d", t.In.Compiler)

	url := "http://127.0.0.1:" + port + "/compile?output=json"
	//  fmt.Println("URL:>", url)

	var payload = data.Encode()

	client := &http.Client{}
	r, err2 := http.NewRequest("POST", url, bytes.NewBufferString(payload)) // <-- URL-encoded payload

	if err2 != nil {
		print("ERR2\n")
		return
	}

	//   r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	resp, err3 := client.Do(r)

	if err3 != nil {
		print("ERR3:")
		fmt.Println(err3)
		return
	}

	// so we close response
	defer resp.Body.Close()
	contents, err4 := ioutil.ReadAll(resp.Body)
	if err4 != nil {
		print("ERR4\n")
		return
	}

	var jsonp = fmt.Sprintf("jscallback(%s);", string(contents))
	//      fmt.Printf("%s\n", jsonp)

	// write to file
	err5 := ioutil.WriteFile("/tmp/"+t.In.Id+".js", []byte(jsonp), 0644)
	if err5 != nil {
		print("ERR5\n")
		return
	}

	if resp.Status != "200 OK" {
		print("ERR6\n")
		return
	}

	// dispatch file
	go perr(dispatch(t.In.Id))

}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":"+gateport, nil)
}
