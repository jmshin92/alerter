package mail

import (
	"fmt"
)

const (
	VendorNaver = "Naver"
	VendorGmail = "Gmail"
)

var (
	vendors = map[string]interface{}{
		VendorNaver: nil,
		VendorGmail: nil,
	}
)

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
	Subject  string `json:"subject,omitempty" yaml:"subject,omitempty"`
	Body     string `json:"body,omitempty" yaml:"body,omitempty"`
	User     string `json:"user,omitempty" yaml:"user,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	Vendor   string `json:"vendor,omitempty" yaml:"vendor,omitempty"`
}

func (this *MailConf) Validate() error {
	if len(this.From) == 0 {
		return fmt.Errorf("from is required")
	}

	if len(this.To) == 0 {
		return fmt.Errorf("to is required")
	}

	if len(this.User) == 0 {
		return fmt.Errorf("user is required")
	}

	if len(this.Password) == 0 {
		return fmt.Errorf("password is required")
	}

	switch this.Vendor {
	case VendorNaver:
	case VendorGmail:
	default:
		return fmt.Errorf("unsupported vendor")
	}

	return nil
}

func Vendors() []string {
	v := make([]string, 0)
	for k := range vendors {
		v = append(v, k)
	}
	return v
}
