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
	"bytes"
	"errors"
	"slices"

	"cse586.kdht/api/kdht"
	"cse586.kdht/given/keys"
)

type KdmRoutingTable struct {
	local   *kdht.NodeInfo
	k       int
	buckets [][]*kdht.NodeInfo
}

func NewRoutingTable(node *kdht.NodeInfo, k int) (kdht.RoutingTable, error) {
	if node.GetId() == nil || len(node.GetId()) != kdht.KeyBytes {
		return nil, errors.New("invalid id")
	}

	if k <= 0 {
		return nil, errors.New("invalid k")
	}

	bucket := make([]*kdht.NodeInfo, 0, k)
	bucket = append(bucket, node)
	table := new(KdmRoutingTable)
	table.local = node
	table.k = k
	table.buckets = make([][]*kdht.NodeInfo, 0, kdht.KeyBits)
	table.buckets = append(table.buckets, bucket)
	return table, nil
}

func (table *KdmRoutingTable) K() int {
	return table.k
}

func (table *KdmRoutingTable) InsertNode(node *kdht.NodeInfo) {
	num := table.bucketNumber(node.GetId())
	idx, bucket := table.findBucket(num)

	for len(bucket) == cap(bucket) {
		if idx != len(table.buckets)-1 {
			return
		}

		table.split()
		idx, bucket = table.findBucket(num)
	}

	table.buckets[idx] = append(bucket, node)
}

func (table *KdmRoutingTable) RemoveNode(key []byte) error {
	if bytes.Equal(table.local.GetId(), key) {
		return kdht.InvalidNodeError
	}

	idx1, idx2 := table.findKey(key)

	if idx1 == -1 || idx2 == -1 {
		return kdht.InvalidNodeError
	}

	table.buckets[idx1] = sliceDelete(table.buckets[idx1], idx2)
	return nil
}

func (table *KdmRoutingTable) Lookup(key []byte) (node *kdht.NodeInfo, ok bool) {
	idx1, idx2 := table.findKey(key)

	if idx1 == -1 || idx2 == -1 {
		return nil, false
	}

	return table.buckets[idx1][idx2], true
}

func (table *KdmRoutingTable) GetNodes(num int) []*kdht.NodeInfo {
	idx, bucket := table.getBucket(num)

	if idx == -1 || len(bucket) == 0 {
		return nil
	}

	return table.buckets[idx]
}

func (table *KdmRoutingTable) ClosestK(key []byte) []*kdht.NodeInfo {
	closest := make([]*kdht.NodeInfo, 0, table.k)

	cmp := func(node1, node2 *kdht.NodeInfo) int {
		dist1 := keys.Distance(key, node1.GetId())
		dist2 := keys.Distance(key, node2.GetId())
		return bytes.Compare(dist1[:], dist2[:])
	}

	idx1 := -1
	init := false

	for {
		if key == nil {
			return closest
		}

		num := table.bucketNumber(key)
		idx2, bucket := table.findBucket(num)

		if init && idx1 == idx2 {
			return closest
		}
		idx1 = idx2
		init = true

		if len(closest)+len(bucket) <= cap(closest) {
			closest = append(closest, bucket...)
			key = nextKey(key)
		} else {
			sorted := cloneSlice(bucket)
			slices.SortFunc(sorted, cmp)
			idx := cap(closest) - len(closest)
			closest = append(closest, sorted[:idx]...)
			return closest
		}
	}
}

func (table *KdmRoutingTable) Buckets() int {
	return len(table.buckets)
}

func (table *KdmRoutingTable) split() {
	idx := len(table.buckets) - 1
	bucket := cloneSlice(table.buckets[idx])
	left := make([]*kdht.NodeInfo, 0, table.k)
	right := make([]*kdht.NodeInfo, 0, table.k)
	table.buckets[idx] = left
	table.buckets = append(table.buckets, right)

	for _, node := range bucket {
		num := table.bucketNumber(node.GetId())
		idx, _ := table.findBucket(num)

		table.buckets[idx] = append(table.buckets[idx], node)
	}
}

func (table *KdmRoutingTable) bucketNumber(key []byte) int {
	dist := keys.Distance(table.local.GetId(), key)
	num := keys.DistanceBucket(dist)
	return num
}

func (table *KdmRoutingTable) getBucket(num int) (int, []*kdht.NodeInfo) {
	idx := cap(table.buckets) - num - 1

	if idx < 0 || idx >= len(table.buckets) {
		return -1, nil
	}

	return idx, table.buckets[idx]
}

func (table *KdmRoutingTable) findBucket(num int) (int, []*kdht.NodeInfo) {
	idx, bucket := table.getBucket(num)

	if idx == -1 {
		idx = len(table.buckets) - 1
		return idx, table.buckets[idx]
	}

	return idx, bucket
}

func (table *KdmRoutingTable) findKey(key []byte) (int, int) {
	num := table.bucketNumber(key)
	idx1, bucket := table.findBucket(num)

	for idx2, node := range bucket {
		if bytes.Equal(node.GetId(), key) {
			return idx1, idx2
		}
	}

	return -1, -1
}

func cloneSlice[Type any](slice []Type) []Type {
	clone := make([]Type, len(slice), cap(slice))
	copy(clone, slice)
	return clone
}

func sliceDelete[Type any](slice []Type, idx int) []Type {
	left := slice[0:idx]
	right := slice[idx+1:]
	return append(left, right...)
}

func nextKey(key []byte) []byte {
	next := cloneSlice(key)

	for idx, b := range next {
		if b == 0x0 {
			continue
		}

		for count := 0; count < 8; count++ {
			mask := byte(1 << (7 - count))

			if b&mask>>(7-count) == 1 {
				next[idx] = b &^ mask
				return next
			}
		}
	}

	return nil
}
