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
)

var (
	ConfPath string
	Vendors bool
)

func main() {
	flag.BoolVar(&Vendors, "v", false, "list of supported vendors")
	flag.StringVar(&ConfPath, "c", "", "config path")
	flag.Parse()

	if Vendors{
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

	alert := mail_factory.NewMail(&c.Mail.MailConf).Send

	err = alerter.NewAlerter(c.Alert.AlerterConfig).
		SetAlert(alert).
		SetTarget(c.TargetUri).
		Run()

	if err != nil {
		logrus.Error(err)
	}
}