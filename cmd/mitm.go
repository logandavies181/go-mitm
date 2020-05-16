package cmd

import (
	"fmt"
	"crypto/tls"
	"io"
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

func OnConnect(ctx *httpproxy.Context, host string) (
	ConnectAction httpproxy.ConnectAction, newHost string) {
	// Apply "Man in the Middle" to all ssl connections. Never change host.
	return httpproxy.ConnectMitm, host
}

func OnRequest(ctx *httpproxy.Context, req *http.Request) (resp *http.Response) {
	// Log proxying requests.
	log.Printf("INFO: Proxy: %s %s", req.Method, req.URL.String())
	readHeaders(req,"Reqest")
	readBody(req,"Request")
	return
}

func OnResponse(ctx *httpproxy.Context, req *http.Request,resp *http.Response) {
	//readHeaders(resp,"Response")
	readBody(resp,"Response")
}

// Reads the headers and dumps to stdout
func readHeaders(r *http.Request,requestOrResponse string) {
	log.Printf("INFO: %s headers:\n",requestOrResponse)
	for k, v := range r.Header {
		fmt.Println(k, v)
	}
}

// Reads the body (if it exists) and dumps to stdout without consuming it
func readBody(r interface{},requestOrResponse string) error {
	// req is a pointer so we're actually copying the value here
	// so we can avoid closing req.Body and messing up the client request
	switch v := r.(type) {
	case http.Request:
		body = v.Body
	case http.Response:
		body = v.Body
	}
	v := r.(type)
	reqCopy := *v
	body := reqCopy.Body
	bodyDat, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if len(bodyDat) != 0 {
		log.Printf("INFO: %s body:\n",requestOrResponse)
		fmt.Println(string(bodyDat))
	}
	return nil
}

func mitmMain() {
	// Create a new proxy with default certificate pair.
	prx, _ := httpproxy.NewProxy()

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Println(err)
		log.Println("Will use auto generated cert instead")
	} else {
		prx.Ca = cert
	}
	// Set handlers.
	prx.OnError = OnError
	prx.OnAccept = OnAccept
	prx.OnConnect = OnConnect
	prx.OnRequest = OnRequest
	prx.OnResponse = OnResponse

	// Listen...
	http.ListenAndServe(":"+port, prx)
}
