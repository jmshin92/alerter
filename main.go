package main

import (
	"flag"
	"fmt"
	"github.com/jmshin92/alerter/pkgs/alerter"
	"github.com/jmshin92/alerter/pkgs/config"
	"github.com/jmshin92/alerter/pkgs/mail"
	"github.com/jmshin92/alerter/pkgs/mail/factory"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	ConfPath string
	Vendors  bool
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}

func main() {
	flag.BoolVar(&Vendors, "v", false, "list of supported vendors")
	flag.StringVar(&ConfPath, "c", "", "config path")
	flag.Parse()

	if Vendors {
		fmt.Println(mail.Vendors())
		os.Exit(0)
	}
	if len(ConfPath) == 0 {
		logrus.Error("config path is mandatory")
		flag.Usage()
		os.Exit(-1)
	}

	c, err := config.GetConfig(ConfPath)
	if err != nil {
		logrus.Error("failed to get config. error: ", err)
		os.Exit(-1)
	}

	alert := func(msg string) error {
		mail, err := mail_factory.NewMail(&c.Mail.MailConf)
		if err != nil {
			logrus.Error(err)
			return err
		}
		mail.SetSubject("Server is not OK")
		mail.SetBody(fmt.Sprintf("[%v] Server is not OK!\n%s", time.Now(), msg))
		return mail.Send()
	}

	recover := func() error {
		mail, err := mail_factory.NewMail(&c.Mail.MailConf)
		if err != nil {
			logrus.Error(err)
			return err
		}
		mail.SetSubject("Server recovered")
		mail.SetBody(fmt.Sprintf("[%v] Server recovered!", time.Now()))
		return mail.Send()
	}

	err = alerter.NewAlerter(c.Alert.AlerterConfig).
		SetAlert(alert).
		SetRecover(recover).
		SetTarget(c.TargetUri).
		Run()

	if err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
}
