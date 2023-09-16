package gomidi

import (
	"context"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"

	"git.mino.one/gomidi-player/internal/utils"
)

// initializePlayer does all the checks required to start playing.
// Returns the available out port on success.
func (p *Player) initializePlayer() (out drivers.Out, err error) {
	// Check if song was set
	if p.currentSong == nil {
		err = ErrNoSong
		return
	}

	// Check if playing
	if p.IsPlaying() {
		err = ErrIsPlaying
		return
	}

	// Get port to play on
	out, err = p.getOutPort()
	if err != nil {
		err = utils.WrapOnError(err, ErrNoOutPort)
		return
	}

	p.isPlaying = true
	p.initializeContext()

	return
}

// getOutPort will get the main out port
func (p *Player) getOutPort() (out drivers.Out, err error) {
	out, err = midi.OutPort(0)
	if err != nil {
		return
	}

	err = out.Open()
	return
}

// initializeContext will create a new context to play with
func (p *Player) initializeContext() {
	p.ctx, p.cancel = context.WithCancelCause(context.Background())
}

// sendStopSignal will signal the player to stop playing
func (p *Player) sendStopSignal(cause error) error {
	if !p.IsPlaying() {
		return ErrIsStopped
	}
	p.cancel(cause)
	return nil
}
