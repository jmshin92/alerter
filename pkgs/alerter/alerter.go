package alerter

import "sync"

type Alerter struct {
	waitGroup sync.WaitGroup
}

func (a *Alerter) Run() {

	a.waitGroup.Add(1)
	go a.run()
	a.waitGroup.Wait()
}

func (a *Alerter) run() {
	defer a.waitGroup.Done()
}
