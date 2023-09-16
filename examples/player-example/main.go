package main

import (
	"bytes"
	"log"

	"git.mino.one/gomidi-player/pkg/gomidi"
)

func main() {
	var song = bytes.NewBuffer(
		[]byte{
			0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x01, 0x00, 0x60, 0x4d, 0x54,
			0x72, 0x6b, 0x00, 0x00, 0x00, 0x3a, 0x00, 0xb0, 0x27, 0x7f, 0x00, 0x07, 0x64, 0x00, 0xff, 0x58,
			0x04, 0x03, 0x02, 0x08, 0x08, 0x00, 0xff, 0x51, 0x03, 0x05, 0x16, 0x15, 0x00, 0x90, 0x21, 0x78,
			0x60, 0x80, 0x21, 0x00, 0x00, 0x90, 0x23, 0x78, 0x60, 0x80, 0x23, 0x00, 0x00, 0x90, 0x23, 0x78,
			0x60, 0x80, 0x23, 0x00, 0x00, 0x90, 0x21, 0x78, 0x60, 0x80, 0x21, 0x00, 0x00, 0xff, 0x2f, 0x00,
		},
	)

	var (
		err    error
		player = gomidi.NewPlayer()
	)

	err = player.SetSongFromSMF(song)
	if err != nil {
		log.Fatal(err)
	}

	err = player.Play()
	if err != nil {
		log.Fatal(err)
	}

	player.Wait()
}