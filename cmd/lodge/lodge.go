package main

import (
	"os"
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
	htmlEnd = `</body></html>`
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
	htmlCSS = `.bar_wrapper {
  background: rgba(0, 0, 0, .1);
  width: 100%;
  min-height: 20px;
  padding:5%;
}
.whole_wrapper {
  background: rgba(0, 0, 0, .1);
  width: 100%;
  min-height: 100%;
  padding:5%;
}
.whole_wrapper .each_card {
  width: 50%;
  align-items: center;
  text-align: center;
  display: flex;
  padding: 10px;
  background: white;
  margin:5% 25%;
  box-shadow: 0 2px 5px 0 rgba(0, 0, 0, 0.16), 0 2px 10px 0 rgba(0, 0, 0, 0.12);
}
.whole_wrapper .each_card .image_container {
  text-align: left;
}
.whole_wrapper .each_card .image_container img {
  width: 50%;
  border-radius: 10px;
}
.whole_wrapper .each_card .right_contents_container {
  display: flex;
  flex-direction: column;
}
.whole_wrapper .each_card .right_contents_container .name_field {
  font-size: 22px;
  font-weight: 900;
  line-height: 30px;
}
.whole_wrapper .each_card .right_contents_container .email_field {
  font-size: 22px;
  line-height: 30px;
}
`
	htmlJS = `let page = 1;
const last_page = 10;
const pixel_offset = 200;
const throttle = (callBack, delay) => {
  let withinInterval;
  return function() {
    const args = arguments;
    const context = this;
    if (!withinInterval) {
      callBack.call(context, args);
      withinInterval = true;
      setTimeout(() => (withinInterval = false), delay);
    }
  };
};

const httpRequestWrapper = (method, URL) => {
  return new Promise((resolve, reject) => {
    const xhr_obj = new XMLHttpRequest();
    xhr_obj.responseType = "json";
    xhr_obj.open(method, URL);
    xhr_obj.onload = () => {
      const data = xhr_obj.response;
      resolve(data);
    };
    xhr_obj.onerror = () => {
      reject("failed");
    };
    xhr_obj.send();
  });
};

const getData = async (page_no = 1) => {
  const data = await httpRequestWrapper(
    "GET",
    "https://sybil.kuracali.com/api/?page=${page_no}&results=10"
  );

  const {results} = data;
  populateUI(results);
};

let handleLoad;

let trottleHandler = () =>{throttle(handleLoad.call(this), 1000)};

document.addEventListener("DOMContentLoaded", () => {
  getData(1);
  window.addEventListener("scroll", trottleHandler);
});

handleLoad =  () => {
  if((window.innerHeight + window.scrollY) >= document.body.offsetHeight - pixel_offset){
    page = page+1;
    if(page<=last_page){
      window.removeEventListener('scroll',trottleHandler)
      getData(page)
      .then((res)=>{
        window.addEventListener('scroll',trottleHandler)
      })
    }
  }
}



const populateUI = data => {
  const barcontainer = document.querySelector('.bar_wrapper');
  container.innerHTML += "Loading...";
  const container = document.querySelector('.whole_wrapper');
  data && 
  data.length && 
  data
  .map((each,index)=>{
    const {name,email,picture} = each;
    const {first} = name;
    const {large} = picture;
    container.innerHTML += '    <div class="each_card">' +
                           '       <div class="image_container">' +
                           '         <img src="${large}" alt="" />' +
                           '       </div>' +
                           '       <div class="right_contents_container">' +
                           '         <div class="name_field">${first}</div>' +
                           '         <div class="email_filed">${email}</div>' +
                           '       </div>' +
                           '    </div>'
  })

}
`
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
	io.WriteString(w, time)
}

func handleWebAppIndexJS(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, htmlJS)
}

func handleWebStaticStylesCSS(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, htmlCSS)
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
	mux.HandleFunc("/api", handleApi)
	mux.HandleFunc("/web/static/styles.css", handleWebStaticStylesCSS)
	mux.HandleFunc("/web/app/index.js", handleWebAppIndexJS)
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
