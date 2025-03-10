package btils

import (
	"unsafe"
)

// Do NOT touch. Otherwise you might run into oob exceptions
const randChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

// In no way shape or form associated with UUIDs defined in rfc4122 (https://datatracker.ietf.org/doc/html/rfc4122)
// Generations are predictable and should not be used for cryptographic applications.
// UID merely stands for "Unique IDentifier" Which is guaranteed with 79.228.162.514.264.337.593.543.950.336 possible
// values, with a 10% first-time-collision probability at 129.209.288.033.988 generations
// and a 50% first-time-collision-probability at 331.411.458.666.437 generations.
type UID [16]byte

func UIDFromString(s string) *UID {
	if len(s) < 16 {
		return nil
	}
	return (*UID)(unsafe.Pointer(unsafe.StringData(s)))
}

func (uid *UID) ToString() string {
	return unsafe.String(unsafe.SliceData(uid[:]), 16)
}

// This function should only be used if you need to validate that the UID does not contain malicious contents (e.g. XSS, SQL injection, etc) otherwise accept uid as-is
func (uid UID) IsValid() bool {
	for i := 0; i < 16; i++ {
		b := uid[i]
		if (b >= 'a' && b <= 'z') ||
			(b >= 'A' && b <= 'Z') ||
			(b >= '0' && b <= '9') ||
			b == '_' || b == '-' {
			continue
		}
		return false
	}
	return true
}

// Might seem counter-intuitive to give a UID, tho this allows rapid uid creation by re-using old UIDs
func NewUID(b *UID) {
	rnd1 := Fastrand()
	rnd2 := Fastrand()
	rnd3 := Fastrand()

	b[0] = randChars[rnd1&63]
	b[1] = randChars[(rnd1>>6)&63]
	b[2] = randChars[(rnd1>>12)&63]
	b[3] = randChars[(rnd1>>18)&63]
	b[4] = randChars[(rnd1>>24)&63]

	b[5] = randChars[rnd2&63]
	b[6] = randChars[(rnd2>>6)&63]
	b[7] = randChars[(rnd2>>12)&63]
	b[8] = randChars[(rnd2>>18)&63]
	b[9] = randChars[(rnd2>>24)&63]

	b[10] = randChars[rnd3&63]
	b[11] = randChars[(rnd3>>6)&63]
	b[12] = randChars[(rnd3>>12)&63]
	b[13] = randChars[(rnd3>>18)&63]
	b[14] = randChars[(rnd3>>24)&63]

	/*
		At this point each random number has 2 bits left over, 30 and 31.
		So we combine all of them to form the 16th character
	*/
	b[15] = randChars[((rnd1>>30)&3)|(((rnd2>>30)&3)<<2)|(((rnd3>>30)&3)<<4)]
}
