/*
Copyright 2021, 2023 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more information.
*/

package router

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"

	"cse586.kdht/api/kdht"
	"google.golang.org/protobuf/proto"
)

// routerCommand is the command to use to run the routing table
// server.  Make sure that kdht-router is in your PATH so that this
// will work.
const routerCommand = "kdht-router"

// connectTimeout is the time to wait for the server to start up and
// connect to a listening socket.  If the server does not reply within
// this period of time, the request failed, so abort and return an
// error.
const connectTimeout = 300000 * time.Millisecond

// socketRouterClient is a proxy object for the router implemented by this API
type socketRouterClient struct {
	// Unix domain socket connected to the routing table server
	c net.Conn
	// This assumes that responses from the server are received in
	// exactly the order that messages are sent, and that the
	// responses to exactly to the process that sent the request.
	// This is tough to ensure without extra work, so all requests
	// simply take the lock, send a message, and collect their
	// response.
	l sync.Mutex
}

// New creates a new routing table connected to a routing table
// server.  It can fail if the routing table server cannot be started
// or does not respond in a timely fashion.
func New(node *kdht.NodeInfo, k int) (kdht.RoutingTable, error) {
	if len(node.Id) != kdht.KeyBytes {
		return nil, errors.New("node.Id is invalid")
	}

	// err is shared between the following goroutine and this
	// function, which is ugly
	var err error
	c := make(chan *socketRouterClient)
	go func() {
		defer close(c)
		var sr *socketRouterClient
		sr, err = connectRouter(node, k)
		if err == nil {
			c <- sr
		}
	}()

	select {
	case <-time.After(connectTimeout):
		return nil, errors.New("Router did not connect")
	case sr := <-c:
		return sr, err
	}
}

// connectRouter does the dirty work of connecting to a server.  It's
// mostly socket wrangling and error handling.  It isn't in New()
// because it was awfully long for an inline goroutine.
func connectRouter(node *kdht.NodeInfo, k int) (sr *socketRouterClient, err error) {
	// Socket addresses starting with @ are configured as abstract
	// sockets on Linux; on other operating systems it may create
	// a socket actually starting with @, I'm not clear.
	a := fmt.Sprintf("@kdhtrouter-%d-%d", os.Getpid(), rand.Int31())

	// unixpacket is SOCK_SEQPACKET on Linux, which simplifies a
	// lot of things for us.
	//
	// Unfortunately, unixpacket doesn't work on macOS, so we're
	// using unix and writing lengths.  Complain to Tim, I guess.
	l, err := net.Listen("unix", a)
	if err != nil {
		return
	}
	defer l.Close()

	// Start the server
	cmd := exec.Command(routerCommand, a)
	if err = cmd.Start(); err != nil {
		return
	}

	// The server should connect to our listening socket almost
	// immediately, so accept its incoming connection.
	var c net.Conn
	if c, err = l.Accept(); err != nil {
		return
	}

	// The init message tells the router how to configure itself.
	// This stuff could all have been provided on the command
	// line, but doing it this way has the side effect of ensuring
	// that communication is actually happening.
	sr = new(socketRouterClient)
	sr.c = c
	_, err = sr.doRequest(&kdht.RouteRequest{Type: kdht.RouteType_INIT, Node: node, I: int32(k)})
	if err != nil {
		c.Close()
		// Don't return the unusable sr
		return nil, err
	}

	return sr, err
}

// doRequest sends a message to the server and collects its response.
// Most of the failures here should never happen (they indicate a
// serialization error or a connection error, both of which are
// unlikely with protobuf and Unix sockets), but some of them could
// indicate client or server code errors, so check carefully.
func (sr *socketRouterClient) doRequest(req *kdht.RouteRequest) (resp *kdht.RouteResponse, err error) {
	// This is really the only thing that has to be serialized!
	sr.l.Lock()
	defer sr.l.Unlock()
	var m []byte
	if m, err = proto.Marshal(req); err != nil {
		return
	}
	if len(m) > int(kdht.RouteMessage_MAX_SIZE) {
		// This is an arbitrary limit, but we should not reach
		// it for reasonable sizes of k in either direction.
		return nil, errors.New("Serialized message was too large to send")
	}
	// Send the message length, hex-encoded
	if _, err = sr.c.Write([]byte(fmt.Sprintf("%04x", len(m)))); err != nil {
		return
	}
	// Send the message to the server
	if _, err = sr.c.Write(m); err != nil {
		return
	}
	// Collect the response; this should maybe time out?
	// First , get and decode the length
	lenbuf := make([]byte, 4)
	if _, err = io.ReadFull(sr.c, lenbuf); err != nil {
		return
	}
	var l int
	fmt.Sscanf(string(lenbuf), "%04x", &l)
	// Get and decode the message itself
	m = make([]byte, l)
	if _, err = io.ReadFull(sr.c, m); err != nil {
		return
	}
	var r kdht.RouteResponse
	if err = proto.Unmarshal(m, &r); err != nil {
		return
	}
	// The lock should protect us from this error on this end, but
	// the server might have messed up.
	if r.Type != req.Type {
		return nil, errors.New("Type mismatch on returned router request")
	}

	// Convert incoming errors to Go errors.  The value of
	// RouteError_NONE must be zero so that the zero value doesn't
	// cause these to trip.
	if r.Error == kdht.RouteError_INVALID {
		return &r, kdht.InvalidNodeError
	}
	if r.Error == kdht.RouteError_STRING {
		return &r, errors.New(r.Str)
	}

	return &r, nil
}

// K implements RoutingTable.K().  This is a little bit bogus because
// K() can't fail, and our router connection _could_ fail, but that
// indicates a bug in the routing table implementation, so all bets
// are off at that point anyway.
//
// It's not clear, maybe this should panic when that happens.  That
// might make sense for this assignment, but in the general case it
// seems like the wrong solution.  These functions should have allowed
// for errors.
func (sr *socketRouterClient) K() int {
	r, err := sr.doRequest(&kdht.RouteRequest{Type: kdht.RouteType_K})
	if err != nil {
		// It's really not clear what to do here, the API
		// doesn't allow for K to fail.  Just ... make
		// something up.
		return 0
	}
	return int(r.I)
}

// InsertNode satisfies RoutingTable.InsertNode().  It just passes its
// argument to the routing table and assumes it worked.
func (sr *socketRouterClient) InsertNode(node *kdht.NodeInfo) {
	// Nothing we can do if this fails, so ...
	sr.doRequest(&kdht.RouteRequest{Type: kdht.RouteType_INSERT_NODE, Node: node})
}

// RemoveNode satisfies RoutingTable.RemoveNode(), by proxing the key
// and returned error message (if any).
func (sr *socketRouterClient) RemoveNode(key []byte) error {
	_, err := sr.doRequest(&kdht.RouteRequest{Type: kdht.RouteType_REMOVE_NODE, Key: key})
	return err
}

// Lookup satisfies RoutingTable.Lookup() by proxying the key and
// returned error message.  A failure in server communication is
// indistinguishable from a lookup of an unknown node.
func (sr *socketRouterClient) Lookup(key []byte) (*kdht.NodeInfo, bool) {
	r, err := sr.doRequest(&kdht.RouteRequest{Type: kdht.RouteType_LOOKUP, Key: key})
	if err != nil {
		return nil, false
	}
	return r.Node, r.Node != nil
}

// GetNodes satisfies RoutingTable.GetNodes() by proxying the key and
// returning the nodelist.  As many other functions here, a
// communication failure is indistinguishable from an empty bucket.
func (sr *socketRouterClient) GetNodes(bucket int) []*kdht.NodeInfo {
	r, err := sr.doRequest(&kdht.RouteRequest{Type: kdht.RouteType_GET_NODES, I: int32(bucket)})
	if err != nil {
		return nil
	}
	return r.Nodes
}

// ClosestK satisfies RoutingTable.ClosestK() by proxying the key and
// returning the nodelist, same as GetNodes().
func (sr *socketRouterClient) ClosestK(key []byte) []*kdht.NodeInfo {
	r, err := sr.doRequest(&kdht.RouteRequest{Type: kdht.RouteType_CLOSEST_K, Key: key})
	if err != nil {
		return nil
	}
	return r.Nodes
}

// Buckets satisfies RoutingTable.Buckets() through the proxy, and
// doesn't know what to do with errors either.
func (sr *socketRouterClient) Buckets() int {
	// Similar problem to K(), in that the protocol doesn't
	// account for failure
	r, err := sr.doRequest(&kdht.RouteRequest{Type: kdht.RouteType_BUCKETS})
	if err != nil {
		return 0
	}
	return int(r.I)
}
