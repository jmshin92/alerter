package mail

import "oss.navercorp.com/clous/indexer/cmd/alerter/mail/naver"

type Mail interface {
	Send() error
	SetFrom(string) Mail
	SetTo(string) Mail
	SetUser(string) Mail
	SetPassword(string) Mail
	SetSubject(string) Mail
	SetBody(string) Mail
	Build() string
}

type MailConf struct {
	From     string `json:"from,omitempty" yaml:"from,omitempty"`
	To       string `json:"to,omitempty" yaml:"to,omitempty"`
	User     string `json:"user,omitempty" yaml:"user,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

type MailVendor int

const (
	NaverMail MailVendor = iota
	Gmail
)

func NewMail(m MailVendor, c *MailConf) Mail {
	switch m {
	case NaverMail:
		return naver.NewNaverMail(c)
	}
	return nil
}
