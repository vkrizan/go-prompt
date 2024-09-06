// +build !windows

package term

import (
	"sync"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

var (
	saveTermios     unix.Termios
	saveTermiosFD   int
	saveTermiosOnce sync.Once
)

func getOriginalTermios(fd int) (unix.Termios, error) {
	var err error
	saveTermiosOnce.Do(func() {
		saveTermiosFD = fd
		var original *unix.Termios
		original, err = termios.Tcgetattr(uintptr(fd))
		saveTermios = *original

	})
	return saveTermios, err
}

// Restore terminal's mode.
func Restore() error {
	o, err := getOriginalTermios(saveTermiosFD)
	if err != nil {
		return err
	}
	return termios.Tcsetattr(uintptr(saveTermiosFD), termios.TCSANOW, &o)
}
