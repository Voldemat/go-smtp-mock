package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/Voldemat/go-smtp-mock/emails"
)

func CreateServerRoutines(
    wg *sync.WaitGroup,
    smtpHost string,
    smtpUser string,
    smtpPort string,
    smtpPassword string,
    queueSize string,
    httpHost string,
    httpPort string,
) *emails.Backend {
    backend, server := emails.CreateSMTPServer(
        smtpHost,
        smtpPort,
        smtpUser,
        smtpPassword,
        queueSize,
    )
    mux := CreateHTTPServer(backend)

	wg.Add(2)

	go func() {
		defer wg.Done()
        http.ListenAndServe(
            httpHost + ":" + httpPort,
            mux,
        )
	}()

	go func() {
		defer wg.Done()
        err := server.ListenAndServe()
        if err != nil {
            log.Fatal(err)
        }
	}()
    return backend
}
