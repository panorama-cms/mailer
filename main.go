package mailer

import (
	"errors"
	"gopkg.in/gomail.v2"
	"strings"
)

// Smtp represents the smtp configuration
type Smtp struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Address represents an email address.
type Address struct {
	Name    string
	Address string
}

// Mail represents the Mail object.
type Mail struct {
	From        Address
	To          []Address
	Bcc         []Address
	Cc          []Address
	ReplyTo     []Address
	Body        string
	Subject     string
	Attachments []string
	Smtp        Smtp
}

// New returns a new Mail object that needs to be set up.
func New() *Mail {
	return &Mail{}
}

func (m *Mail) AddAddress(address string, name string) {
	m.To = append(m.To, Address{
		Name:    name,
		Address: address,
	})
}

func (m *Mail) AddBCC(address string, name string) {
	m.Bcc = append(m.Bcc, Address{
		Name:    name,
		Address: address,
	})
}

func (m *Mail) AddCC(address string, name string) {
	m.Cc = append(m.Cc, Address{
		Name:    name,
		Address: address,
	})
}

func (m *Mail) AddAttachment(path string) {
	m.Attachments = append(m.Attachments, path)
}

func (m *Mail) AddStringAttachment(data string, filename string) {
	m.Attachments = append(m.Attachments, data+";filename="+filename)
}

func (m *Mail) ClearAddresses() {
	m.To = []Address{}
}

func (m *Mail) ClearCCs() {
	m.Cc = []Address{}
}

func (m *Mail) ClearBCCs() {
	m.Bcc = []Address{}
}

func (m *Mail) ClearAllRecipients() {
	m.ClearAddresses()
	m.ClearCCs()
	m.ClearBCCs()
}

func (m *Mail) ClearAttachments() {
	m.Attachments = []string{}
}

func (m *Mail) ClearReplyTos() {
	m.ReplyTo = []Address{}
}

func (m *Mail) convertAddress(address Address) string {
	return address.Name + " <" + address.Address + ">"
}

func (m *Mail) convertAddresses(addresses []Address) []string {
	var converted []string
	for _, address := range addresses {
		converted = append(converted, m.convertAddress(address))
	}
	return converted
}

func (m *Mail) SetFrom(address string, name string) {
	m.From = Address{
		Name:    name,
		Address: address,
	}
}

func (m *Mail) SetSubject(subject string) {
	m.Subject = subject
}

func (m *Mail) SetBody(body string) {
	m.Body = body
}

func (m *Mail) SetSmtp(host string, port int, username string, password string) {
	m.Smtp = Smtp{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (m *Mail) AddReplyTo(address string, name string) {
	m.ReplyTo = append(m.ReplyTo, Address{
		Name:    name,
		Address: address,
	})
}

func (m *Mail) validate() error {
	if m.From.Address == "" {
		return errors.New("from is not set")
	}
	if m.Subject == "" {
		return errors.New("subject is not set")
	}
	if m.Body == "" {
		return errors.New("body is not set")
	}
	if m.Smtp == (Smtp{}) {
		return errors.New("smtp is not set")
	}
	if m.Smtp.Host == "" {
		return errors.New("host is not set")
	}
	if m.Smtp.Port == 0 {
		return errors.New("port is not set")
	}
	if m.Smtp.Username == "" {
		return errors.New("username is not set")
	}
	if m.Smtp.Password == "" {
		return errors.New("password is not set")
	}

	return nil
}

func (m *Mail) Send() error {
	if err := m.validate(); err != nil {
		return err
	}
	msg := gomail.NewMessage()

	from := m.convertAddress(m.From)
	msg.SetHeader("From", from)

	to := strings.Join(m.convertAddresses(m.To), ", ")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", m.Body)
	for _, attachment := range m.Attachments {
		msg.Attach(attachment)
	}

	d := gomail.NewDialer(m.Smtp.Host, m.Smtp.Port, m.Smtp.Username, m.Smtp.Password)

	// Send the email
	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
