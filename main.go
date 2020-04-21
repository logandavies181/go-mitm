package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-httpproxy/httpproxy"
)

func OnError(ctx *httpproxy.Context, where string,
	err *httpproxy.Error, opErr error) {
	// Log errors.
	log.Printf("ERR: %s: %s [%s]", where, err, opErr)
}

func OnAccept(ctx *httpproxy.Context, w http.ResponseWriter,
	r *http.Request) bool {
	// Handle local request has path "/info"
	if r.Method == "GET" && !r.URL.IsAbs() && r.URL.Path == "/info" {
		w.Write([]byte("This is go-httpproxy."))
		return true
	}
	return false
}

func OnAuth(ctx *httpproxy.Context, authType string, user string, pass string) bool {
	// Auth test user.
	return true
}

func OnConnect(ctx *httpproxy.Context, host string) (
	ConnectAction httpproxy.ConnectAction, newHost string) {
	// Apply "Man in the Middle" to all ssl connections. Never change host.
	return httpproxy.ConnectMitm, host
}

func OnRequest(ctx *httpproxy.Context, req *http.Request) (
	resp *http.Response) {
	// Log proxying requests.
	log.Printf("INFO: Proxy: %s %s", req.Method, req.URL.String())
	log.Printf("DEEB: headers:")
	for k, v := range req.Header {
		log.Println(k, v)
	}

	if req.Body != nil {
		reqCopy := *req
		reqBody := reqCopy.Body
		/*reqBody, err := req.GetBody()
		if err != nil {
			log.Println("getbody err")
		}*/
		reqBodyDat, err := ioutil.ReadAll(reqBody)
		//_, err = ioutil.ReadAll(reqBody)
		if err != nil {
			log.Println("err2")
			log.Println(err)
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(reqBodyDat))
		log.Printf("DEEB: %s", string(reqBodyDat))
	}
	return
}

func OnResponse(ctx *httpproxy.Context, req *http.Request,
	resp *http.Response) {
	// Add header "Via: go-httpproxy".
	resp.Header.Add("Via", "go-httpproxy")
}

func main() {
	// Create a new proxy with default certificate pair.
	prx, _ := httpproxy.NewProxy()

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Println(err)
	}
	prx.Ca = cert
	// Set handlers.
	prx.OnError = OnError
	prx.OnAccept = OnAccept
	//prx.OnAuth = OnAuth
	prx.OnConnect = OnConnect
	prx.OnRequest = OnRequest
	prx.OnResponse = OnResponse

	// Listen...
	http.ListenAndServe(":8080", prx)
}
