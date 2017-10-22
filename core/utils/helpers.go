package utils

import "crypto/rand"

func ArrayOfBytes(i int, b byte) (p []byte) {

	for i != 0 {

		p = append(p, b)
		i--
	}
	return
}

func FitBytesInto(d []byte, i int) []byte {

	if len(d) < i {

		dif := i - len(d)

		return append(ArrayOfBytes(dif, 0), d...)
	}

	return d[:i]
}

func StripByte(d []byte, b byte) []byte {
	// Given a byte with b padding
	// functions returns the bytes with padding stripped

	for i, bb := range d {

		if bb != b {
			return d[i:]
		}
	}

	return nil
}

func Max(a, b int) int {

	if a >= b {

		return a
	}

	return b
}

func Min(a, b int) int {

	if a <= b {

		return a
	}

	return b
}

func RandomInt(a, b int) int {

	var bytes = make([]byte, 1)
	rand.Read(bytes)

	per := float32(bytes[0]) / 256.0
	dif := Max(a, b) - Min(a, b)

	return Min(a, b) + int(per*float32(dif))
}

// From http://devpy.wordpress.com/2013/10/24/create-random-string-in-golang/
func RandomString(n int) string {

	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
