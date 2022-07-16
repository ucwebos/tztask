package utils

func MurmurHash2(key string) (hash uint32) {
	var seed uint32 = 1532637697
	const m uint32 = 0x5bd1e995
	const r = 24

	var l int = len(key)
	var h uint32 = seed ^ uint32(l)

	var data = []byte(key)

	var k uint32

	for l >= 4 {
		k = uint32(data[0]) + uint32(data[1])<<8 + uint32(data[2])<<16 + uint32(data[3])<<24

		k *= m
		k ^= k >> r
		k *= m

		h *= m
		h ^= k

		data = data[4:]
		l -= 4
	}

	switch l {
	case 3:
		h ^= uint32(data[2]) << 16
		fallthrough
	case 2:
		h ^= uint32(data[1]) << 8
		fallthrough
	case 1:
		h ^= uint32(data[0])
		h *= m
	}

	h ^= h >> 13
	h *= m
	h ^= h >> 15

	return h
}
