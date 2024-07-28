package sanitize

import (
	"encoding/hex"
	"fmt"
)

func isHexUpper(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')
}

func ConvertQoutedPrintable(s string) string {
	idx := 0
	resp := []byte{}
	var err error
	hexInput := []byte{0, 0}

	for {
		if idx >= len(s) {
			break
		}

		c := s[idx]
		if c != '=' {
			resp = append(resp, c)
			idx++
			continue
		}

		if idx+2 >= len(s) {
			resp = append(resp, c)
			idx++
			continue
		}

		hexInput[0] = s[idx+1]
		hexInput[1] = s[idx+2]

		if hexInput[0] == '\r' && hexInput[1] == '\n' {
			idx += 3
			continue
		}

		if !isHexUpper(hexInput[0]) || !isHexUpper(hexInput[1]) {
			resp = append(resp, c)
			idx++
			continue
		}

		resp, err = hex.AppendDecode(resp, hexInput)
		if err != nil {
			fmt.Printf("DECODEING ERROR, Failed to decode %s to a byte: %s", string(hexInput), err)
			resp = append(resp, c)
			idx++
			continue
		}

		idx += 3
		continue
	}

	return string(resp)
}
