package ticker

import (
	"time"
)

const (
	kDelay          = true
	kNoDelay        = false
	defalutInterval = 60
)

func TickerAfterDone(interval int, function func()) {
	Ticker(interval, function)
}

func Ticker(interval int, function func()) {
	function()
	tickerDone(interval, function)
}

func TickerDelayed(interval int, function func()) {
	tickerDone(interval, function)
}

func AsyncTicker(interval int, function func()) {
	go function()
	tickerDone(interval, function)
}

func tickerDone(interval int, function func()) {

	if interval <= 0 {
		interval = defalutInterval
	}

	go func() {
		eventsTick := time.NewTicker(time.Duration(interval) * time.Second)
		defer eventsTick.Stop()

		for {
			select {
			case <-eventsTick.C:
				function()
			}
		}
	}()
}
