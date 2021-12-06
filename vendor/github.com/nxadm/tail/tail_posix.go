//go:build !windows
// +build !windows

package tail

import (
	"os"
)

func OpenFile(name string) (file *os.File, err error) {
	return os.Open(name)
}
