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
	"bytes"
	"testing"

	"cse586.kdht/api/kdht"
	"cse586.kdht/impl"
)

// The tests in this file test basic routing table insertions to
// ensure that nodes are inserted (or not inserted, as the case may
// be) into the routing table appropriately.  Because it is not in the
// impl package and it does not depend on the exact definition of your
// routing table, it must necessarily use other methods on your
// routing table (such as Buckets() and GetNodes()), which were placed
// there precisely to be used for testing.  You may wish to create
// simpler tests that use the internal structure of your routing table
// in order to test your table without these extra method
// dependencies.  You may use these tests as a model for that purpose.
//
// None of the tests in this package check the ordering properties of
// the routing table that are required by the Kademlia paper.  Your
// implementation is not required to order nodes within a bucket in
// any particular order.

// TestInsertAdjacentK2 creates a routing table with k=2, then inserts
// a single neighbor on the same side of the DHT as the routing table
// node.  It ensures that both nodes go into the same k-Bucket, and
// that the k-Bucket is properly filled.
func TestInsertAdjacentK2(t *testing.T) {
	n := &kdht.NodeInfo{Id: leftCenter}
	kt, error := impl.NewRoutingTable(n, 2)
	if kt == nil || error != nil {
		t.Fatal("Could not create routing table")
	}

	n2 := &kdht.NodeInfo{Id: leftCenterNeighborLeft}
	kt.InsertNode(n2)

	if kt.Buckets() != 1 {
		t.Errorf("Routing table with one neighbor and k = 2 does not have one bucket")
	}

	nodes := kt.GetNodes(159)
	if len(nodes) != 2 {
		t.Errorf("Bucket 159 does not have 2 nodes")
	}
	// The two nodes returned should be nodes n and n2
	if !bytes.Equal(nodes[0].Id, n.Id) &&
		!bytes.Equal(nodes[1].Id, n.Id) {
		t.Errorf("Node n was not in the routing table")
	}
	if !bytes.Equal(nodes[0].Id, n2.Id) &&
		!bytes.Equal(nodes[1].Id, n2.Id) {
		t.Errorf("Node n2 was not in the routing table")
	}
}

// TestInsertTwoK2 creates a routing table with k=2 and then inserts
// one node on the same side of the center of the key space as the
// routing table node, and then one node on the opposite side of the
// key space.  This should result in the following routing table:
//
// | bucket 159               | bucket 158               |
// | nr                       | n nl                     |
//
// This test does NOT actually make sure that the ordering
// requirements in the Kademlia paper are maintained, as they are not
// required for this project.
func TestInsertTwoK2(t *testing.T) {
	n := &kdht.NodeInfo{Id: leftCenter}
	nl := &kdht.NodeInfo{Id: leftCenterNeighborLeft}
	nr := &kdht.NodeInfo{Id: leftCenterNeighborRight}

	kt, error := impl.NewRoutingTable(n, 2)
	if kt == nil || error != nil {
		t.Fatal("Could not create routing table")
	}
	kt.InsertNode(nl)
	kt.InsertNode(nr)

	if kt.Buckets() != 2 {
		t.Errorf("Routing table with two neighbors and k = 2 does not have two buckets")
	}

	nodes := kt.GetNodes(159)
	if len(nodes) != 1 {
		t.Errorf("Bucket 0 does not have 1 node")
	}
	if !bytes.Equal(nodes[0].Id, nr.Id) {
		t.Errorf("Node nr was not in the routing table at level 0")
	}

	nodes = kt.GetNodes(158)
	if len(nodes) != 2 {
		t.Errorf("Bucket 1 does not have 2 nodes")
	}
	// The two nodes returned should be nodes n and n2
	if !bytes.Equal(nodes[0].Id, n.Id) &&
		!bytes.Equal(nodes[1].Id, n.Id) {
		t.Errorf("Node n was not in the routing table at level 1")
	}
	if !bytes.Equal(nodes[0].Id, nl.Id) &&
		!bytes.Equal(nodes[1].Id, nl.Id) {
		t.Errorf("Node nr was not in the routing table at level 1")
	}
}

// TestInsertTwoFarSideK2 creates a routing table with k=2 and then
// inserts two nodes on the far side of the routing table from the
// routing table node.  This differs from TestInsertTwoK2 in where the
// resulting nodes are placed in the table.  It should result in this
// table:
//
// | bucket 159               | bucket 158               |
// | n1 n2                    | n                        |
func TestInsertTwoFarSideK2(t *testing.T) {
	n := &kdht.NodeInfo{Id: leftCenter}
	kt, error := impl.NewRoutingTable(n, 2)
	if kt == nil || error != nil {
		t.Fatal("Could not create routing table")
	}

	id1 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	id2 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	n1 := &kdht.NodeInfo{Id: id1}
	n2 := &kdht.NodeInfo{Id: id2}

	kt.InsertNode(n1)
	kt.InsertNode(n2)

	if kt.Buckets() != 2 {
		t.Errorf("Routing table does not have two buckets after inserting three far nodes with k = 2")
	}

	nodes := kt.GetNodes(159)
	if len(nodes) != 2 {
		t.Errorf("Far side of routing table does not have two nodes with k = 2")
	}
	// The two nodes returned should be nodes n1 and n2
	if !bytes.Equal(nodes[0].Id, n1.Id) &&
		!bytes.Equal(nodes[1].Id, n1.Id) {
		t.Errorf("Node n1 was not in the routing table at level 1")
	}
	if !bytes.Equal(nodes[0].Id, n2.Id) &&
		!bytes.Equal(nodes[1].Id, n2.Id) {
		t.Errorf("Node n2 was not in the routing table at level 1")
	}

	nodes = kt.GetNodes(158)
	if len(nodes) != 1 {
		t.Errorf("Near side of routing table does not have one node: %d", len(nodes))
	}
	if !bytes.Equal(nodes[0].Id, n.Id) {
		t.Errorf("Node n is not the only node on the near side")
	}
}

// TestInsertDiscardFarSideK2 creates a routing table with k=2 and
// then inserts three nodes on the far side of the key space from the
// routing table node.  This should cause only the FIRST TWO of those
// nodes to be stored in the routing table, as the third node will be
// discarded since the k-bucket into which it would otherwise be
// inserted already has k entries.  The resulting routing table is the
// same as the routing table for TestInsertTwoFarSideK2.
func TestInsertDiscardFarSideK2(t *testing.T) {
	n := &kdht.NodeInfo{Id: leftCenter}
	kt, error := impl.NewRoutingTable(n, 2)
	if kt == nil || error != nil {
		t.Fatal("Could not create routing table")
	}
	id1 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	id2 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	id3 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3}
	n1 := &kdht.NodeInfo{Id: id1}
	n2 := &kdht.NodeInfo{Id: id2}
	n3 := &kdht.NodeInfo{Id: id3}

	kt.InsertNode(n1)
	kt.InsertNode(n2)
	kt.InsertNode(n3)

	if kt.Buckets() != 2 {
		t.Errorf("Routing table does not have two buckets after inserting three far nodes with k = 2")
	}

	nodes := kt.GetNodes(159)
	if len(nodes) != 2 {
		t.Errorf("Far side of routing table does not have two nodes with k = 2")
	}
	// The two nodes returned should be nodes n1 and n2
	if !bytes.Equal(nodes[0].Id, n1.Id) &&
		!bytes.Equal(nodes[1].Id, n1.Id) {
		t.Errorf("Node n1 was not in the routing table at level 1")
	}
	if !bytes.Equal(nodes[0].Id, n2.Id) &&
		!bytes.Equal(nodes[1].Id, n2.Id) {
		t.Errorf("Node n2 was not in the routing table at level 1")
	}

	nodes = kt.GetNodes(158)
	if len(nodes) != 1 {
		t.Errorf("Near side of routing table does not have one node: %d", len(nodes))
	}
}
*/
