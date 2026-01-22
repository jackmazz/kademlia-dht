/*
Copyright 2021 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more package.
*/

package router

import (
	"bytes"
	"crypto/sha1"
	"testing"

	"cse586.kdht/api/kdht"
)

// We only need to provide a small number of tests here, because the
// general router testing suite can be used to test the rest of this
// functionality (through, e.g., Autograder).  These tests should test
// a couple of cases to ensure things like correct encoding of errors,
// basic connectivity to the router server, etc.

func TestRouterConnect(t *testing.T) {
	key := sha1.Sum([]byte("RADIO"))
	rt, err := New(&kdht.NodeInfo{Id: key[:], Address: ""}, 3)
	if err != nil || rt == nil {
		t.Fatalf("Could not create router object: %v", err)
	}
}

func TestRouterK(t *testing.T) {
	const k = 7

	key := sha1.Sum([]byte("Pawn Shop"))
	rt, _ := New(&kdht.NodeInfo{Id: key[:], Address: ""}, k)

	rtk := rt.K()
	if rtk != k {
		t.Errorf("K was unexpected: %d (expected %d)", rtk, k)
	}
}

func TestLookupSelf(t *testing.T) {
	key := sha1.Sum([]byte("Dogs and Chaplains"))
	rt, _ := New(&kdht.NodeInfo{Id: key[:], Address: ""}, 3)

	n, ok := rt.Lookup(key[:])
	if !ok || bytes.Compare(n.Id, key[:]) != 0 {
		t.Errorf("Could not look up self, or ID differed: %v %s", ok, n)
	}
}

func TestLookupUnknown(t *testing.T) {
	key := sha1.Sum([]byte("REQUIEM FOR A FRIEND"))
	rt, _ := New(&kdht.NodeInfo{Id: key[:], Address: ""}, 3)

	id := sha1.Sum(key[:])
	_, ok := rt.Lookup(id[:])
	if ok {
		t.Error("Lookup of invalid node succeeded?")
	}
}
