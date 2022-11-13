package main

import (
	"os"
	"path"
	"path/filepath"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"github.com/charlesap/sybil/pkg/lodge"
)

const (
	htmlStart = `<html><body>`
	htmlInfiniteStart = `<!DOCTYPE html>
<html>

<head>
	<title>Your Lodge</title>
	<meta charset="UTF-8" />
	<link rel="stylesheet" href="web/static/styles.css"/>
	<script src="web/app/index.js" defer></script>
</head>

<body>
	<div id="bar">
		<div class="bar_wrapper">

		</div>
	</div>
	<div id="app">
		<div class="whole_wrapper">

		</div>
	</div>

	</script>
`
	htmlEnd = `</body></html>`
	httpPort  = ":80"
)

var (
	flgVerbose             = false
	flgProduction          = false
	flgRedirectHTTPToHTTPS = false
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	time := time.Now().String()
	fmt.Printf(r.URL.String())
	io.WriteString(w, htmlInfiniteStart)
	io.WriteString(w, time)
	io.WriteString(w, htmlEnd)
}

func handleApi(w http.ResponseWriter, r *http.Request) {
	time := time.Now().String()
//	fmt.Printf("API REQUEST")
//	fmt.Printf(r.URL.String())
	io.WriteString(w, `[{ "name" : "bob", "time" : "`)
	io.WriteString(w, time)
	io.WriteString(w, `", "email" : "none", "picture" : "none" },`)
	io.WriteString(w, ` { "name" : "bob", "time" : "`)
	io.WriteString(w, time)
	io.WriteString(w, `", "email" : "none", "picture" : "none" }]`)
}

func handleWebApp(w http.ResponseWriter, r *http.Request) {
	Filename := path.Base(r.URL.String())
	http.ServeFile(w, r, filepath.Join(".", "web/app", Filename))
}

func handleWebStatic(w http.ResponseWriter, r *http.Request) {
	Filename := path.Base(r.URL.String())
	http.ServeFile(w, r, filepath.Join(".", "web/static", Filename))
}

func makeServerFromMux(mux *http.ServeMux) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/", handleApi)
	mux.HandleFunc("/web/static/", handleWebStatic)
	mux.HandleFunc("/web/app/", handleWebApp)
	return makeServerFromMux(mux)

}

func makeHTTPToHTTPSRedirectServer() *http.Server {
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	}
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)
	return makeServerFromMux(mux)
}

func parseFlags() {
	flag.BoolVar(&flgVerbose, "v", false, "if true, more logging")
	flag.BoolVar(&flgProduction, "p", false, "if true, we start HTTPS server")
	flag.BoolVar(&flgRedirectHTTPToHTTPS, "r", false, "if true, we redirect HTTP to HTTPS")
	flag.Parse()
}



func udpResponse(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	time := time.Now().String()
	responseStr := fmt.Sprintf("time received: %v. Your message: %v!", time, string(buf))

	udpServer.WriteTo([]byte(responseStr), addr)
}

func slingUDP() {
	udpServer, err := net.ResolveUDPAddr("udp", ":1053")

	if err != nil {
		println("ResolveUDPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		println("Listen failed:", err.Error())
		os.Exit(1)
	}

	//close the connection
	defer conn.Close()

	_, err = conn.Write([]byte("This is a UDP message"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		os.Exit(1)
	}

	println(string(received))
}

func handleUDP(){

	// listen to incoming udp packets
	udpServer, err := net.ListenPacket("udp", ":1969")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	for {
		buf := make([]byte, 1024)
		_, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go udpResponse(udpServer, addr, buf)
	}

}


func main() {

	parseFlags()
	var m *autocert.Manager

	var httpsSrv *http.Server
	if flgProduction {
		hostPolicy := func(ctx context.Context, host string) error {
			// Note: change to your real host
			allowedHost := "sybil.kuracali.com"
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only '%s' host is allowed, have '%s'", allowedHost, host)
		}

		dataDir := "."
		m = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
		}

		httpsSrv = makeHTTPServer()
		httpsSrv.Addr = ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		go func() {
			fmt.Printf("Starting HTTPS server on %s\n", httpsSrv.Addr)
			err := httpsSrv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
			}
		}()
	}

	var httpSrv *http.Server
	if flgRedirectHTTPToHTTPS {
		httpSrv = makeHTTPToHTTPSRedirectServer()
	} else {
		httpSrv = makeHTTPServer()
	}
	// allow autocert handle Let's Encrypt callbacks over http
	if m != nil {
		httpSrv.Handler = m.HTTPHandler(httpSrv.Handler)
	}


        go handleUDP()

	httpSrv.Addr = httpPort
	fmt.Printf("Starting HTTP server on %s\n", httpPort)
	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}



	baseName := filepath.Base(os.Args[0])

	lodge.Emit(baseName)
}
