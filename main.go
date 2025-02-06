package main

import (
	"github.com/Voldemat/go-smtp-mock/emails"
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
    smtpHost := os.Getenv("SMTP_HOST")
    if smtpHost == "" {
        log.Fatal("SMTP_HOST is not defined")
    }
    smtpPort := os.Getenv("SMTP_PORT")
    if smtpPort == "" {
        log.Fatal("SMTP_PORT is not defined")
    }
    smtpUser := os.Getenv("SMTP_USER")
    if smtpUser == "" {
        log.Fatal("SMTP_USER is not defined")
    }
    smtpPassword := os.Getenv("SMTP_PASSWORD")
    if smtpPassword == "" {
        log.Fatal("SMTP_PASSWORD is not defined")
    }
    backend, server := emails.CreateSMTPServer(
        smtpHost,
        smtpPort,
        smtpUser,
        smtpPassword,
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
