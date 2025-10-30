package extract

import (
	"bufio"
	"bytes"
	"io"
)

func UrLs(r io.Reader) ([]string, error) {
	br := bufio.NewReader(r)
	data, err := io.ReadAll(br)
	if err != nil {
		return nil, err
	}

	var out []string
	// Search for "http://" or "https://" occurrences.
	i := 0
	for i < len(data) {
		j := indexHTTP(data, i)
		if j == -1 {
			break
		}

		end := j
		for end < len(data) && allowedURLByte(data[end]) {
			end++
		}
		if end == j {
			i = j + 1
			continue
		}

		out = append(out, string(data[j:end]))
		i = end
	}

	return out, nil
}

func indexHTTP(b []byte, from int) int {
	i := from
	for i < len(b) {
		if b[i] != 'h' {
			i++
			continue
		}
		if i+7 < len(b) && bytes.Equal(b[i:i+8], []byte("https://")) {
			return i
		}
		i++
	}
	return -1
}

func allowedURLByte(c byte) bool {
	if c < 0x20 || c == 0x7f {
		return false
	}
	switch c {
	case ' ', '"', '\'', '<', '>', '`', '\\':
		return false
	}
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
		return true
	}
	switch c {
	case '-', '.', '_', '~', ':', '/', '?', '#', '[', ']', '@', '!', '$', '&', '(', ')', '*', '+', ',', ';', '=': // RFC3986
		return true
	case '%', '|':
		return true
	}
	return false
}
