package toolbox

import "fmt"

// StringByteSize converts a size (in bytes) to bytes, kibibytes (KiB),
// mebibyte (MiB), gibibyte (GiB), tebibyte (TiB), pebibyte (PiB) or
// exbibyte (EiB)
// The returned string includes the unit (nnn bytes, n.nn KiB, n.nn MiB, ...)
func StringByteSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d byte", size)
	}
	s := float64(size)
	for _, unit := range []string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB"} {
		s = s / 1024.0
		if s < 1024.0 {
			return fmt.Sprintf("%.2f %s", s, unit)
		}
	}
	// This won't happen since int64 can't be bigger than ~4.00 EiB
	return fmt.Sprintf("%.2f EiB", s)
}
