package mail_factory

import (
	"github.com/jmshin92/alerter/pkgs/mail"
	"github.com/jmshin92/alerter/pkgs/mail/naver"
)

func NewMail(c *mail.MailConf) mail.Mail {
	switch c.Vendor {
	case mail.VendorNaver:
		return naver.NewNaverMail(c)
	}
	return nil
}
