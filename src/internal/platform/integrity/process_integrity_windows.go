//go:build windows

package integrity

/*
#cgo LDFLAGS: -ladvapi32
#include "integrity.h"
*/
import "C"

const (
	UntrustedRID        = 0x00000000
	LowRID              = 0x00001000
	MediumRID           = 0x00002000
	HighRID             = 0x00003000
	SystemRID           = 0x00004000
	ProtectedProcessRID = 0x00005000
)

func GetProcessLevel(pid uint32) uint32 {
	return uint32(C.GetProcessLevel(C.uint(pid)))
}
