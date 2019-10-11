package mail_factory

import (
	"fmt"
	"github.com/jmshin92/alerter/pkgs/mail"
	"github.com/jmshin92/alerter/pkgs/mail/naver"
)

func NewMail(c *mail.MailConf) (mail.Mail, error) {
	switch c.Vendor {
	case mail.VendorNaver:
		return naver.NewNaverMail(c)
	default:
		return nil, fmt.Errorf("unsupported vendor[%v]", c.Vendor)
	}
}
