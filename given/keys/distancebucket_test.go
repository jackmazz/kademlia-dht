/*
Copyright 2021 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more package.
*/

package keys

import (
	"testing"

	"cse586.kdht/api/kdht"
)

func TestConsistentLengthDistanceBucket(t *testing.T) {
	if kdht.KeyBytes != 20 {
		t.Error("KeyDistance tests assume key length of 20 bytes")
	}
}

func TestDistanceBucketEmptyID(t *testing.T) {
	b := DistanceBucket([20]byte{})
	if b != 0 {
		t.Error("Distance of empty difference was not zero")
	}
}

func TestDistanceBucketMSB(t *testing.T) {
	b := DistanceBucket([20]byte{0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	if b != 159 {
		t.Error("Distance of most significant bit is not 159")
	}
}

func TestDistanceBucket79(t *testing.T) {
	b := DistanceBucket([20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	if b != 79 {
		t.Errorf("Distance of bit 79 is not 79: %d", b)
	}
}

func TestDistanceBucket80(t *testing.T) {
	b := DistanceBucket([20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	if b != 80 {
		t.Errorf("Distance of bit 80 is not 80: %d", b)
	}
}

func TestArbitraryBucket(t *testing.T) {
	b := DistanceBucket([20]byte{
		0x29, 0xd9, 0x50, 0xb1, 0x38, 0xe8, 0x04, 0x10, 0x7d, 0xf4,
		0x01, 0xa4, 0x5f, 0x95, 0x16, 0x19, 0x52, 0x33, 0xe1, 0xbb,
	})
	if b != 157 {
		t.Errorf("Distance bucket of arbitrary value is not 157: %d", b)
	}
}
