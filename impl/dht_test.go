package impl

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"cse586.kdht/api/kdht"
	"cse586.kdht/given/keys"
)

const (
	Address1 = "localhost:4586"
	Address2 = "localhost:5486"
	Address3 = "localhost:1986"
	Address4 = "localhost:5905"
	Address5 = "localhost:1943"
)

func TestDHT_Standard(t *testing.T) {
	niters := 1
	kmin := 2
	kmax := 2
	alphamin := 2
	alphamax := 2
	bufferTime := 3 * time.Second

	for iter := 0; iter < niters; iter++ {
		for k := kmin; k <= kmax; k++ {
			for alpha := alphamin; alpha <= alphamax; alpha++ {
				func(k int, alpha int) {
					fmt.Printf("Testing k = %v, alpha = %v\n", k, alpha)

					node1, err1 := NewNode(byteToKey(0x10), Address1, k, alpha, []string{Address2})
					node2, err2 := NewNode(byteToKey(0x20), Address2, k, alpha, []string{Address3})
					node3, err3 := NewNode(byteToKey(0x30), Address3, k, alpha, []string{Address4})
					node4, err4 := NewNode(byteToKey(0x40), Address4, k, alpha, []string{Address5})
					node5, err5 := NewNode(byteToKey(0x50), Address5, k, alpha, []string{Address1})
					nodes := []*KdmNode{node1, node2, node3, node4, node5}
					errs := []error{err1, err2, err3, err4, err5}

					defer func() {
						for i, node := range nodes {
							if node == nil {
								continue
							}
							err := node.Shutdown()
							if err != nil {
								t.Logf("(node %v shutdown failed) %v", i, err)
							}
						}
					}()

					for i, err := range errs {
						if err != nil {
							t.Fatalf("(node%v creation failed) %v", i, err)
						}
					}

					exp1 := []*kdht.NodeInfo{node1.info, node3.info, node2.info, node5.info, node4.info}
					exp2 := []*kdht.NodeInfo{node2.info, node3.info, node1.info, node4.info, node5.info}
					exp3 := []*kdht.NodeInfo{node3.info, node2.info, node1.info, node5.info, node4.info}
					exp4 := []*kdht.NodeInfo{node4.info, node5.info, node1.info, node2.info, node3.info}
					exp5 := []*kdht.NodeInfo{node5.info, node4.info, node1.info, node3.info, node2.info}
					exps := [][]*kdht.NodeInfo{exp1, exp2, exp3, exp4, exp5}

					time.Sleep(bufferTime)

					for i := 0; i < len(nodes); i++ {
						for j := 0; j < len(nodes); j++ {
							t.Logf("(node%v) find node%v (k = %v, alpha = %v)\n", i+1, j+1, k, alpha)

							exp := exps[j]
							if k <= len(exp) {
								exp = exp[0:k]
							}

							act, err := nodes[i].FindNode(nodes[j].info.Id)
							if err != nil {
								t.Fatalf("(node%v FindNode failed) %v", i, err)
							}

							if len(act) != len(exp) {
								msg := "exp and act differ in length:\n"
								msg += sprintInfos("exp", exp) + "\n"
								msg += sprintInfos("act", act)
								t.Fatalf("%v", msg)
							}

							for k := 0; k < len(exp); k++ {
								flag := false
								for l := 0; l < len(act); l++ {
									if bytes.Equal(exp[k].Id, act[l].Id) {
										flag = true
									}
								}
								if !flag {
									msg := "exp and act differ:\n"
									msg += sprintInfos("exp", exp) + "\n"
									msg += sprintInfos("act", act)
									t.Fatalf("%v\n\n", msg)
								}
							}

							t.Logf("passed\n\n")
						}
					}
				}(k, alpha)
			}
		}
	}
}

func TestDHT_Store(t *testing.T) {
	niters := 1
	k := 2
	alpha := 1
	bufferTime := 3 * time.Second

	for iter := 0; iter < niters; iter++ {
		func(k int, alpha int) {
			node1, err1 := NewNode(byteToKey(0x10), Address1, k, alpha, []string{Address2})
			node2, err2 := NewNode(byteToKey(0x20), Address2, k, alpha, []string{Address3})
			node3, err3 := NewNode(byteToKey(0x30), Address3, k, alpha, []string{Address4})
			node4, err4 := NewNode(byteToKey(0x40), Address4, k, alpha, []string{Address5})
			node5, err5 := NewNode(byteToKey(0x50), Address5, k, alpha, []string{Address1})
			nodes := []*KdmNode{node1, node2, node3, node4, node5}
			errs := []error{err1, err2, err3, err4, err5}

			defer func() {
				for i, node := range nodes {
					if node == nil {
						continue
					}
					err := node.Shutdown()
					if err != nil {
						t.Logf("(node %v shutdown failed) %v", i, err)
					}
				}
			}()

			for i, err := range errs {
				if err != nil {
					t.Fatalf("(node%v creation failed) %v", i, err)
				}
			}

			val1 := []byte("val1")
			val2 := []byte("val2")
			val3 := []byte("val3")
			val4 := []byte("val4")
			val5 := []byte("val5")
			vals := [][]byte{val1, val2, val3, val4, val5}

			time.Sleep(bufferTime)

			for i, exp := range vals {
				key := keys.Compute(exp)
				fmt.Printf("Testing:\n    act = %v\n    key = %v\n\n", exp, key)

				err := node1.Store(exp)
				if err != nil {
					t.Fatalf("(node%v Store failed) %v", i, err)
				}

				time.Sleep(bufferTime)

				flag := false
				for j, node := range nodes {
					act, ok := node.localStorage[string(key)]
					if ok {
						if bytes.Equal(exp, act) {
							t.Logf("node%v is storing val%v", j+1, i+1)
							flag = true
						} else {
							msg := fmt.Sprintf("node%v stored an incorrect value for val%v:\n", j+1, i+1)
							msg += fmt.Sprintf("    exp: %v\n", exp)
							msg += fmt.Sprintf("    act: %v\n", act)
							t.Fatalf("%v", msg)
						}
					}
				}

				if !flag {
					t.Fatalf("val%v wasn't stored\n", i+1)
				}

				flag = false
				for j, node := range nodes {
					act, info, err := node.FindValue(key)
					if errors.Is(err, kdht.ValueError) {
						continue
					}

					if err != nil {
						t.Fatalf("(node%v FindValue failed) %v", j+1, err)
					}

					if bytes.Equal(exp, act) {
						t.Logf("node%v found val%v from %v", j+1, i+1, info.Id)
						flag = true
					} else {
						msg := fmt.Sprintf("node%v stored an incorrect value for val%v:\n", j+1, i+1)
						msg += fmt.Sprintf("    exp: %v\n", exp)
						msg += fmt.Sprintf("    act: %v\n", act)
						t.Fatalf("%v", msg)
					}
				}

				if !flag {
					t.Fatalf("val%v wasn't found\n", i+1)
				}

				t.Logf("passed\n\n")
			}
		}(k, alpha)
	}
}

func TestDHT_Neighbors(t *testing.T) {
	niters := 1
	k := 2
	alpha := 1
	bufferTime := 3 * time.Second

	for iter := 0; iter < niters; iter++ {
		func(k int, alpha int) {
			node1, err1 := NewNode(byteToKey(0x10), Address1, k, alpha, []string{Address2})
			node2, err2 := NewNode(byteToKey(0x20), Address2, k, alpha, []string{Address3})
			node3, err3 := NewNode(byteToKey(0x30), Address3, k, alpha, []string{Address4})
			node4, err4 := NewNode(byteToKey(0x40), Address4, k, alpha, []string{Address5})
			node5, err5 := NewNode(byteToKey(0x50), Address5, k, alpha, []string{Address1})
			nodes := []*KdmNode{node1, node2, node3, node4, node5}
			errs := []error{err1, err2, err3, err4, err5}

			defer func() {
				for i, node := range nodes {
					if node == nil {
						continue
					}
					err := node.Shutdown()
					if err != nil {
						t.Logf("(node %v shutdown failed) %v", i, err)
					}
				}
			}()

			for i, err := range errs {
				if err != nil {
					t.Fatalf("(node%v creation failed) %v", i, err)
				}
			}

			exp1 := []*kdht.NodeInfo{node2.info, node5.info}
			exp2 := []*kdht.NodeInfo{node1.info, node3.info}
			exp3 := []*kdht.NodeInfo{node2.info, node4.info}
			exp4 := []*kdht.NodeInfo{node3.info, node5.info}
			exp5 := []*kdht.NodeInfo{node4.info, node1.info}
			exps := [][]*kdht.NodeInfo{exp1, exp2, exp3, exp4, exp5}

			time.Sleep(bufferTime)

			for i, node := range nodes {
				fmt.Printf("Testing: node%v\n", i+1)

				exp := exps[i]
				act := node.Neighbors()

				if len(act) != len(exp) {
					msg := "exp and act differ in length:\n"
					msg += sprintInfos("exp", exp) + "\n"
					msg += sprintInfos("act", act)
					t.Fatalf("%v", msg)
				}

				for k := 0; k < len(exp); k++ {
					flag := false
					for l := 0; l < len(act); l++ {
						if bytes.Equal(exp[k].Id, act[l].Id) {
							flag = true
							break
						}
					}
					if !flag {
						msg := "exp and act differ:\n"
						msg += sprintInfos("exp", exp) + "\n"
						msg += sprintInfos("act", act)
						t.Fatalf("%v\n\n", msg)
					}
				}

				t.Logf("passed\n\n")
			}
		}(k, alpha)
	}
}

func TestDHT_Complex(t *testing.T) {
	k := 2
	alpha := 1
	bufferTime := 3 * time.Second

	fmt.Printf("Testing k = %v, alpha = %v\n", k, alpha)

	node1, err1 := NewNode(byteToKey(0x10), Address1, k, alpha, []string{Address2})
	node2, err2 := NewNode(byteToKey(0x20), Address2, k, alpha, []string{Address1})
	node3, err3 := NewNode(byteToKey(0x30), Address3, k, alpha, []string{Address2, Address4})
	node4, err4 := NewNode(byteToKey(0x40), Address4, k, alpha, []string{Address5})
	node5, err5 := NewNode(byteToKey(0x50), Address5, k, alpha, []string{Address4})
	nodes := []*KdmNode{node1, node2, node3, node4, node5}
	errs := []error{err1, err2, err3, err4, err5}

	defer func() {
		for i, node := range nodes {
			if node == nil {
				continue
			}
			err := node.Shutdown()
			if err != nil {
				t.Logf("(node %v shutdown failed) %v", i, err)
			}
		}
	}()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("(node%v creation failed) %v", i, err)
		}
	}

	time.Sleep(bufferTime)

	closest, err := node1.FindNode(node5.info.Id)
	if err != nil {
		t.Fatalf("(node1 FindNode failed) %v", err)
	}

	flag := false
	for _, info := range closest {
		if bytes.Equal(info.Id, node5.info.Id) {
			flag = true
		}
	}

	if !flag {
		t.Fatalf("node5 wasn't found\n")
	}

	t.Logf("node1 found node5\n")
	t.Logf("passed\n\n")
}

func sprintInfos(msg string, infos []*kdht.NodeInfo) string {
	str := fmt.Sprintf("%v:\n", msg)
	lf := ""
	for _, info := range infos {
		str += fmt.Sprintf("%v    %v", lf, info.Id)
		lf = "\n"
	}
	return str
}

func byteToKey(byt byte) []byte {
	key := make([]byte, kdht.KeyBytes)
	key[0] = byt
	return key
}
