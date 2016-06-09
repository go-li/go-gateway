package main

import (
	"os"
	"crypto/sha256" 
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const defaulthost = "127.0.0.1"
const sloppyhost = "yolo.sloppy.zone"


const gateport = "80"

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
	Security bool
	FromJson bool
}
type paste struct {
	Body     string `json:"body"`
}

func perr(e error) {
	fmt.Println(e)
}

func hello(w http.ResponseWriter, req *http.Request) {

	var myhost string

	// first, decide content type, decode json or url encoded, put to t

	var t outside

	if len(req.Header.Get("Content-Type")) == 16 {
		// this is json
		// First, decode the json

		decoder := json.NewDecoder(req.Body)

		err0 := decoder.Decode(&t)
		if err0 != nil {
			print("ERR0\n")
			return
		}

	// Handle pastebins separately
	if (t.In.Compiler == 0) {
		var hsum = sha256.Sum256([]byte(t.In.Body))
		var p = fmt.Sprintf("%x%x%x%x", hsum[0],hsum[1],hsum[2],hsum[3])

		var pasta paste
		pasta.Body = t.In.Body

writr, err5 := os.Create("/tmp/samplecache/"+p+".js")
if err5 != nil {
            fmt.Println(err5);return;
}
		enc := json.NewEncoder(writr)
		fmt.Fprint(writr, "jscallback(")
        if err := enc.Encode(&pasta); err != nil {
            fmt.Println(err);return;
        }
		fmt.Fprint(writr, ");")

		print(p)
		print(":SUM\n")
		return
	}

		// myhost is localhost, because from database
		myhost = defaulthost

		// Port check
		if (t.In.Compiler > 8080) || (t.In.Compiler < 8000) {
			print("ERR1\n")
			return
		}

		// Next, security check
		if (t.In.Compiler != 8080) && osbanned(t.In.Body) {
			print("Route to safe sandbox\n")
			// route to safe sandbox
			t.In.Compiler = 8080
			t.In.Security = true
		}

		t.In.FromJson = true
	} else {
		// this is normal, probably options request
		w.Header().Set("Access-Control-Allow-Origin", "*");
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Allow", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "X-apikey")

		// this is url encoded
		if len(req.Method) != 4 {

			fmt.Fprintf(w, "")
			return
		}


		t.In.Body = req.FormValue("body")
		fmt.Sscanf(req.FormValue("version"), "%d", &t.In.Version)
		fmt.Sscanf(req.FormValue("compiler"), "%d", &t.In.Compiler)

		// Next, security enable
		if osbanned(t.In.Body) {
			t.In.Security = true
		}

		// myhost is based on cloud port
		if t.In.Compiler == 7000 {
			myhost = sloppyhost
			t.In.Compiler = 8078
		} else if t.In.Compiler == 7001 {
			myhost = "play.golang.mx"
			t.In.Compiler = 80
		} else if t.In.Compiler == 7003 {
			myhost = sloppyhost
			t.In.Compiler = 8333
		} else if t.In.Compiler == 7002 {
			myhost = "play.golang.org"
			t.In.Compiler = 80
		} else {
			//unknown cloud port
			print("ERR1\n")
			return
		}

	}



	// Now, build url encoded request for the actual playground docker

	data := url.Values{}
	data.Set("version", "2")
	data.Add("body", t.In.Body)

	if (t.In.Security) {
		data.Set("sec", "y")
	} else {
		data.Set("sec", "n")
	}



	var port = fmt.Sprintf(":%d", t.In.Compiler)
	if (t.In.Compiler == 80) {
		port = "";
	}


	url := "http://" + myhost + port + "/compile?output=json"
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

	if resp.Status != "200 OK" {
		print("ERR6:"+resp.Status+"\n")
		var body = make([]byte,1000);
		resp.Body.Read(body)
		fmt.Println(string(body))
		return
	}

	// so we close response
	defer resp.Body.Close()
	contents, err4 := ioutil.ReadAll(resp.Body)
	if err4 != nil {
		print("ERR4\n")
		return
	}


	if t.In.FromJson {



		var jsonp = fmt.Sprintf("jscallback(%s);", string(contents))
		//      fmt.Printf("%s\n", jsonp)

		// write to file
		err5 := ioutil.WriteFile("/tmp/"+t.In.Id+".js", []byte(jsonp), 0644)
		if err5 != nil {
			print("ERR5\n")
			return
		}

		// dispatch file
		go perr(dispatch(t.In.Id))

		return
	} else {
		fmt.Fprintf(w, "%s", contents)
		return

	}

}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":"+gateport, nil)
}
