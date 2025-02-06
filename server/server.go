package server

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/Voldemat/go-smtp-mock/emails"
)

type CreateServerRoutinesArgs struct {
    Wg *sync.WaitGroup
    SmtpHost string
    SmtpPort int
    SmtpUser string
    SmtpPassword string
    HttpHost string
    HttpPort int
    QueueSize string
}

func CreateServerRoutines(
    args CreateServerRoutinesArgs,
) *emails.Backend {
    backend, server := emails.CreateSMTPServer(
        args.SmtpHost,
        args.SmtpPort,
        args.SmtpUser,
        args.SmtpPassword,
        args.QueueSize,
    )
    mux := CreateHTTPServer(backend)
	args.Wg.Add(2)
	go func() {
		defer args.Wg.Done()
        http.ListenAndServe(
            args.HttpHost + ":" + strconv.Itoa(args.HttpPort),
            mux,
        )
	}()

	go func() {
		defer args.Wg.Done()
        err := server.ListenAndServe()
        if err != nil {
            log.Fatal(err)
        }
	}()
    return backend
}
