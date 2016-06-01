package main

import (
	"io/ioutil"
	"bytes"
        "encoding/json"
	"fmt"
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
	In inside      `json:"0"`
}

type inside struct {
    Id   string      `json:"_id"`
    Version int `json:"version"`
    Body string `json:"body"`
    Compiler string `json:"compiler"`
}

func hello(w http.ResponseWriter, req *http.Request) {

// First, decode the json

    decoder := json.NewDecoder(req.Body)
    var t outside
    err1 := decoder.Decode(&t)
	if (err1 != nil) {print("ERR1");return;}


    data := url.Values{}
    data.Set("version", "2")
    data.Add("body", t.In.Body)


    url := "http://127.0.0.1:"+t.In.Compiler+"/compile?output=json"
  //  fmt.Println("URL:>", url)

	var payload = data.Encode()


    client := &http.Client{}
    r, err2 := http.NewRequest("POST", url, bytes.NewBufferString(payload)) // <-- URL-encoded payload

	if (err2 != nil) {print("ERR2");return;}

//   r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

    resp, err3 := client.Do(r)

	if (err3 != nil) {print("ERR3:");fmt.Println(err3);return;}


	// so we close response
	defer resp.Body.Close()
        contents, err4 := ioutil.ReadAll(resp.Body)
	if (err4 != nil) {print("ERR4");return;}


	var jsonp = fmt.Sprintf("jscallback(%s);", string(contents))
  //      fmt.Printf("%s\n", jsonp)


	// write to file
	err5 := ioutil.WriteFile("/tmp/"+t.In.Id+".js", []byte(jsonp), 0644)
	if (err5 != nil) {print("ERR5");return;}

	if resp.Status != "200 OK" {print("ERR6");return;}

	// dispatch file
	go dispatch(t.In.Id)


}




func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":"+gateport, nil)
}

