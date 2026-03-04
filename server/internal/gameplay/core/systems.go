package core

import (
	"time"

	"github.com/rouzbehsbz/zurvan"
)

type NetworkFlushSystem struct{}

func (n *NetworkFlushSystem) Update(w *zurvan.World, dt time.Duration) {
	sessionsManager := SessionManagerRes(w)
	sessionsManager.FlushAll()
}

type NetworkReceiveSystem struct{}

func (n *NetworkReceiveSystem) Update(w *zurvan.World, dt time.Duration) {
	dispatcher := DispatcherRes(w)
	nonBlockingPackets := NonBlockingPacketsRes(w)

	for len(nonBlockingPackets) > 0 {
		packet := <-nonBlockingPackets

		_ = dispatcher.Dispatch(packet)
	}
}
