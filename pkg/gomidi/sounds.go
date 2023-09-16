package gomidi

import (
	"context"
	"errors"
	"runtime"
	"time"

	"gitlab.com/gomidi/midi/v2/drivers"

	"git.mino.one/gomidi-player/internal/utils"
)

// playOn will play the current song in the given out port
func (p *Player) playOn(out drivers.Out) {
	defer p.donePlaying()

	// Drivers may invoke CGO
	// Makes sure thread is locked to avoid weird errors
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer utils.IgnoreError(out.Close)

	// Creates timer to properly time sound writes
	var sleep = time.NewTimer(0)
	defer sleep.Stop()

	// Makes sure channel is drained
	<-sleep.C

	// Play each unit of song
	for i, unit := range p.currentSong[p.currentUnitIndex:] {
		p.currentDuration += unit.Interval
		p.currentUnitIndex = i
		if unit.Interval > 0 {
			sleep.Reset(unit.Interval)
			select {
			case <-sleep.C:
				break
			case <-p.ctx.Done():
				return
			}
		}
		_ = out.Send(unit.Message)
	}
}

// donePlaying does the cleaning up after playOn ends
func (p *Player) donePlaying() {
	_ = p.sendStopSignal(errDone)
	var err = context.Cause(p.ctx)

	switch {
	case errors.Is(err, errDone):
		fallthrough
	case errors.Is(err, errStopped):
		p.currentUnitIndex = 0
		p.currentDuration = 0
	}

	p.isPlaying = false
}
