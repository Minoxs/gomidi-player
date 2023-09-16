# Gomidi Player

[![Go Reference](https://pkg.go.dev/badge/github.com/Minoxs/gomidi-player.svg)](https://pkg.go.dev/github.com/Minoxs/gomidi-player)

This MIDI player is built using [gomidi](https://gitlab.com/gomidi/midi), and is made to be used as a library.
The functionality is fairly basic, but it should cover the basics of reading an SMF file and playing its sounds like a music player would.

The player is run in a separate goroutine, and is safe to be manipulated by multiple goroutines.
Feel free to modify this to fit your needs, or use as a starting point. Credits are, as always, greatly appreciated.

# TODO

- [ ] Turn the output driver into a parameter
