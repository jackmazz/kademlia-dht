func (node KdmNode) processFindNode(msg *kdht.Message, conn net.Conn) {
	req := kdht.Message{}
	req.Sender = node.info
	req.Type = kdht.MessageType_FIND_NODE
	req.Key = msg.Key

	ctcs := node.routingTable.ClosestK(msg.Key)
	// reduce ctcs to size alpha
	found := make(map[string]bool)
	c := 0
	mut := sync.Mutex{}

	for _, ctc := range ctcs {
		go func(ctc *kdht.NodeInfo) {
			rsp, err := node.contactAddress(&req, ctc.Address)
			mut.Lock()
			c++
			if err == nil {
				found[string(rsp.Key)] = true
			}
			mut.Unlock()
		}(ctc)
	}

	for len(found) < node.routingTable.K() || c < node.alpha {
		// wait
	}

	// sort found's values by distance
	// reduce found's values to size k

	rsp := kdht.Message{}
	rsp.Sender = node.info
	rsp.Type = kdht.MessageType_NODES
	rsp.Nodes = nil // put found's values here
	node.sendMessage(&rsp, conn)
}