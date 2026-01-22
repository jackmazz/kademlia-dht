func (node KdmNode) FindNode(id []byte) ([]*kdht.NodeInfo, error) {
	visited := make(map[string]bool)
	visited[string(node.info.Id)] = true
	closest := node.closestKContacts(id)
	printInfos("initial closest nodes", closest)
	// fmt.Printf("======================\n")
	// fmt.Printf("| ITERATIONS STARTED |\n")
	// fmt.Printf("======================\n\n")
	for {
		var found []*kdht.NodeInfo
		c := 0
		for _, info := range closest {
			// fmt.Printf("contacting %v\n", info.Id)

			if c >= node.alpha {
				// fmt.Printf("already contacted alpha (%v) nodes\n\n", c)
				break
			}
			_, ok := visited[string(info.Id)]
			if ok {
				// fmt.Printf("already visited\n\n")
				continue
			}

			req := kdht.Message{}
			req.Sender = node.info
			req.Type = kdht.MessageType_FIND_NODE
			req.Key = id

			rsp, err := node.contactAddress(&req, info.Address)
			visited[string(info.Id)] = true
			if err != nil {
				// fmt.Printf("an error occured when contacting\n\n")
				continue
			}

			found = append(found, rsp.Nodes...)
			c++
			printInfos("found nodes", rsp.Nodes)
		}

		printInfos("all found nodes", found)

		flag := false
		for _, finfo := range found {
			farr := keys.Distance(finfo.Id, id)
			fdist := farr[:]

			present := false
			for _, info := range closest {
				if bytes.Equal(info.Id, finfo.Id) {
					present = true
					break
				}
			}
			if present {
				continue
			}

			if len(closest) < node.routingTable.K() {
				closest = append(closest, finfo)
				flag = true
				continue
			}

			var max []byte
			maxidx := -1
			for idx, info := range closest {
				arr := keys.Distance(info.Id, id)
				dist := arr[:]
				if max == nil || bytes.Compare(max, dist) < 0 {
					max = dist
					maxidx = idx
				}
			}

			if bytes.Compare(fdist, max) < 0 {
				closest[maxidx] = finfo
				flag = true
			}
		}

		if !flag {
			printInfos("final closest nodes", closest)
			// fmt.Printf("=======================\n")
			// fmt.Printf("| ITERATIONS COMPLETE |\n")
			// fmt.Printf("=======================\n\n")
			return closest, nil
		}

		printInfos("current closest nodes", closest)
		// fmt.Printf("==================\n")
		// fmt.Printf("| NEXT ITERATION |\n")
		// fmt.Printf("==================\n\n")
	}
}