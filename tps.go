package tps

import (
	"time"
	"math/rand"
)

type Tps struct {
	done chan struct{}
	ticker *time.Ticker
	events chan struct{}
	p float64
}

func New(tps float64) *Tps {
	return &Tps{
		done: make(chan struct{}, 1),
		ticker: time.NewTicker(time.Millisecond),
		events: make(chan struct{}),
		p: tps / 1000,
	}
}

func (t *Tps) Run(cancel <-chan struct{}) {
	for {
		select {
		case <-cancel:
			t.done <-struct{}{}
			return
		case <-t.ticker.C:
			if sendEvent(t.p) {
				select {
				case t.events <- struct{}{}:
				default:
					// drop it on the floor
				}
			}
		}
	}
}

func (t *Tps) Done() <-chan struct{} {
	return t.done
}

func (t *Tps) Events() <-chan struct{} {
	return t.events
}

func sendEvent(p float64) bool {
	return rand.Float64() < p
}