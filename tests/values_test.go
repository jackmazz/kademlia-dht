/*
Copyright 2021, 2023 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more information.
*/

package tests

/*
import (
	"testing"

	"cse586.kdht/api/kdht"
	"cse586.kdht/given/keys"
)

// leftCenter is an address just to the "left of center" in the
// numbering space, if 00...00 is at the far right (as in the Kademlia
// paper).
var leftCenter = []byte{0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// leftCenterNeighborLeft is an address just to the "left" (toward
// 11...11) of leftCenter in the address space.
var leftCenterNeighborLeft = []byte{0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

// leftCenterNeighborRight is an address just to the right of
// leftCenter in the address space.
var leftCenterNeighborRight = []byte{
	0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
}

// TestConsistentKeySize just makes sure that the key size expected by
// the tests is the key size configured for the k-DHT.  If you choose
// to play with different key sizes, tests that use the values in this
// file will fail.
func TestConsistentKeySize(t *testing.T) {
	h := keys.Compute([]byte(""))
	if kdht.KeyBytes != 20 || len(h) != 20 {
		t.Error("Given routing tests assume key length of 20 bytes")
	}
}
*/
