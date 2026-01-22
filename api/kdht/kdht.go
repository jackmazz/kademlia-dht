/*
Copyright 2021, 2023 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more information.
*/

// The api/kdht package contains constants, interfaces, and messages
// related to implementing our k-DHT protocol.
package kdht

// The following error values are defined here so that they can be
// used with errors.Is().  If you want to provide more detailed
// information, look at fmt.Errorf() and %w.

// ShutdownError indicates that an operation was attempted on a
// node which has completed or is executing Shutdown().
var ShutdownError = &kError{"This node is already shut down"}

// StorageError is returned by Node.Store() when a value cannot be
// stored to K nodes, either due to failures or simply too few nodes
// in the DHT.
var StorageError = &kError{"The value could not be stored in K nodes"}

// ValueError indicates that a value could not be retrieved by
// Node.FindValue() because no reachable node stores the value.
var ValueError = &kError{"The value could not be found"}

var InvalidNodeError = &kError{"The node does not exist"}

// Node represents an instance of a k-DHT node, and the operations that
// can be performed on that node.  If the node has been shut down, any
// operation invoked on the node should return a ShutdownError.
type Node interface {
	// Ping sends a ping message containing the given message to
	// the given node id.
	//
	// If this node has been shut down, the node cannot be
	// contacted, or the node does not respond to the ping, an
	// error is returned.
	Ping(id []byte, message []byte) error

	// Store stores the given value into the DHT at its address
	// (computed as the SHA-1 sum of the value itself) by sending
	// a store message to the K nodes closest to the address.
	//
	// If fewer than K nodes can be found to store the value,
	// StorageError is returned (even if it was successfully
	// stored on some nodes).
	Store(value []byte) error

	// FindNode looks up the K nodes closest to the given node
	// ID, per the algorithm specified in the project handout.
	// Nodes that do not respond are not included in the returned
	// set.  If fewer than K active nodes cannot be found, fewer
	// than K nodes are returned.  Note that this method cannot
	// fail if this node is active, as it can return a slice
	// containing only  itself.
	FindNode(id []byte) ([]*NodeInfo, error)

	// FindValue retrieves the value for a given ID from the DHT.
	// As with Kademlia, the process for FindValue is identical to
	// FindNode, except that a node that stores the value will
	// send it instead of a list of K nodes.  If the value is
	// found, both the value and the node at which it was found
	// are returned to the caller.
	//
	// If the value cannot be found, ValueError is returned
	// instead.
	FindValue(id []byte) ([]byte, NodeInfo, error)

	// Shutdown stops this k-DHT node.  Its listening socket is
	// closed, and any ongoing operations stop as soon as is
	// practical, returning ShutdownError.
	//
	// This method returns ShutdownError if the node is already
	// shut down.
	Shutdown() error

	// Neighbors returns a slice of all of the neighbors known to
	// this node.
	//
	// The returned slice may be empty if this node does not know
	// of any neighbors.
	Neighbors() []*NodeInfo
}

// RoutingTable is an interface to a routing table at a particular
// k-DHT Node.  You must implement RoutingTable as a separate
// interface so that it can be tested separately from your Node.  Your
// Node should use an instance of your routing table in its operation.
//
// The routing table must be safe for concurrent access.  You may wish
// to utilize data structures from the sync package, such as sync.Map
// or sync.Mutex, to accomplish this.  You are NOT required to
// maintain any particular semantics for concurrent access (e.g., if
// InsertNode is called while a GetNodes operation is ongoing), but
// your implementation must not corrupt its data structures or crash.
// You may use whatever semantics for concurrent access make sense for
// your project.  (It seems likely that simple mutual exclusion is the
// route to take, however.)
//
// Remember that your implementation may provide other features and
// methods than those listed here if they are helpful to you.
type RoutingTable interface {
	// K returns the value of K used by this routing table.
	K() int

	// InsertNode adds the node having the specified information
	// to the routing table.  This operation cannot fail, although
	// the node may not be inserted if the corresponding k-bucket
	// is already full.
	InsertNode(node *NodeInfo)

	// RemoveNode removes the node having the specified
	// identifier from the routing table.
	//
	// This operation fails and returns InvalidNodeError if the
	// node is not present in the table or local node's ID is
	// attempted to be removed.
	RemoveNode(key []byte) error

	// Lookup finds a particular node in the local routing table
	// by its ID, if it exists.  If it does not exist, nothing
	// happens and ok is false.
	Lookup(key []byte) (node *NodeInfo, ok bool)

	// Retrieve the nodes stored in a numbered bucket.  Bucket 0
	// is the bucket representing nodes that differ only in the
	// least-significant bit, and bucket 159 is the bucket
	// representing nodes that differ in the most-significant bit.
	// This is the same notation used in the Kademlia paper, and
	// the same behavior provided by given/keys.DistanceBucket().
	//
	// This operation returns nil if the specified bucket is empty.
	GetNodes(bucket int) []*NodeInfo

	// Return the K closest nodes to the specified key.  They may
	// be in a single bucket, or may be spread across several
	// buckets.
	//
	// This operation cannot fail and cannot return an empty slice
	// because the local node is always present.
	ClosestK(key []byte) []*NodeInfo

	// Buckets returns the number of non-empty buckets in this
	// routing table.  These buckets necessarily start with bucket
	// number 159 and work their way down to 160 - Buckets().  The
	// return value of this function is in 0 <= Buckets() < 160.
	Buckets() int
}

// kError is an internal type that represents an error in a KDHT
// operation.
type kError struct {
	msg string
}

// Error satisfies the error type requirement.
func (e *kError) Error() string {
	return string(e.msg)
}
