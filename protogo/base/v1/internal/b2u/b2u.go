package b2u

func B264X1(b []byte) uint64 {
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

func B264X2(b []byte) (uint64, uint64) {
	return B264X1(b), B264X1(b[8:])
}

func B232X1(b []byte) uint32 {
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

func U264X1(u uint64) [8]byte {
	return [8]byte{
		byte(u >> 56), byte(u >> 48), byte(u >> 40), byte(u >> 32),
		byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u),
	}
}

func U264X2(hi, lo uint64) [16]byte {
	return [16]byte{
		byte(hi >> 56), byte(hi >> 48), byte(hi >> 40), byte(hi >> 32),
		byte(hi >> 24), byte(hi >> 16), byte(hi >> 8), byte(hi),

		byte(lo >> 56), byte(lo >> 48), byte(lo >> 40), byte(lo >> 32),
		byte(lo >> 24), byte(lo >> 16), byte(lo >> 8), byte(lo),
	}
}

func U232X1(u uint32) [4]byte {
	return [4]byte{
		byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u),
	}
}

func B20(b []byte) (uint64, uint64, uint32) {
	return B264X1(b), B264X1(b[8:]), B232X1(b[16:])
}

func U20(b0, b1 uint64, b2 uint32) [20]byte {
	return [20]byte{
		byte(b0 >> 56), byte(b0 >> 48), byte(b0 >> 40), byte(b0 >> 32),
		byte(b0 >> 24), byte(b0 >> 16), byte(b0 >> 8), byte(b0),

		byte(b1 >> 56), byte(b1 >> 48), byte(b1 >> 40), byte(b1 >> 32),
		byte(b1 >> 24), byte(b1 >> 16), byte(b1 >> 8), byte(b1),

		byte(b2 >> 24), byte(b2 >> 16), byte(b2 >> 8), byte(b2),
	}
}
