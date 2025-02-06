package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/Voldemat/go-smtp-mock/server"
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
	smtpPortString := os.Getenv("SMTP_PORT")
	if smtpPortString == "" {
		log.Fatal("SMTP_PORT is not defined")
	}
	smtpPort, err := strconv.ParseInt(smtpPortString, 10, 64)
	if err != nil {
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
	httpPortString := os.Getenv("SMTP_HTTP_PORT")
	if httpPortString == "" {
		log.Fatal("SMTP_HTTP_PORT is not defined")
	}
	httpPort, err := strconv.ParseInt(httpPortString, 10, 64)
	if err != nil {
		log.Fatal("SMTP_HTTP_PORT is not defined")
	}
	var wg sync.WaitGroup
	wg.Add(2)

	server.CreateServerRoutines(server.CreateServerRoutinesArgs{
		Wg:           &wg,
		SmtpHost:     smtpHost,
		SmtpPort:     int(smtpPort),
		SmtpUser:     smtpUser,
		SmtpPassword: smtpPassword,
		HttpHost:     os.Getenv("HTTP_HOST"),
		HttpPort:     int(httpPort),
		QueueSize:    queueSize,
	},
	)
	wg.Wait()
}
