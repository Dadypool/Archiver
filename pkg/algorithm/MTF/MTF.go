package MTF

func Encode(data []byte) []byte {
	// 1. Init alphabet
	alphabet := make([]byte, 256)
	for i := range alphabet {
		alphabet[i] = byte(i)
	}

	encoded := make([]byte, 0, len(data))

	for _, c := range data {
		// 2. Find index
		idx := byte(0)
		for i, b := range alphabet {
			if b == c {
				idx = byte(i)
				break
			}
		}

		// 3. Index to result
		encoded = append(encoded, idx)

		// 4. Move to front
		if idx > 0 {
			if idx == 255 {
				copy(alphabet[1:], alphabet[:255])
			} else {
				copy(alphabet[1:idx+1], alphabet[:idx])
			}
			alphabet[0] = c
		}
	}

	return encoded
}

func Decode(encoded []byte) []byte {
	// 1. Init alphabet
	alphabet := make([]byte, 256)
	for i := range alphabet {
		alphabet[i] = byte(i)
	}

	decoded := make([]byte, 0, len(encoded))

	for _, idx := range encoded {
		// 2. Find index
		c := alphabet[idx]
		decoded = append(decoded, c)

		// 3. Move to front
		if idx > 0 {
			if idx == 255 {
				copy(alphabet[1:], alphabet[:255])
			} else {
				copy(alphabet[1:idx+1], alphabet[:idx])
			}
			alphabet[0] = c
		}
	}

	return decoded
}
