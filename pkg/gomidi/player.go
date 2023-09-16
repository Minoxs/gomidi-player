// Package gomidi deals with the lower level stuff required to play MIDI songs.
package gomidi

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // Register Player driver
	"gitlab.com/gomidi/midi/v2/smf"

	"git.mino.one/gomidi-player/internal/utils"
)

var (
	ErrNoSong     = errors.New("song is not set")
	ErrIsPlaying  = errors.New("already playing a song")
	ErrIsStopped  = errors.New("already stopped")
	ErrInvalidSMF = errors.New("failed to read tracks from SMF file")
	ErrNoOutPort  = errors.New("failed to get sound out port")

	errDone    = errors.New("done playing song")
	errStopped = errors.New("stopped playing")
	errPaused  = errors.New("paused song")
)

type (
	// songUnit is the smallest unit that can be played.
	songUnit struct {
		smf.Message
		Interval time.Duration
	}

	// Player plays songs in MIDI format.
	// Initialize using midi.NewPlayer.
	Player struct {
		mutexPlay sync.RWMutex
		isPlaying bool
		ctx       context.Context
		cancel    context.CancelCauseFunc

		currentDuration  time.Duration
		totalDuration    time.Duration
		currentUnitIndex int
		currentSong      []songUnit

		DefaultBPM    float64
		DefaultVolume uint8
	}
)

// NewPlayer initializes the Player.
func NewPlayer() *Player {
	return &Player{
		ctx:           utils.UnavailableContext(),
		cancel:        func(cause error) {},
		DefaultBPM:    180,
		DefaultVolume: 100,
	}
}

// IsPlaying returns whether the player is currently playing something.
func (p *Player) IsPlaying() bool {
	return p.isPlaying
}

// SetSongFromSMF takes in a file in SMF format and returns units that can be played.
func (p *Player) SetSongFromSMF(buffer io.Reader) (err error) {
	p.mutexPlay.Lock()
	defer p.mutexPlay.Unlock()

	if p.IsPlaying() {
		return ErrIsPlaying
	}

	p.currentUnitIndex = 0
	p.currentDuration = 0
	p.currentSong = nil

	var events = SMFReader{}
	err = events.ReadSMF(buffer)
	if err != nil {
		return
	}

	p.currentSong, p.totalDuration = events.GetSongUnits()
	return
}

// Play will play the song set.
// This is a non-blocking operation. Use Wait to wait until song finishes.
func (p *Player) Play() error {
	p.mutexPlay.Lock()
	defer p.mutexPlay.Unlock()

	var out, err = p.initializePlayer()
	if err != nil {
		return err
	}

	go p.playOn(out)
	return nil
}

// Stop stops the song currently playing.
func (p *Player) Stop() (err error) {
	return p.sendStopSignal(errStopped)
}

// Pause pauses the song currently playing.
func (p *Player) Pause() (err error) {
	return p.sendStopSignal(errPaused)
}

// Wait will block until the player stops playing.
func (p *Player) Wait() {
	for p.IsPlaying() {
		time.Sleep(50 * time.Millisecond)
	}
}

// SongDuration is the total duration of the song.
func (p *Player) SongDuration() time.Duration {
	return p.totalDuration
}

// SongCurrent is the current time on the song.
func (p *Player) SongCurrent() time.Duration {
	return p.currentDuration
}

// SongRemaining is the remaining duration of the song.
func (p *Player) SongRemaining() time.Duration {
	return p.totalDuration - p.currentDuration
}
