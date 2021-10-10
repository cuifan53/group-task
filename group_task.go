package group_task

import (
	"sync"
)

func NewGroupTask(options ...Option) *GroupTask {
	p := &GroupTask{
		groupMap:  make(map[string]chan func()),
		queueSize: 100,
	}

	for _, v := range options {
		v(p)
	}

	return p
}

type Option func(p *GroupTask)

func WithQueueSize(queueSize int) Option {
	return func(p *GroupTask) {
		p.queueSize = queueSize
	}
}

func WithPanicHandler(panicHandler func(interface{})) Option {
	return func(p *GroupTask) {
		p.panicHandler = panicHandler
	}
}

type GroupTask struct {
	mu           sync.Mutex
	groupMap     map[string]chan func()
	queueSize    int
	panicHandler func(interface{})
}

func (p *GroupTask) Do(group string, task func()) {
	p.mu.Lock()
	defer p.mu.Unlock()

	c, ok := p.groupMap[group]
	if ok {
		c <- task
		return
	}

	ch := make(chan func(), p.queueSize)
	p.groupMap[group] = ch
	go p.do(group, ch)
	ch <- task
}

func (p *GroupTask) do(group string, ch chan func()) {
	defer func() {
		if err := recover(); err != nil {
			if p.panicHandler != nil {
				p.panicHandler(err)
			}
		}
		p.mu.Lock()
		delete(p.groupMap, group)
		p.mu.Unlock()
	}()

	for task := range ch {
		task()
		if len(ch) == 0 {
			close(ch)
			return
		}
	}
}
