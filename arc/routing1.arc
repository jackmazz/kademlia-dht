/*
Copyright 2021 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more package.
*/

package impl

import (
	"errors"

	"cse586.kdht/api/kdht"
)

type KdmRoutingTable struct {
	localId []byte
	root    *KdmRoutingNode
	k       int
	nbkts   int
}

type KdmRoutingNode struct {
	bitidx int
	left   *KdmRoutingNode
	right  *KdmRoutingNode
	bkt    []*kdht.NodeInfo
	cap    int
}

// NewRoutingTable returns an instance of a RoutingTable suitable for
// the node identified by the specified NodeInfo with the specified
// parameter k for its k-buckets.
//
// This function returns an error if the provided NodeInfo is invalid
// (e.g., its key is nil or not kdht.KeyBytes long), or if the
// specified k is <= 0.
func NewRoutingTable(node *kdht.NodeInfo, k int) (kdht.RoutingTable, error) {
	if node.GetId() == nil || len(node.GetId()) != kdht.KeyBytes {
		return nil, errors.New("invalid id")
	}

	if k <= 0 {
		return nil, errors.New("invalid k")
	}

	tbl := new(KdmRoutingTable)
	tbl.localId = node.GetId()
	tbl.root = newRoutingNode(node, 0, k)
	tbl.k = k
	tbl.nbkts = 1
	return tbl, nil
}

// K returns the value of K used by this routing table.
func (tbl *KdmRoutingTable) K() int {
	return tbl.k
}

// InsertNode adds the node having the specified information
// to the routing table.  This operation cannot fail, although
// the node may not be inserted if the corresponding k-bucket
// is already full.
func (tbl *KdmRoutingTable) InsertNode(node *kdht.NodeInfo) {
	tbl.insertNodeAux(node, tbl.root)
}

// RemoveNode removes the node having the specified
// identifier from the routing table.
//
// This operation fails and returns InvalidNodeError if the
// node is not present in the table or local node's ID is
// attempted to be removed.
func (tbl *KdmRoutingTable) RemoveNode(key []byte) error {
	return nil
}

// Lookup finds a particular node in the local routing table
// by its ID, if it exists.  If it does not exist, nothing
// happens and ok is false.
func (tbl *KdmRoutingTable) Lookup(key []byte) (node *kdht.NodeInfo, ok bool) {
	return nil, false
}

// Retrieve the nodes stored in a numbered bucket.  Bucket 0
// is the bucket representing nodes that differ only in the
// least-significant bit, and bucket 159 is the bucket
// representing nodes that differ in the most-significant bit.
// This is the same notation used in the Kademlia paper, and
// the same behavior provided by given/keys.DistanceBucket().
//
// This operation returns nil if the specified bucket is empty.
func (tbl *KdmRoutingTable) GetNodes(bucket int) []*kdht.NodeInfo {
	rnode := tbl.root
	for i := kdht.KeyBits - 1; i > bucket && rnode.right != nil; i-- {
		rnode = rnode.right
	}
	if rnode.left != nil {
		return rnode.left.bkt
	}
	return rnode.bkt
}

// Return the K closest nodes to the specified key.  They may
// be in a single bucket, or may be spread across several
// buckets.
//
// This operation cannot fail and cannot return an empty slice
// because the local node is always present.
func (tbl *KdmRoutingTable) ClosestK(key []byte) []*kdht.NodeInfo {
	return nil
}

// Buckets returns the number of non-empty buckets in this
// routing table.  These buckets necessarily start with bucket
// number 159 and work their way down to 160 - Buckets().  The
// return value of this function is in 0 <= Buckets() < 160.
func (tbl *KdmRoutingTable) Buckets() int {
	return tbl.nbkts
}

func (tbl *KdmRoutingTable) insertNodeAux(node *kdht.NodeInfo, start *KdmRoutingNode) {
	dist := xorBytes(tbl.localId, node.GetId())
	rnode := start
	for {
		if rnode.bkt != nil {
			if !rnode.isFull() {
				rnode.bkt = append(rnode.bkt, node)
				return
			}
			tbl.split(rnode)
			tbl.nbkts++
		} else {
			bit := bitAt(dist, rnode.bitidx)
			if bit == 0x0 {
				rnode = rnode.right
			} else if bit == 0x1 && !rnode.left.isFull() {
				rnode = rnode.left
			} else {
				return
			}
		}
	}
}

func (tbl *KdmRoutingTable) split(rnode *KdmRoutingNode) {
	bitidx := rnode.bitidx + 1
	if bitidx >= kdht.KeyBits {
		return
	}

	bkt := rnode.bkt
	rnode.left = newRoutingNode(nil, bitidx, rnode.cap)
	rnode.right = newRoutingNode(nil, bitidx, rnode.cap)
	rnode.bkt = nil
	for i := 0; i < len(bkt); i++ {
		tbl.insertNodeAux(bkt[i], rnode)
	}
}

func newRoutingNode(node *kdht.NodeInfo, bitidx, cap int) *KdmRoutingNode {
	rnode := new(KdmRoutingNode)
	rnode.bitidx = bitidx
	rnode.left = nil
	rnode.right = nil
	rnode.bkt = make([]*kdht.NodeInfo, 0)
	rnode.cap = cap
	if node != nil {
		rnode.bkt = append(rnode.bkt, node)
	}
	return rnode
}

func (rnode *KdmRoutingNode) isFull() bool {
	return len(rnode.bkt) == rnode.cap
}

func xorBytes(arr1, arr2 []byte) []byte {
	xorarr := make([]byte, len(arr1))
	for i := 0; i < len(arr1); i++ {
		xorarr[i] = arr1[i] ^ arr2[i]
	}
	return xorarr
}

func bitAt(arr []byte, idx int) byte {
	b := arr[idx/8]
	bits := make([]byte, 8)
	bits[0] = b & 0x80 >> 7
	bits[1] = b & 0x40 >> 6
	bits[2] = b & 0x20 >> 5
	bits[3] = b & 0x10 >> 4
	bits[4] = b & 0x08 >> 3
	bits[5] = b & 0x04 >> 2
	bits[6] = b & 0x02 >> 1
	bits[7] = b & 0x01
	return bits[idx%8]
}
