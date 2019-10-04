package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/jmshin92/alerter/pkgs/mail"
)

const (
	SubjectFail    = "Server failed"
	SubjectRecover = "Server recovered"
)

var (
	confPath  string
	lastAlert *Alert
	c         *Conf
	m         mail.Mail
)

type Conf struct {
	Interval           int           `json:"interval" yaml:"interval"`
	IntervalDuration   time.Duration `json:"-" yaml:"-"`
	AlertDelay         int           `json:"alert_interval" yaml:"alert"`
	AlertDelayDuration time.Duration `json:"-" yaml:"-"`
	Address            string        `json:"address" yaml:"address"`
	MailVendor         string        `json:"mail_vendor" yaml:"mail_vendor"`
	Mail               mail.MailConf `json:"mail" yaml:"mail"`
}

type Alert struct {
	LastAlertTime time.Time
	CreatedTime   time.Time
}

func loadConf(p string) (*Conf, error) {
	c = &Conf{}
	confBytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(confBytes, c); err != nil {
		return nil, err
	}
	return c, nil
}

func main() {
	flag.StringVar(&confPath, "c", "", "config path")
	flag.Parse()

	if len(confPath) == 0 {
		logrus.Errorln("config path is required")
		flag.Usage()
		os.Exit(-1)
	}

	conf, err := loadConf(confPath)
	if err != nil {
		err = errors.Wrapf(err, "failed to load config file[%v]", confPath)
		logrus.Errorln(err)
		os.Exit(-1)
	}

	interval := time.Duration(conf.Interval) * time.Second
	for {
		select {
		case <-time.After(interval):
		}

		resp, err := http.Get(conf.Address)
		if err != nil {
			err = errors.Wrapf(err, "failed to send request to [%v]", conf.Address)
			alert(conf, err.Error())
			continue
		}

		logrus.Infoln("Status code[%v]", resp.StatusCode)

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("server status[%v] is not OK[%v]", resp.StatusCode, http.StatusOK)
			alert(conf, err.Error())
			continue
		}

		recoverAlert()
	}
}

func alert(conf *Conf, data string) error {
	now := time.Now()
	if lastAlert == nil {
		lastAlert = &Alert{
			CreatedTime: now,
		}
	}

	alertPivot := lastAlert.LastAlertTime.Add(conf.AlertDelayDuration)
	if alertPivot.After(now) {
		logrus.Infoln("skip alert until pivot[%v]", alertPivot)
		return nil
	}

	err := mail.NewMail(mail.NaverMail).
		SetFrom(conf.Mail.From).
		SetTo(conf.Mail.To).
		SetUser(conf.Mail.User).
		SetPassword(conf.Mail.Password).
		SetSubject(SubjectFail).
		SetBody(data).
		Send()

	if err != nil {
		fmt.Println(err)
		return err
	}

	lastAlert.LastAlertTime = now
	return nil
}

func recoverAlert() {
	if lastAlert == nil {
		return
	}

	err := mail.NewMail(mail.NaverMail, conf).
		SetSubject(SubjectFail).
		SetBody("recovered").
		Send()

	if err != nil {
		fmt.Println(err)
		return err
	}

	lastAlert = nil
}
