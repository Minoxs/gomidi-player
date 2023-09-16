package gomidi

import (
	"io"
	"sort"
	"time"

	"gitlab.com/gomidi/midi/v2/smf"

	"git.mino.one/gomidi-player/internal/utils"
)

type SMFReader struct {
	currentSong []smf.TrackEvent
}

// ReadSMF reads an SMF file and parses all the track events.
func (e *SMFReader) ReadSMF(buffer io.Reader) (err error) {
	e.currentSong = make([]smf.TrackEvent, 0, 100)
	return utils.WrapOnError(
		smf.ReadTracksFrom(buffer).Do(e.readEvent).Error(),
		ErrInvalidSMF,
	)
}

// readEvent reads a single event from the track.
func (e *SMFReader) readEvent(event smf.TrackEvent) {
	if event.Message.IsPlayable() {
		e.currentSong = append(e.currentSong, event)
	}
}

// GetSongUnits parses the track events and returns the playable units.
func (e *SMFReader) GetSongUnits() (units []songUnit, totalDuration time.Duration) {
	e.sortSong()
	units = make([]songUnit, len(e.currentSong))
	totalDuration = 0
	for i := 0; i < len(units); i++ {
		var event = e.currentSong[i]
		units[i] = songUnit{
			Message:  event.Message,
			Interval: time.Microsecond * time.Duration(event.AbsMicroSeconds-totalDuration.Microseconds()),
		}
		totalDuration += units[i].Interval
	}
	return
}

// sortSong makes sure the song is ordered by time.
func (e *SMFReader) sortSong() {
	sort.SliceStable(
		e.currentSong, func(i, j int) bool {
			return e.currentSong[i].AbsMicroSeconds < e.currentSong[j].AbsMicroSeconds
		},
	)
}
