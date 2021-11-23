package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/go-apache-test/lib/test"
)

func Start(port string) *http.Server {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: ":" + port, Handler: mux}

	mux.HandleFunc("/", handle)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return srv
}

func Stop(srv *http.Server) {
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

func handle(res http.ResponseWriter, req *http.Request) {
	xForwardedHost := req.URL.Host
	if len(req.Header["X-Forwarded-Host"]) > 0 {
		xForwardedHost = req.Header["X-Forwarded-Host"][0]
	}

	m1 := regexp.MustCompile(`(:[0-9]{0,4})`)
	pathname := req.URL.Path
	if req.URL.RawQuery != "" {
		pathname += "?" + req.URL.RawQuery
	}
	e := test.Expectations{
		Backend:  m1.ReplaceAllString(req.Host, ""),
		Host:     xForwardedHost,
		Pathname: pathname,
	}
	jsonByte, err := json.Marshal(e)

	if err != nil {
		res.Header().Set("Content-Type", "text/html")
		res.Write([]byte(fmt.Sprint(err)))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonByte)
}
