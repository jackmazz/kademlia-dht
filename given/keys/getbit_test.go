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

var emptyID = make([]byte, 20)

func TestConsistentLengthGetBit(t *testing.T) {
	if kdht.KeyBytes != 20 {
		t.Error("GetBit tests assume key length of 20 bytes")
	}
}

func TestGetLowestBit(t *testing.T) {
	if GetBit(emptyID, 0) != 0 {
		t.Error("Empty ID lowest bit was not zero")
	}

	id := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x1}
	if GetBit(id, 0) != 1 {
		t.Error("Lowest bit was not one")
	}
	if GetBit(id, 1) != 0 {
		t.Error("Second bit was not zero")
	}
}

func TestGetHighestBit(t *testing.T) {
	if GetBit(emptyID, kdht.KeyBits-1) != 0 {
		t.Error("Highest bit of empty ID was not zero")
	}

	id := []byte{0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if GetBit(id, kdht.KeyBits-1) != 1 {
		t.Error("Highest bit was not one")
	}
	if GetBit(id, kdht.KeyBits-2) != 0 {
		t.Error("Second highest bit was not zero")
	}
}

func TestEighthBit(t *testing.T) {
	if GetBit(emptyID, 8) != 0 {
		t.Error("Eighth bit of empty ID was not zero")
	}

	id := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x1, 0}
	if GetBit(id, 7) != 0 {
		t.Error("Seventh bit was not zero")
	}
	if GetBit(id, 8) != 1 {
		t.Error("Eighth bit was not one")
	}
	if GetBit(id, 9) != 0 {
		t.Error("Ninth bit was not zero")
	}
}

func TestTwentiethBit(t *testing.T) {
	if GetBit(emptyID, 20) != 0 {
		t.Error("Twentieth bit of empty ID was not zero")
	}

	id := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x10, 0, 0}
	if GetBit(id, 19) != 0 {
		t.Error("Nineteenth bit was not zero")
	}
	if GetBit(id, 20) != 1 {
		t.Error("Twentieth bit was not one")
	}
	if GetBit(id, 21) != 0 {
		t.Error("Twenty-First bit was not zero")
	}
}
