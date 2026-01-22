// Return the K closest nodes to the specified key. They may
// be in a single bucket, or may be spread across several
// buckets.
//
// This operation cannot fail and cannot return an empty slice
// because the local node is always present.
func (table *KdmRoutingTable) ClosestK(key []byte) []*kdht.NodeInfo {
	closest := make([]*kdht.NodeInfo, 0, table.k)
	id := key
	pos := 0

	for len(closest) < cap(closest) {
		num := table.bucketNumber(id)
		idx := table.bucketIndex(num)
		bucket := table.buckets[idx]

		if len(closest)+len(bucket) <= cap(closest) {
			for _, node := range bucket {
				if !bytes.Equal(node.GetId(), key) {
					closest = append(closest, node)
				}
			}
		} else {
			// sort by distance
			break
		}

		pos++
		idx2 := pos / len(key)

		// set id to next proper id
	}

	return closest
}