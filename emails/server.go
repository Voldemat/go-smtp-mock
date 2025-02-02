package emails

import (
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type Email struct {
	From string   `json:"from"`
	To   []string `json:"to"`
	Body []byte   `json:"body"`
}

type Backend struct {
	username string
	password string
	Emails   chan Email
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{bkd: bkd}, nil
}

func (bkd *Backend) AddEmail(email Email) {
	if len(bkd.Emails) == cap(bkd.Emails) {
		<-bkd.Emails
	}
	bkd.Emails <- email
}

type Session struct {
	bkd   *Backend
	email Email
}

func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

func (s *Session) Auth(mech string) (sasl.Server, error) {
	server := sasl.NewPlainServer(
		func(identity string, username string, password string) error {
			if username != s.bkd.username || password != s.bkd.password {
				return errors.New("Invalid username or password")
			}
			return nil
		},
	)
	return server, nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.email.From = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.email.To = append(s.email.To, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	s.email.Body = b
	return nil
}

func (s *Session) Reset() {
	s.bkd.AddEmail(s.email)
	s.email = Email{}
}

func (s *Session) Logout() error {
	return nil
}

func CreateSMTPServer(
	host string,
	port string,
	username string,
	password string,
	queueSizeString string,
) (*Backend, *smtp.Server) {
    queueSize, err := strconv.ParseInt(queueSizeString, 10, 64)
    if err != nil {
        panic(err)
    }
	be := &Backend{
		username: username,
		password: password,
		Emails:   make(chan Email, queueSize),
	}
	s := smtp.NewServer(be)
	s.Addr = host + ":" + port
	s.Domain = host
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	return be, s
}
