package main

import (
	"go-smtp-mock/emails"
	"log"
	"net/http"
	"os"
	"sync"
)



func main() {
    queueSize := os.Getenv("SMTP_QUEUE_SIZE")
    if queueSize == "" {
        queueSize = "10"
    }
    backend, server := emails.CreateSMTPServer(
        os.Getenv("SMTP_HOST"),
        os.Getenv("SMTP_PORT"),
        os.Getenv("SMTP_USER"),
        os.Getenv("SMTP_PASSWORD"),
        queueSize,
    )
    mux := CreateHTTPServer(backend)
    wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
        http.ListenAndServe(
            os.Getenv("HTTP_HOST") + ":" + os.Getenv("HTTP_PORT"),
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

	wg.Wait()
}
