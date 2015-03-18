// Package negroni-rollbar is a middleware for Negroni that reports panics to rollbar.com.
//
//  package main
//
//  import (
//    "github.com/codegangsta/negroni"
//    "github.com/jfbus/negroni-rollbar"
//  )
//
//
//  func main() {
//    n := negroni.Classic()
//    n.Use(rollbar.Report(rollbar.Config{Token: ROLLBAR_TOKEN}))
//
//    m := pat.New()
//    m.Get("/panic", func() {
//      panic("an error occured")
//    })
//    n.UseHandler(m)
//    n.Run()
//  }
package rollbar

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	rb "github.com/stvp/rollbar"
)

type Config struct {
	Token string
}

type report struct{}

// Report returns a middleware that recovers from any panics, sends the error to rollbar and writes a HTTP 500 response.
func Report(cfg Config) negroni.Handler {
	rb.Token = cfg.Token
	rb.Environment = "production"
	return &report{}
}

func (m *report) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {

			stack := rb.BuildStack(3)
			if e, ok := err.(error); ok {
				rb.RequestErrorWithStack(rb.CRIT, r, e, stack)
			} else {
				rb.RequestErrorWithStack(rb.CRIT, r, fmt.Errorf("%s", err), stack)
			}
			var str string
			for _, f := range stack {
				str += fmt.Sprintf("File \"%s\" line %d in %s\n", f.Filename, f.Line, f.Method)
			}

			log.Printf("PANIC: %s\n%s", err, str)

			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("500 Internal Server Error"))
		}
	}()

	next(rw, r)
}
