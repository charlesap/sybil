package main

import (
	"path"
	"path/filepath"
	"context"
	"crypto/tls"
	"strconv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"


	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"


	"golang.org/x/text/language"
	"golang.org/x/crypto/acme/autocert"

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
			<div class="bar_card">

			</div>
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
	bundle                 *i18n.Bundle
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
//	time := time.Now().String()
//	fmt.Printf(r.URL.String())
	io.WriteString(w, htmlInfiniteStart)
//	io.WriteString(w, time)
	io.WriteString(w, htmlEnd)
}

func handleApi(w http.ResponseWriter, r *http.Request) {

	name:="bob"

//	lang := r.FormValue("lang")
	lang := "es"
	accept := r.Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(bundle, lang, accept)

	helloPerson := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "HelloPerson",
			Other: "Hello {{.Name}}",
		},
		TemplateData: map[string]string{
			"Name": name,
		},
	})


//	al:=r.Header.Get("Accept-Language")
//	m := language.NewMatcher([]language.Tag{language.English, language.French})
//	desired, _, _ := language.ParseAcceptLanguage(al)
//	tag, i, conf := m.Match(desired...)
//	fmt.Println("case B", tag, i, conf) // fr-u-rg-dezzzz instead of en-u-rg-gbzzzz

	pg := r.URL.Query().Get("page")
	p, err1 := strconv.Atoi(pg)
	if err1 != nil { p = 1 }
	rs := r.URL.Query().Get("results")
	n, err2 := strconv.Atoi(rs)
	if err2 != nil { n = 1 }
//	time := time.Now().String()
//	fmt.Printf("API REQUEST")
	fmt.Printf("%s %s %d %s %d \n",r.URL.String(),pg,p,rs,n)
	io.WriteString(w, `{"results":[`)
	for i:=0;i<n;i++ {
		s:=fmt.Sprintf("%s %d %s %s %s%d%s",`{ "name" : "bob`,((p-1)*n)+i,`", "time" : "`,helloPerson, `", "email" : "bob@bob.com", "picture" : "img/`,((p-1)*n)+i,`.jpg" }`)
		io.WriteString(w, s)
		if i < (n-1) {
			io.WriteString(w, `,`)
		}
	}
	io.WriteString(w, `],"info":{"seed":"X","results":`)
	io.WriteString(w, rs)
	io.WriteString(w, `,"page":`)
	io.WriteString(w, pg)
	io.WriteString(w, `,"version":"0.1"}}`)
}

func handleImg(w http.ResponseWriter, r *http.Request) {
	Filename := path.Base(r.URL.String())
	http.ServeFile(w, r, filepath.Join(".", "examples", Filename))
}

func handleWebApp(w http.ResponseWriter, r *http.Request) {
	Filename := path.Base(r.URL.String())
	http.ServeFile(w, r, filepath.Join(".", "web/app", Filename))
}

func handleWebStatic(w http.ResponseWriter, r *http.Request) {
	Filename := path.Base(r.URL.String())
	http.ServeFile(w, r, filepath.Join(".", "web/static", Filename))
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/favicon.ico")
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
	mux.HandleFunc("/img/", handleImg)
	mux.HandleFunc("/web/static/", handleWebStatic)
	mux.HandleFunc("/web/app/", handleWebApp)
	mux.HandleFunc("/favicon.ico", handleFavicon)
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



func Webmain() {

	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// No need to load active.en.toml since we are providing default translations.
	bundle.MustLoadMessageFile("web/static/active.en.toml")
	bundle.MustLoadMessageFile("web/static/active.es.toml")

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


	httpSrv.Addr = httpPort
	fmt.Printf("Starting HTTP server on %s\n", httpPort)
	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}



}

