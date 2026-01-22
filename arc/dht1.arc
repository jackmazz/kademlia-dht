/*
Copyright 2021, 2023 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more information.
*/

package impl

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"slices"

	"cse586.kdht/api/kdht"
	"google.golang.org/protobuf/proto"
)

type KdmNode struct {
	info            *kdht.NodeInfo
	routingTable    kdht.RoutingTable
	alpha           int
	listener        net.Listener
	requestChannel  chan *kdht.Message
	responseChannel chan *kdht.Message
}

const Network = "tcp"
const HeaderLength = 2

// NewNode returns an instance of a Node that is fully prepared to
// particpate in a k-DHT and serve requests.  The created node MUST be
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
func NewNode(key []byte, addr string, k int, alpha int, neighbors []string) (kdht.Node, error) {
	ln, err := net.Listen(Network, addr)
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

	node := new(KdmNode)
	node.info = info
	node.routingTable = table
	node.alpha = alpha
	node.listener = ln
	node.requestChannel = make(chan *kdht.Message)

	go node.processMessages()
	go node.listenForMessages()

	return node, nil
}

func (node *KdmNode) Ping(id []byte, value []byte) error {
	tgt, ok := node.routingTable.Lookup(id)
	if !ok {
		return kdht.InvalidNodeError
	}

	msg := newMessage(node.info, kdht.MessageType_PING)
	msg.Value = value
	err := sendMessage(msg, tgt.Address)
	if err != nil {
		return err
	}

}

func (node *KdmNode) Store(value []byte) error {
	return errors.New("not implemented")
}

func (node *KdmNode) FindNode(id []byte) ([]*kdht.NodeInfo, error) {
	return nil, errors.New("not implemented")
}

func (node *KdmNode) FindValue(id []byte) ([]byte, kdht.NodeInfo, error) {
	return nil, kdht.NodeInfo{}, errors.New("not implemented")
}

func (node *KdmNode) Shutdown() error {
	return errors.New("not implemented")
}

func (node *KdmNode) Neighbors() []*kdht.NodeInfo {
	return nil
}

func (node *KdmNode) processMessages() {
	handlePing := func(msg *kdht.Message) {
		_, ok := node.routingTable.Lookup(msg.Sender.Id)
		if !ok {
			node.routingTable.InsertNode(msg.Sender)
		}
		response := newMessage(msg.Sender, kdht.MessageType_ACK)
		response.Value = msg.Value
		sendMessage(response, msg.Sender.Address)
	}

	for msg := range node.requestChannel {
		switch msg.Type {
		case kdht.MessageType_PING:
			handlePing(msg)
		}
	}
}

func (node *KdmNode) listenForMessages() {
	for {
		conn, err := node.listener.Accept()
		if err == net.ErrClosed {
			return
		}

		if err == nil {
			go func(conn net.Conn) {
				defer func() { recover() }()

				for {
					buf := make([]byte, HeaderLength)
					_, err := io.ReadFull(conn, buf)
					if err != nil {
						return
					}

					l := binary.BigEndian.Uint16(buf)
					buf = make([]byte, l)
					_, err = io.ReadFull(conn, buf)
					if err != nil {
						return
					}

					msg := new(kdht.Message)
					err = proto.Unmarshal(buf, msg)
					if err != nil {
						return
					}

					switch msg.Type {
					case kdht.MessageType_PING:
					case kdht.MessageType_STORE:
					case kdht.MessageType_FIND_NODE:
					case kdht.MessageType_FIND_VALUE:
						node.requestChannel <- msg
					case kdht.MessageType_ACK:
					case kdht.MessageType_GET:
					case kdht.MessageType_VALUE:
					case kdht.MessageType_NODES:
						node.responseChannel <- msg
					}
				}
			}(conn)
		}
	}
}

func newMessage(sender *kdht.NodeInfo, msgType kdht.MessageType) *kdht.Message {
	msg := new(kdht.Message)
	msg.Sender = sender
	msg.Type = msgType
	msg.Key = nil
	msg.Value = nil
	msg.Nodes = nil
	return msg
}

func sendMessage(msg *kdht.Message, addr string) error {
	conn, err := net.Dial(Network, addr)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	l := len(data)
	data = slices.Concat(make([]byte, HeaderLength), data)
	binary.BigEndian.PutUint16(data[:HeaderLength], uint16(l))

	var wl int
	for wl < int(l) {
		n, err := conn.Write(data[wl:])
		if err != nil {
			return err
		}
		wl += n
	}

	conn.Close()
	return nil
}
