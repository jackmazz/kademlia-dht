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
	"bytes"
	"crypto/sha1"
	"testing"

	"cse586.kdht/api/kdht"
)

func TestConsistentLengthDistance(t *testing.T) {
	if kdht.KeyBytes != 20 {
		t.Error("Distance tests assume key length of 20 bytes")
	}
}

func TestZeroDistance(t *testing.T) {
	d := Distance(emptyID, emptyID)

	if bytes.Compare(d[:], emptyID) != 0 {
		t.Error("Comparison of two empty IDs yielded nonzero distance")
	}

	id := sha1.Sum([]byte("TestZeroDistance"))
	d = Distance(id[:], id[:])
	if bytes.Compare(d[:], emptyID) != 0 {
		t.Error("Comparison of an arbitrary ID to itself yielded nonzero distance")
	}
}

func TestNonzeroDistance(t *testing.T) {
	id1 := sha1.Sum([]byte("TestNonzeroDistance"))
	id2 := sha1.Sum([]byte("Another arbitrary key"))

	d := Distance(id1[:], emptyID)
	if bytes.Compare(d[:], emptyID) == 0 {
		t.Error("An arbitary key had zero distance to the empty key")
	}

	d = Distance(id2[:], emptyID)
	if bytes.Compare(d[:], emptyID) == 0 {
		t.Error("An arbitary key had zero distance to the empty key")
	}

	d = Distance(id1[:], id2[:])
	if bytes.Compare(d[:], emptyID) == 0 {
		t.Error("Two different keys had zero distance to each other")
	}

}

func TestKnownDistance(t *testing.T) {
	id1 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x1}
	id2 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x2}
	d1_2 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x3}

	d := Distance(id1, id2)
	if bytes.Compare(d[:], d1_2) != 0 {
		t.Error("Known distance between two one-bit IDs was incorrect")
	}

	idA := []byte{
		0x9f, 0x2c, 0xbb, 0xcb, 0xa1, 0x05, 0xf1, 0x39, 0x0d, 0xb7,
		0x9a, 0x26, 0x3a, 0xf5, 0xbf, 0x95, 0x54, 0xf9, 0x35, 0xa4,
	}
	idB := []byte{
		0xb6, 0xf5, 0xeb, 0x7a, 0x99, 0xed, 0xf5, 0x29, 0x70, 0x43,
		0x9b, 0x82, 0x65, 0x60, 0xa9, 0x8c, 0x06, 0xca, 0xd4, 0x1f,
	}
	dA_B := []byte{
		0x29, 0xd9, 0x50, 0xb1, 0x38, 0xe8, 0x04, 0x10, 0x7d, 0xf4,
		0x01, 0xa4, 0x5f, 0x95, 0x16, 0x19, 0x52, 0x33, 0xe1, 0xbb,
	}
	d = Distance(idA, idB)
	if bytes.Compare(d[:], dA_B) != 0 {
		t.Errorf("Known distance between two arbitrary IDs was incorrect")
	}
}
