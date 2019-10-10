package config

import (
	"encoding/json"
	"fmt"
	"github.com/jmshin92/alerter/pkgs/mail"
)

const (
	defaultSubject = "Server is not OK"
)

type Mail struct {
	mail.MailConf `yaml:",inline"`
}

func NewMail() *Mail {
	return &Mail{
		MailConf: mail.MailConf{
			Subject: defaultSubject,
		},
	}
}

func (this *Mail) UnmarshalJSON(data []byte) error {
	if this == nil {
		return fmt.Errorf("alert is nil")
	}

	type Alias Mail
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(this),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if err := this.MailConf.Validate(); err != nil {
		return err
	}

	return nil
}

func (this *Mail) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if this == nil {
		return fmt.Errorf("alert is nil")
	}

	type Alias Mail
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(this),
	}

	if err := unmarshal(aux.Alias); err != nil {
		return err
	}

	if err := this.MailConf.Validate(); err != nil {
		return err
	}

	return nil
}
