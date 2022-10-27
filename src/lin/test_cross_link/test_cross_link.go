package main

import (
	"fmt"
	"lin/lin_common"
	"math/rand"
)


type cross_link_inf struct {

}

func (pthis*cross_link_inf)Ntf_node_in_view(node_id int, node_id_in_view int) {
	//fmt.Println(node_id, " view ", node_id_in_view)
}
func (pthis*cross_link_inf)Ntf_node_out_view(node_id int, node_id_out_view int) {
	//fmt.Println(node_id, " out view ", node_id_out_view)
}

func main() {

	var cli cross_link_inf
	clm := lin_common.Crosslink_mgr_constructor(&cli)
	// test add
	fmt.Println("add 100 step 2")
	for i := 2; i <= 100; i += 2 {
		node_tmp := &lin_common.Crosslink_node_param {
			X : float32(i),
			Y : float32(i),
			ViewRange : 10,
		}
		clm.Crosslink_mgr_add(node_tmp)
		//fmt.Println(clm)
		clm.Check()
		fmt.Print(".")
	}
	for i := 1; i <= 99; i += 2 {
		node_tmp := &lin_common.Crosslink_node_param {
			X : float32(i),
			Y : float32(i),
			ViewRange : 10,
		}
		clm.Crosslink_mgr_add(node_tmp)
		//fmt.Println(clm)
		clm.Check()
		fmt.Print(".")
	}

	// test rand add
	fmt.Println("\nadd rand 1000")
	rand.Seed(0)
	for i := 1; i <= 1000; i ++ {
		x := float32(rand.Int() % 100)
		y := float32(rand.Int() % 100)

		node_tmp := &lin_common.Crosslink_node_param {
			X : x,
			Y : y,
			ViewRange : 10,
		}
		clm.Crosslink_mgr_add(node_tmp)
		//fmt.Println(clm)
		clm.Check()
		fmt.Print(".")
		if i % 100 == 0 {
			fmt.Println()
		}
	}

	node0 := &lin_common.Crosslink_node_param {
		X : 0,
		Y : 0,
		ViewRange : 10,
	}
	node0_id := clm.Crosslink_mgr_add(node0)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node1 := &lin_common.Crosslink_node_param {
		X : 1,
		Y : 1,
		ViewRange : 10,
	}
	clm.Crosslink_mgr_add(node1)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node2 := &lin_common.Crosslink_node_param {
		X : 11,
		Y : 11,
		ViewRange : 10,
	}
	clm.Crosslink_mgr_add(node2)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node3 := &lin_common.Crosslink_node_param {
		X : 16,
		Y : 16,
		ViewRange : 10,
	}
	clm.Crosslink_mgr_add(node3)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node4 := &lin_common.Crosslink_node_param {
		X : 21,
		Y : 21,
		ViewRange : 10,
	}
	clm.Crosslink_mgr_add(node4)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node5 := &lin_common.Crosslink_node_param {
		X : -6,
		Y : -6,
		ViewRange : 10,
	}
	clm.Crosslink_mgr_add(node5)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node6 := &lin_common.Crosslink_node_param {
		X : -11,
		Y : -11,
		ViewRange : 10,
	}
	clm.Crosslink_mgr_add(node6)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node7 := &lin_common.Crosslink_node_param {
		X : -16,
		Y : -16,
		ViewRange : 10,
	}
	clm.Crosslink_mgr_add(node7)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node8 := &lin_common.Crosslink_node_param {
		X : -21,
		Y : -21,
		ViewRange : 10,
	}
	clm.Crosslink_mgr_add(node8)
	//fmt.Println(clm)
	clm.Check()

	fmt.Println("add", clm.Crosslink_get_node_count())


	// update pos
	fmt.Println("\n test update pos 0-100")
	for i := 0; i < 100; i ++ {
		clm.Crosslink_mgr_update_pos(node0_id, float32(i), float32(i))
		//fmt.Println(clm)
		clm.Check()
		if i % 100 == 0 {
			fmt.Println()
		}
		fmt.Print(".")
	}
	//fmt.Println(clm)
	fmt.Println("\n test update pos 100-0")
	for i := 0; i < 100; i ++ {
		clm.Crosslink_mgr_update_pos(node0_id, float32(100 - i), float32(100 - i))
		//fmt.Println(clm)
		clm.Check()
		if i % 100 == 0 {
			fmt.Println()
		}
		fmt.Print(".")
	}

	// update pos rand 100
	fmt.Println("\n test update pos rand 100")
	rand.Seed(0)
	for i := 0; i < 1000; i ++ {
		x := float32(rand.Int() % 100)
		y := float32(rand.Int() % 100)
		clm.Crosslink_mgr_update_pos(node0_id, x, y)
		//fmt.Println(clm)
		clm.Check()
		if i % 100 == 0 {
			fmt.Println()
		}
		fmt.Print(".")
	}

	// udpate pos rand 1000
	fmt.Println("\n test update pos rand 1000")
	rand.Seed(0)
	for i := 0; i < 1000; i ++ {
		x := float32(rand.Int() % 1000)
		y := float32(rand.Int() % 1000)
		clm.Crosslink_mgr_update_pos(node0_id, x, y)
		//fmt.Println(clm)
		clm.Check()
		if i % 100 == 0 {
			fmt.Println()
		}
		fmt.Print(".")
	}

	// add rand 1000
	fmt.Println("\n test add rand 1000")
	map_del := make(map[int]int)
	rand.Seed(0)
	for i := 1; i <= 100; i ++ {
		x := float32(rand.Int() % 1000)
		y := float32(rand.Int() % 1000)

		node_tmp := &lin_common.Crosslink_node_param {
			X : x,
			Y : y,
			ViewRange : 10,
		}
		id := clm.Crosslink_mgr_add(node_tmp)
		//fmt.Println(clm)
		clm.Check()
		map_del[id] = id
		if i % 100 == 0 {
			fmt.Println()
		}
		fmt.Print(".")
	}
	fmt.Println("add", clm.Crosslink_get_node_count())

	// udpate pos rand 1000
	fmt.Println("\n test update pos rand 1000")
	rand.Seed(0)
	for i := 0; i < 1000; i ++ {
		x := float32(rand.Int() % 1000)
		y := float32(rand.Int() % 1000)
		clm.Crosslink_mgr_update_pos(node0_id, x, y)
		//fmt.Println(clm)
		clm.Check()
		if i % 100 == 0 {
			fmt.Println()
		}
		fmt.Print(".")
	}

	// del
	fmt.Println("\n test del")
	for k, _ := range map_del {
		clm.Crosslink_mgr_del(k)
		clm.Check()
		fmt.Print(".")
	}

	fmt.Println("delete", clm.Crosslink_get_node_count())

	fmt.Println("\n\n\n delete\n")
	clm.Crosslink_mgr_del(node0_id)
	//fmt.Println(clm)
	clm.Check()

	return
}
