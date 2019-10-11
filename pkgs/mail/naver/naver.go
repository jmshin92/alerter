package naver

import (
	"log"
	"net/smtp"

	"github.com/jmshin92/alerter/pkgs/mail"
	"github.com/sirupsen/logrus"
)

const (
	Addr = "smtp.naver.com:587"
	Host = "smtp.naver.com"
)

type NaverMail struct {
	*mail.MailConf

	boundary string
}

func NewNaverMail(c *mail.MailConf) (*NaverMail, error) {
	m := &NaverMail{
		MailConf: c,
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

func (this *NaverMail) Send() error {
	msg := this.Build()
	logrus.Info(msg)
	err := smtp.SendMail(Addr,
		smtp.PlainAuth("", this.From, this.Password, Host),
		this.From, []string{this.To}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}
	return nil
}

func (this *NaverMail) SetFrom(v string) mail.Mail {
	this.From = v
	return this
}

func (this *NaverMail) SetTo(v string) mail.Mail {
	this.To = v
	return this
}

func (this *NaverMail) SetUser(v string) mail.Mail {
	this.User = v
	return this
}

func (this *NaverMail) SetPassword(v string) mail.Mail {
	this.Password = v
	return this
}

func (this *NaverMail) SetSubject(v string) mail.Mail {
	this.Subject = v
	return this
}

func (this *NaverMail) SetBody(v string) mail.Mail {
	this.Body = v
	return this
}

func (this *NaverMail) Build() string {
	return "From: " + this.From + "\n" +
		"To: " + this.To + "\n" +
		"Subject: " + this.Subject + "\n" +
		"Content-Type: multipart/alternative;" + "\n" +
		"       boundary=\"" + this.boundary + "\"" + "\n\n" +
		this.Body + "\n"
}
