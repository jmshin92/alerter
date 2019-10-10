package config

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmshin92/alerter/pkgs/alerter"
)

const (
	defaultCheckInterval = 10
	defaultAlertInterval = 60
)

type Alert struct {
	*alerter.AlerterConfig `json:"-" yaml:"-"`
	CheckInterval          int `json:"check_interval,omitempty" yaml:"check_interval,omitempty"`
	AlertInterval          int `json:"alert_interval,omitempty" yaml:"alert_interval,omitempty"`
}

func NewAlert() *Alert {
	a := &Alert{
		CheckInterval: defaultCheckInterval,
		AlertInterval: defaultAlertInterval,
	}
	a.AlerterConfig = &alerter.AlerterConfig{}
	a.convertTime()
	return a
}

func (this *Alert) UnmarshalJSON(data []byte) error {
	if this == nil {
		return fmt.Errorf("alert is nil")
	}

	type Alias Alert
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(this),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	this.convertTime()
	if err := this.validate(); err != nil {
		return err
	}
	return nil
}

func (this *Alert) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if this == nil {
		return fmt.Errorf("alert is nil")
	}

	type Alias Alert
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(this),
	}

	if err := unmarshal(aux.Alias); err != nil {
		return err
	}

	this.convertTime()
	if err := this.validate(); err != nil {
		return err
	}
	return nil
}

func (this *Alert) convertTime() {
	this.AlertIntervalDuration = time.Duration(this.AlertInterval) * time.Second
	this.CheckIntervalDuration = time.Duration(this.CheckInterval) * time.Second
}

func (this *Alert) validate() error {
	if this.CheckInterval <= 0 {
		return fmt.Errorf("check_interval should be larger than 0")
	}
	if this.AlertInterval < 0 {
		return fmt.Errorf("alert_interval cannot be negative")
	}
	return nil
}
