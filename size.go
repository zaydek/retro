package main

import (
	"fmt"
)

// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format
func ByteCount(b int64) string {
	const u = 1024

	if b < u {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(u), 0
	for n := b / u; n >= u; n /= u {
		div *= u
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func main() {
	fmt.Println(ByteCount(1023 * 1024))
}
