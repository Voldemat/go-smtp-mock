package server

import (
	"encoding/json"
	"github.com/Voldemat/go-smtp-mock/emails"
	"github.com/Voldemat/go-smtp-mock/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateHTTPServer(smtpBackend *emails.Backend) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/clear-queue", func(w http.ResponseWriter, r *http.Request) {
        smtpBackend.ClearQueue()
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		message := []byte("null")
		w.Header().Add("Content-Length", strconv.Itoa(len(message)))
		w.Write(message)
	})
	mux.HandleFunc("/get-last-email", func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		timeout, err := time.ParseDuration(queryValues.Get("timeoutMS") + "ms")
		if err != nil {
			w.WriteHeader(400)
			message := []byte("Invalid timeoutMS param")
			w.Header().Add("Content-Type", "text/plain")
			w.Header().Add("Content-Length", strconv.Itoa(len(message)))
			w.Write(message)
			return
		}
		email := utils.GetValueWithTimeout(smtpBackend.Emails, timeout)
		if email == nil {

			w.WriteHeader(200)
			w.Header().Add("Content-Type", "application/json")
			message := []byte("null")
			w.Header().Add("Content-Length", strconv.Itoa(len(message)))
			w.Write(message)
			return
		}
		message, err := json.Marshal(email)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			message := []byte("Internal server error")
			w.Header().Add("Content-Type", "text/plain")
			w.Header().Add("Content-Length", strconv.Itoa(len(message)))
			w.Write(message)
			return
		}
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Length", strconv.Itoa(len(message)))
		w.Write(message)
	})
	return mux
}
