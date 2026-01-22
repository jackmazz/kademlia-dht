func (node KdmNode) FindNode(id []byte) ([]*kdht.NodeInfo, error) {
	visited := make(map[string]bool)
	visited[string(node.info.Id)] = true
	closest := node.closestKContacts(id)
	vmut := &sync.Mutex{}

	printInfos("initial closest nodes", closest)
	fmt.Printf("======================\n")
	fmt.Printf("| ITERATIONS STARTED |\n")
	fmt.Printf("======================\n\n")

	for {
		ch := make(chan *kdht.Message)

		lim := node.alpha
		if lim > len(closest) {
			lim = len(closest)
		}
		count := 0
		cmut := &sync.Mutex{}

		for idx := 0; idx < lim; idx++ {
			go func(info *kdht.NodeInfo) {
				fmt.Printf("contacting %v\n\n", info.Id)

				defer func() {
					recover()
					cmut.Lock()
					count++
					cmut.Unlock()
					fmt.Printf("done contacting %v\n\n", info.Id)
				}()

				vmut.Lock()
				_, ok := visited[string(info.Id)]
				if ok {
					fmt.Printf("already visited %v\n\n", info.Id)
					vmut.Unlock()
					return
				}
				visited[string(info.Id)] = true
				vmut.Unlock()

				req := kdht.Message{}
				req.Sender = node.info
				req.Type = kdht.MessageType_FIND_NODE
				req.Key = id

				rsp, err := node.contactAddress(&req, info.Address)
				if err != nil {
					fmt.Printf("couldn't contact %v\n\n", info.Id)
					return
				}

				ch <- rsp
				fmt.Printf("got a response from %v\n\n", info.Id)
			}(closest[idx])
		}
		ch <- nil

		flag := false
		for rsp := range ch {
			if rsp != nil {
				printInfos(fmt.Sprintf("recieved from %v", rsp.Sender.Id), rsp.Nodes)
				for _, info := range rsp.Nodes {
					if infosContains(closest, info) {
						continue
					}

					if len(closest) < node.routingTable.K() {
						closest = append(closest, info)
						flag = true
						continue
					}

					idx, max := infosMaxDistance(closest, id)
					dist := distanceAsSlice(info.Id, id)
					if bytes.Compare(dist, max) < 0 {
						closest[idx] = info
						flag = true
					}
				}
			}

			if flag {
				close(ch)
				break
			}

			cmut.Lock()
			if count == lim {
				close(ch)
				printInfos("final closest nodes", closest)
				fmt.Printf("=======================\n")
				fmt.Printf("| ITERATIONS COMPLETE |\n")
				fmt.Printf("=======================\n\n")
				return closest, nil
			}
			cmut.Unlock()
		}

		printInfos("current closest nodes", closest)
		fmt.Printf("==================\n")
		fmt.Printf("| NEXT ITERATION |\n")
		fmt.Printf("==================\n\n")
	}
}