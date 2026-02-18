//go:build windows

package window

/*
#cgo LDFLAGS: -luser32
#include "window.h"
*/
import "C"

func HasVisibleWindow(pid uint32) bool {
	return C.HasVisibleWindow(C.uint(pid)) != 0
}
