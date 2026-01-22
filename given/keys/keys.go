/*
Copyright 2021, 2023 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more information.
*/

// The keys package contains helper functions for performing various
// operations on KDHT keys.  The purpose of this package is to remove
// some of the effort of key manipulation, so that you can concentrate
// on implementing the DHT structure itself.
package keys

import (
	"crypto/sha1"

	"cse586.kdht/api/kdht"
)

// Compute computes and returns the key for obj
func Compute(obj []byte) []byte {
	k := sha1.Sum(obj)
	return k[:]
}

// GetBit returns the value of the given bit in a key (or distance).
//
// This function returns a meaningless value if b >= kdht.KeyBits or
// if len(x) != kdht.KeyBytes.
func GetBit(x []byte, b int) int {
	if b >= kdht.KeyBits || len(x) != kdht.KeyBytes {
		return 0
	}
	// Bit 0 is in x[kdht.KeyBytes-1]
	v := x[kdht.KeyBytes-1-(b/8)]
	if v&byte(1<<(b%8)) == 0 {
		return 0
	} else {
		return 1
	}
}

// Distance computes the distance between two keys x and y
//
// This function may crash if x and y are not valid keys.
func Distance(x []byte, y []byte) (d [kdht.KeyBytes]byte) {
	for i := 0; i < kdht.KeyBytes; i++ {
		d[i] = x[i] ^ y[i]
	}
	return
}

// Compute the k-bucket of the first bit difference in the distance
// between two keys.
//
// We treat the key as a big-endian integer, so the zero bit of byte
// KeyBytes - 1 of the key is the lowest-order bit.  Two keys
// differing in the most significant bit go into bucket kdht.KeyBits -
// 1, and two identical keys go into bucket 0.  This matches the
// description of k-buckets in the Kademlia paper (see Section 2.2).
//
// The job of this function is therefore to find the _highest order
// bit_ that is set in the given key distance, and return the index of
// that bit.
func DistanceBucket(d [kdht.KeyBytes]byte) int {
	for i := 0; i < kdht.KeyBytes; i++ {
		// If this byte is identical, the difference must be
		// in another bucket
		if d[i] == 0 {
			continue
		}

		b := d[i]
		for j := 7; j >= 0; j-- {
			// Remember that the MSB of this byte is the
			// _closest_ to the MSB of the total distance,
			// so we check bit j but the distance from the
			// LSB of byte i is (8 - j).
			if b&(1<<j) != 0 {
				// This is the first non-zero bit in
				// the distance.
				return kdht.KeyBits - (8*i + (8 - j))
			}
		}
	}

	// The distance between these two hashes is zero.
	return 0
}
