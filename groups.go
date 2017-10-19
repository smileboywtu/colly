package colly

import (
	"sync"
	"time"
)

// Group several spider together and schedule them all

// trigger for spider to schedule
type Trigger func()

// Runnable Spider
type Runnable func()

type Group struct {
	// spider list
	spiders []Runnable

	// group name
	Name string

	// group duration
	period time.Duration

	// Parallelism limit for group
	Parallelism int

	// wait chan
	waitChan chan bool

	wg   *sync.WaitGroup
	lock *sync.Mutex
}

// New Group method create a new group
func NewGroup(Name string, Parallelism int) *Group {

	chansize := 1
	if Parallelism > 0 {
		chansize = Parallelism
	}

	return &Group{
		Name:        Name,
		spiders:     make([]Runnable, 0, 10),
		Parallelism: chansize,
		period:      0 * time.Second,
		waitChan:    make(chan bool, chansize),
		wg:          &sync.WaitGroup{},
		lock:        &sync.Mutex{},
	}

}

// Add new Spider to Group
func (g *Group) AddSpider(Runner Runnable) {
	g.lock.Lock()
	g.spiders = append(g.spiders, Runner)
	g.lock.Unlock()
}

// Wait all Spider Done
func (g *Group) Wait() {
	g.wg.Wait()
}

// Wrapper for single spider
func (g *Group) RunSpider(index int, Runner Runnable) {

	// wait done or wait another channel done
	g.waitChan <- true
	defer func() {
		<-g.waitChan
	}()

	g.wg.Add(1)
	defer func() {
		g.wg.Done()
	}()

	// run spider
	Runner()
}

// Run all Spiders in a single turn without trigger and period
func (g *Group) RunPending() {

	for index, runner := range g.spiders {
		g.RunSpider(index, runner)
	}
}

// Run all spider forever and wait for scheduler to schedule
func (g *Group) RunForever() {

	for true {
		// run now
		g.RunPending()

		time.Sleep(g.period)
	}

}

// TimeSchedule run the group with time period
func (g *Group) TimeSchedule(Duration time.Duration) *Group {
	g.period = Duration
	// chain design
	return g
}
