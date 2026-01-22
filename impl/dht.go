/*
Copyright 2021, 2023 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requestuires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more information.
*/

package impl

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"slices"
	"sync"

	"cse586.kdht/api/kdht"
	"cse586.kdht/given/keys"
	"google.golang.org/protobuf/proto"
)

type KdmNode struct {
	info         *kdht.NodeInfo
	alpha        int
	routingTable kdht.RoutingTable
	localStorage map[string][]byte
	storageMutex *sync.Mutex
	listener     net.Listener
	closed       bool
}

const network = "tcp"
const headerLength = 2

// NewNode returns an instance of a Node that is fully prepared to
// particpate in a k-DHT and serve requestuests.  The created node MUST be
// listening on the specified address before this method returns.  It
// MUST NOT wait for any other node to be contacted or the routing
// table to be populated before returning.
//
// The k and alpha parameters represent k and alpha in the Kademlia
// paper.  They should be used when creating your routing table,
// maintaining the k-buckets of the routing table, and contacting
// nodes for lookup.
//
// Each of the specified neighbors should be contacted and added to
// the node's routing table if they respond.  The node should then
// execute a FindNode for itself on each of these neighbors in order
// to pre-populate its routing table.
//
// This function returns an error if the new node cannot be created

func NewNode(key []byte, addr string, k int, alpha int, neighbors []string) (*KdmNode, error) {
	if alpha < 1 {
		return nil, errors.New("invalid alpha")
	}

	ln, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	info := new(kdht.NodeInfo)
	info.Id = key
	info.Address = addr

	table, err := NewRoutingTable(info, k)
	if err != nil {
		return nil, err
	}

	node := KdmNode{}
	node.info = info
	node.alpha = alpha
	node.routingTable = table
	node.localStorage = make(map[string][]byte)
	node.storageMutex = &sync.Mutex{}
	node.listener = ln
	node.closed = false

	go node.listenForRequests()
	for _, neighbor := range neighbors {
		go func(addr string) {
			request := kdht.Message{}
			request.Sender = node.info
			request.Type = kdht.MessageType_PING
			node.contactAddress(&request, addr)
		}(neighbor)
	}

	return &node, nil
}

func (node *KdmNode) Ping(id []byte, message []byte) error {
	target, ok := node.routingTable.Lookup(id)
	if !ok {
		return kdht.InvalidNodeError
	}

	request := kdht.Message{}
	request.Sender = node.info
	request.Type = kdht.MessageType_PING
	request.Value = message

	_, err := node.contactAddress(&request, target.Address)
	if err != nil {
		return err
	}
	return nil
}

func (node *KdmNode) Store(val []byte) error {
	if node.closed {
		return kdht.ShutdownError
	}

	key := keys.Compute(val)
	_, _, closest := node.nodeLookup(key, false)

	for _, info := range closest {
		go func(info *kdht.NodeInfo) {
			request := kdht.Message{}
			request.Sender = node.info
			request.Type = kdht.MessageType_STORE
			request.Key = key
			request.Value = val
			node.contactAddress(&request, info.Address)
		}(info)
	}

	if len(closest) < node.routingTable.K() {
		return kdht.StorageError
	}
	return nil
}

func (node *KdmNode) FindNode(id []byte) ([]*kdht.NodeInfo, error) {
	if node.closed {
		return nil, kdht.ShutdownError
	}

	_, _, closest := node.nodeLookup(id, false)
	return closest, nil
}

func (node *KdmNode) FindValue(id []byte) ([]byte, kdht.NodeInfo, error) {
	if node.closed {
		return nil, kdht.NodeInfo{}, kdht.ShutdownError
	}

	val, sender, _ := node.nodeLookup(id, true)
	if val == nil {
		return nil, kdht.NodeInfo{}, kdht.ValueError
	}

	return val, *sender, nil
}

func (node *KdmNode) Shutdown() error {
	if node.closed {
		return kdht.ShutdownError
	}

	node.closed = true
	return node.listener.Close()
}

func (node *KdmNode) Neighbors() []*kdht.NodeInfo {
	if node.closed {
		return []*kdht.NodeInfo{}
	}

	neighbors := []*kdht.NodeInfo{}
	maxbucket := kdht.KeyBits - 1
	minbucket := maxbucket - node.routingTable.Buckets()
	for i := minbucket; i <= maxbucket; i++ {
		bucket := node.routingTable.GetNodes(i)
		if bucket == nil {
			continue
		}

		for _, info := range bucket {
			if bytes.Equal(info.Id, node.info.Id) {
				continue
			}
			neighbors = append(neighbors, info)
		}
	}
	return neighbors
}

func (node *KdmNode) processPing(message *kdht.Message, conn net.Conn) {
	response := kdht.Message{}
	response.Sender = node.info
	response.Type = kdht.MessageType_ACK
	response.Value = message.Value
	node.sendMessage(&response, conn)
}

func (node *KdmNode) processStore(message *kdht.Message, conn net.Conn) {
	node.storeValue(message.Key, message.Value)
	response := kdht.Message{}
	response.Sender = node.info
	response.Type = kdht.MessageType_ACK
	node.sendMessage(&response, conn)
}

func (node *KdmNode) processGet(message *kdht.Message, conn net.Conn) {
	response := kdht.Message{}
	response.Sender = node.info
	val, ok := node.accessValue(message.Key)
	if !ok {
		response.Type = kdht.MessageType_ACK
	} else {
		response.Type = kdht.MessageType_VALUE
		response.Value = val
	}
	node.sendMessage(&response, conn)
}

func (node *KdmNode) processFindNode(message *kdht.Message, conn net.Conn) {
	response := kdht.Message{}
	response.Sender = node.info
	response.Type = kdht.MessageType_NODES
	response.Nodes = node.routingTable.ClosestK(message.Key)
	node.sendMessage(&response, conn)
}

func (node *KdmNode) processFindValue(message *kdht.Message, conn net.Conn) {
	val, ok := node.accessValue(message.Key)
	if !ok {
		node.processFindNode(message, conn)
		return
	}

	response := kdht.Message{}
	response.Sender = node.info
	response.Type = kdht.MessageType_VALUE
	response.Value = val
	node.sendMessage(&response, conn)
}

func (node *KdmNode) listenForRequests() {
	for {
		conn, err := node.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err == nil {
			go func(conn net.Conn) {
				defer func() { recover() }()

				for {
					message, err := node.recieveMessage(conn)
					if err != nil {
						return
					}

					go node.routingTable.InsertNode(message.Sender)
					switch message.Type {
					case kdht.MessageType_PING:
						go node.processPing(message, conn)
					case kdht.MessageType_STORE:
						go node.processStore(message, conn)
					case kdht.MessageType_GET:
						go node.processGet(message, conn)
					case kdht.MessageType_FIND_NODE:
						go node.processFindNode(message, conn)
					case kdht.MessageType_FIND_VALUE:
						go node.processFindValue(message, conn)
					}
				}
			}(conn)
		}
	}
}

func (node *KdmNode) nodeLookup(id []byte, tog bool) ([]byte, *kdht.NodeInfo, []*kdht.NodeInfo) {
	closest := node.routingTable.ClosestK(id)
	visited := make(map[string]bool)
	visited[string(node.info.Id)] = true
	vmut := &sync.Mutex{}

	markVisited := func(id []byte) {
		vmut.Lock()
		visited[string(id)] = true
		vmut.Unlock()
	}

	isVisited := func(id []byte) bool {
		vmut.Lock()
		_, ok := visited[string(id)]
		vmut.Unlock()
		return ok
	}

	var contactClosest func([]*kdht.NodeInfo, int) ([]byte, *kdht.NodeInfo, []*kdht.NodeInfo)
	contactClosest = func(slice []*kdht.NodeInfo, idx int) ([]byte, *kdht.NodeInfo, []*kdht.NodeInfo) {
		if len(slice) == 0 {
			return nil, nil, closest
		}

		ch := make(chan *kdht.Message)
		for _, info := range slice {
			go func(info *kdht.NodeInfo) {
				defer recover()

				if isVisited(info.Id) {
					ch <- nil
					return
				}
				markVisited(info.Id)

				request := kdht.Message{}
				request.Sender = node.info
				if tog {
					request.Type = kdht.MessageType_FIND_VALUE
				} else {
					request.Type = kdht.MessageType_FIND_NODE
				}
				request.Key = id

				response, err := node.contactAddress(&request, info.Address)
				if err != nil {
					ch <- nil
					return
				}

				ch <- response
			}(info)
		}

		flag := false
		count := 0
		for response := range ch {
			if response != nil {
				if tog && response.Value != nil {
					close(ch)
					return response.Value, response.Sender, nil
				}

				for _, info := range response.Nodes {
					if containsNode(closest, info) {
						continue
					}

					if len(closest) < node.routingTable.K() {
						closest = append(closest, info)
						flag = true
						continue
					}

					dist := sliceOfDistance(info.Id, id)
					maxidx, maxdist := furthestNode(closest, id)
					if bytes.Compare(dist, maxdist) < 0 {
						closest[maxidx] = info
						flag = true
					}
				}
			}

			count++
			if count >= len(slice) {
				close(ch)
				break
			}
		}

		if !flag {
			return contactClosest(sliceAtMost(closest[idx+1:], node.alpha), idx+1)
		}
		return contactClosest(sliceAtMost(closest, node.alpha), 0)
	}
	return contactClosest(sliceAtMost(closest, node.alpha), 0)
}

func (node *KdmNode) storeValue(key []byte, val []byte) {
	node.storageMutex.Lock()
	node.localStorage[string(key)] = val
	node.storageMutex.Unlock()
}

func (node *KdmNode) accessValue(key []byte) ([]byte, bool) {
	node.storageMutex.Lock()
	val, ok := node.localStorage[string(key)]
	node.storageMutex.Unlock()
	return val, ok
}

func (node *KdmNode) contactAddress(message *kdht.Message, addr string) (*kdht.Message, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	err = node.sendMessage(message, conn)
	defer func() { conn.Close() }()
	if err != nil {
		return nil, err
	}

	response, err := node.recieveMessage(conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (node *KdmNode) sendMessage(message *kdht.Message, conn net.Conn) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	length := len(data)
	data = slices.Concat(make([]byte, headerLength), data)
	binary.BigEndian.PutUint16(data[:headerLength], uint16(length))

	var total int
	for total < int(length) {
		nbytes, err := conn.Write(data[total:])
		if err != nil {
			return err
		}
		total += nbytes
	}

	return nil
}

func (node *KdmNode) recieveMessage(conn net.Conn) (*kdht.Message, error) {
	buf := make([]byte, headerLength)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint16(buf)
	buf = make([]byte, length)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}

	message := new(kdht.Message)
	err = proto.Unmarshal(buf, message)
	if err != nil {
		return nil, err
	}

	go node.routingTable.InsertNode(message.Sender)
	return message, nil
}

func containsNode(nodes []*kdht.NodeInfo, target *kdht.NodeInfo) bool {
	for _, info := range nodes {
		if bytes.Equal(info.Id, target.Id) {
			return true
		}
	}
	return false
}

func furthestNode(nodes []*kdht.NodeInfo, key []byte) (int, []byte) {
	var maxdist []byte
	maxidx := -1
	for idx, info := range nodes {
		dist := sliceOfDistance(info.Id, key)
		if maxdist == nil || bytes.Compare(maxdist, dist) < 0 {
			maxdist = dist
			maxidx = idx
		}
	}
	return maxidx, maxdist
}

func sliceOfDistance(key1 []byte, key2 []byte) []byte {
	dist := keys.Distance(key1, key2)
	return dist[:]
}

func sliceAtMost[T any](slice []T, maxlen int) []T {
	retslice := slice
	if maxlen <= len(slice) {
		retslice = slice[0:maxlen]
	}
	return retslice
}
