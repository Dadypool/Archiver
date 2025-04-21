package RLE

func Encode(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}

	var encoded []byte
	current := data[0]
	count := 1

	for i := 1; i < len(data); i++ {
		if data[i] == current && count < 255 {
			count++
		} else {
			encoded = append(encoded, byte(count), current)
			current = data[i]
			count = 1
		}
	}
	encoded = append(encoded, byte(count), current)

	return encoded
}

func Decode(encoded []byte) []byte {
	var decoded []byte
	for i := 0; i < len(encoded); i += 2 {
		count := int(encoded[i])
		symbol := encoded[i+1]
		for j := 0; j < count; j++ {
			decoded = append(decoded, symbol)
		}
	}
	return decoded
}
