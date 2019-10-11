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

type AlertFunc func(string) error
type RecoverFunc func() error
type CheckFunc func(string) error

func defaultAlert(msg string) error {
	fmt.Fprintln(os.Stderr, msg)
	return nil
}

func defaultRecover() error {
	fmt.Fprintln(os.Stdout, "Server recovered")
	return nil
}

func defaultCheck(t string) error {
	resp, err := http.Get(t)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Server is not OK. status:[%v] %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	logrus.Info("Server is OK")
	return nil
}

type AlerterConfig struct {
	CheckIntervalDuration time.Duration
	AlertIntervalDuration time.Duration
}

type Alerter struct {
	conf        *AlerterConfig
	target      string
	alertFunc   AlertFunc
	recoverFunc RecoverFunc
	checkFunc   CheckFunc

	waitGroup sync.WaitGroup
	lock      sync.RWMutex
	running   bool
	ctx       context.Context
	close     func()

	lastAlert *alert
}

type alert struct {
	LastTime time.Time
	Msg      string
}

func NewAlerter(conf *AlerterConfig) *Alerter {
	a := &Alerter{
		conf:        conf,
		alertFunc:   defaultAlert,
		recoverFunc: defaultRecover,
		checkFunc:   defaultCheck,
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

	a.alertFunc = f
	return a
}

func (a *Alerter) SetRecover(f RecoverFunc) *Alerter {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.recoverFunc = f
	return a
}

func (a *Alerter) SetCheck(f CheckFunc) *Alerter {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.checkFunc = f
	return a
}

func (a *Alerter) validate() error {
	if a.alertFunc == nil {
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

		err := a.checkFunc(a.target)
		if err != nil {
			a.sendAlert(err.Error())
		} else {
			a.sendRecover()
		}
	}
}

func (a *Alerter) sendAlert(msg string) {
	now := time.Now()

	if a.lastAlert == nil {
		a.lastAlert = &alert{
			Msg: msg,
		}
	} else {
		alertPivot := a.lastAlert.LastTime.Add(a.conf.AlertIntervalDuration)
		if alertPivot.After(now) {
			logrus.Warn("Skip sending alert during alert interval[%v]. (left: %v)", a.conf.AlertIntervalDuration, alertPivot.Sub(now))
			return
		}
	}
	logrus.Info("send alert msg: ", msg)
	a.alertFunc(msg)
	a.lastAlert.LastTime = now
}

func (a *Alerter) sendRecover() {
	if a.lastAlert != nil {
		a.lastAlert = nil
		logrus.Info("send recovery msg")
		a.recoverFunc()
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
