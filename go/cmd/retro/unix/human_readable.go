package unix

import "fmt"

// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format
func HumanReadable(byteCount int64) string {
	const unit = 1024
	if byteCount < unit {
		return fmt.Sprintf("%d B", byteCount)
	}
	div, exp := int64(unit), 0
	for n := byteCount / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(byteCount)/float64(div), "KMGTPE"[exp])
}
