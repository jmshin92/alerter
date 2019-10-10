package alerter

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type AlertFunc func() error
type CheckFunc func(string) error

func defaultAlert() error {
	fmt.Fprintln(os.Stderr, "error alert")
	return nil
}

func defaultCheck(t string) error {
	resp, err := http.Get(t)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code[%v] is not OK[%v]", resp.StatusCode, http.StatusOK)
	}

	logrus.Info("Server is OK")
	return nil
}

type AlerterConfig struct {
	CheckIntervalDuration time.Duration
	AlertIntervalDuration time.Duration
}

type Alerter struct {
	conf   *AlerterConfig
	target string
	alert  AlertFunc
	check  CheckFunc

	waitGroup sync.WaitGroup
	lock      sync.RWMutex
	running   bool
	ctx       context.Context
	close     func()
}

func NewAlerter(conf *AlerterConfig) *Alerter {
	a := &Alerter{
		conf: conf,
		alert: defaultAlert,
		check: defaultCheck,
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.ctx = ctx
	a.close = func() {
		cancel()
		a.waitGroup.Wait()
	}

	return a
}

func (a *Alerter) SetTarget(target string) *Alerter {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.target = target
	return a
}

func (a *Alerter) SetAlert(f AlertFunc) *Alerter {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.alert = f
	return a
}

func (a *Alerter) SetCheck(f CheckFunc) *Alerter {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.check = f
	return a
}

func (a *Alerter) validate() error {
	if a.alert == nil {
		return fmt.Errorf("alert function is not set")
	}
	if len(a.target) == 0 {
		return fmt.Errorf("target is not set")
	}
	return nil
}

func (a *Alerter) Run() error {
	if err := a.RunAsync(); err != nil {
		return err
	}
	a.waitGroup.Wait()
	return nil
}

func (a *Alerter) RunAsync() error {
	if a.running {
		return fmt.Errorf("alerter is already running")
	}

	if err := a.validate(); err != nil {
		return err
	}

	var runGroup sync.WaitGroup
	runGroup.Add(1)
	go a.run(&runGroup)
	runGroup.Wait()
	return nil
}

func (a *Alerter) run(runGroup *sync.WaitGroup) {
	a.waitGroup.Add(1)
	a.running = true
	defer func() {
		a.running = false
		a.waitGroup.Done()
	}()
	runGroup.Done()

	for {
		select {
		case <-a.ctx.Done():
			logrus.Info("received stop signal")
			return
		case <-time.After(a.conf.CheckIntervalDuration):
		}
		logrus.Info("Check server")

		err := a.check(a.target)
		if err != nil {
			logrus.Error("server seems to be not ok")
			a.alert()
		}

	}
}

func (a *Alerter) Close() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	logrus.Info("closing alerter...")
	a.close()
	logrus.Debug("closed alerter!")
	return nil
}
