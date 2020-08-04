package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

// Get env var or default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(getEnv("RF_FORWARD_URL", "https://api.ipgeolocation.io"))

	debug, _ := strconv.ParseBool(getEnv("RF_DEBUG", "false"))

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host
	// We need to clear the remote addr field so the getip endpoint works properly
	originalRemoteAddr := req.RemoteAddr
	req.RemoteAddr = ""

	if debug {
		log.Println(":::START:Forwarding Request:::")
		log.Printf("URI: %s\n", req.URL)
		log.Printf("Host: %s\n", req.URL.Host)
		log.Printf("Path: %s\n", req.URL.Path)
		log.Printf("URI: %s\n", req.URL.RequestURI())
		log.Printf("Body: %s\n", req.Body)
		log.Printf("originalRemoteAddr: %s\n", originalRemoteAddr)
		log.Printf("FullRequest: %s\n", req)
		log.Println(":::END:Forwarding Request:::")
	}

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

/*
	Entry
*/

func main() {

	port := getEnv("RF_PORT", "8080")

	// start server
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
