package impl

/*
import (
	"fmt"
	"slices"
	"testing"

	"cse586.kdht/api/kdht"
	"cse586.kdht/impl"
)

func TestRouting_Standard(t *testing.T) {
	me := &kdht.NodeInfo{Address: "me", Id: byteToKey(0x10)}
	a := &kdht.NodeInfo{Address: "a", Id: byteToKey(0x90)}
	b := &kdht.NodeInfo{Address: "b", Id: byteToKey(0x20)}
	c := &kdht.NodeInfo{Address: "c", Id: byteToKey(0xD0)}
	d := &kdht.NodeInfo{Address: "d", Id: byteToKey(0xB0)}
	e := &kdht.NodeInfo{Address: "e", Id: byteToKey(0x30)}
	f := &kdht.NodeInfo{Address: "f", Id: byteToKey(0x40)}
	g := &kdht.NodeInfo{Address: "g", Id: byteToKey(0x70)}

	table, err := impl.NewRoutingTable(me, 2)
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}

	testBucketContents(t, table, 159, me)
	testBucketContents(t, table, 158)

	table.InsertNode(a)
	testBucketContents(t, table, 159, me, a)
	testBucketContents(t, table, 158)

	table.InsertNode(b)
	testBucketContents(t, table, 159, a)
	testBucketContents(t, table, 158, me, b)
	testBucketContents(t, table, 157)

	table.InsertNode(c)
	testBucketContents(t, table, 159, a, c)
	testBucketContents(t, table, 158, me, b)
	testBucketContents(t, table, 157)

	table.InsertNode(d)
	testBucketContents(t, table, 159, a, c)
	testBucketContents(t, table, 158, me, b)
	testBucketContents(t, table, 157)

	table.InsertNode(e)
	testBucketContents(t, table, 159, a, c)
	testBucketContents(t, table, 158)
	testBucketContents(t, table, 157, e, b)
	testBucketContents(t, table, 156, me)
	testBucketContents(t, table, 155)

	table.InsertNode(f)
	testBucketContents(t, table, 159, a, c)
	testBucketContents(t, table, 158, f)
	testBucketContents(t, table, 157, e, b)
	testBucketContents(t, table, 156, me)
	testBucketContents(t, table, 155)

	logClosestK(t, table, me)
	logClosestK(t, table, a)
	logClosestK(t, table, b)
	logClosestK(t, table, c)
	logClosestK(t, table, d)
	logClosestK(t, table, e)
	logClosestK(t, table, f)
	logClosestK(t, table, g)

	testLookup(t, table, me, true)
	testLookup(t, table, a, true)
	testLookup(t, table, b, true)
	testLookup(t, table, c, true)
	testLookup(t, table, d, false)
	testLookup(t, table, e, true)
	testLookup(t, table, f, true)
	testLookup(t, table, g, false)

	testRemove(t, table, a, true)
	testBucketContents(t, table, 159, c)
	testRemove(t, table, b, true)
	testBucketContents(t, table, 157, e)
	testRemove(t, table, c, true)
	testBucketContents(t, table, 159)
	testRemove(t, table, d, false)
	testRemove(t, table, e, true)
	testBucketContents(t, table, 157)
	testRemove(t, table, f, true)
	testBucketContents(t, table, 157)
	testRemove(t, table, g, false)
	testRemove(t, table, me, false)
}

func testBucketContents(
	t *testing.T,
	table kdht.RoutingTable, num int,
	nodes ...*kdht.NodeInfo) {

	bucket := table.GetNodes(num)

	if len(nodes) == 0 && bucket != nil {
		t.Fatalf("FAILURE: bucket %v should have been empty", num)
	}

	if len(bucket) != len(nodes) {
		t.Fatalf("FAILURE: bucket %v should have had %v node(s), but had %v instead",
			num, len(nodes), len(bucket))
	}

	for _, node := range nodes {
		if !slices.Contains(bucket, node) {
			t.Fatalf("FAILURE: bucket %v did not contain node %v",
				num, node.GetAddress())
		}
	}
}

func testRemove(
	t *testing.T,
	table kdht.RoutingTable,
	node *kdht.NodeInfo, exp bool) {

	err := table.RemoveNode(node.GetId())
	if exp && err != nil {
		t.Fatalf("ERROR: %s", err)
	}

	if !exp {
		return
	}

	testLookup(t, table, node, false)
}

func testLookup(
	t *testing.T,
	table kdht.RoutingTable,
	node *kdht.NodeInfo, exp bool) {

	fnode, ok := table.Lookup(node.GetId())
	if exp != ok {
		if exp {
			t.Fatalf("FAILURE: node %v was expected to be found but it wasn't",
				node.GetAddress())
		} else {
			t.Fatalf("FAILURE: node %v was not expected to be found but it was",
				node.GetAddress())
		}
	}

	if !exp {
		return
	}

	if fnode == nil {
		t.Fatalf("FAILURE: node %v was expected to be found but Lookup returned nil",
			node.GetAddress())
	}

	if node != fnode {
		t.Fatalf("FAILURE: node %v was expected to be found but node %v was found",
			node.GetAddress(), fnode.GetAddress())
	}
}

func logClosestK(
	t *testing.T,
	table kdht.RoutingTable,
	node *kdht.NodeInfo) {

	closest := table.ClosestK(node.GetId())
	sep := ""

	str := fmt.Sprintf("Closest k for %v: (", node.GetAddress())
	for _, cnode := range closest {
		str += fmt.Sprintf("%v%v", sep, cnode.GetAddress())
		sep = ","
	}
	str += ")"

	t.Logf("%v; len = %v", str, len(closest))
}
*/
