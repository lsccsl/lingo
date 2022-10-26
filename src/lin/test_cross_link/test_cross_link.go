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


type cross_link_node_inf struct {
	x float32
	y float32
	data interface{}
	view_range float32
}
func (pthis*cross_link_node_inf)Get_node_x()float32 {
	return pthis.x
}
func (pthis*cross_link_node_inf)Get_node_y()float32 {
	return pthis.y
}
func (pthis*cross_link_node_inf)Get_node_data()interface{} {
	return pthis.data
}
func (pthis*cross_link_node_inf) Get_view_range()float32 {
	return pthis.view_range
}


func main() {

	var cli cross_link_inf
	clm := lin_common.Crosslink_mgr_constructor(&cli)

	// test add
	fmt.Println("add 100 step 2")
	for i := 2; i <= 100; i += 2 {
		node_tmp := &cross_link_node_inf {
			x : float32(i),
			y : float32(i),
			view_range : 10,
		}
		clm.Crosslink_mgr_add(node_tmp)
		//fmt.Println(clm)
		clm.Check()
		fmt.Print(".")
	}
	for i := 1; i <= 99; i += 2 {
		node_tmp := &cross_link_node_inf {
			x : float32(i),
			y : float32(i),
			view_range : 10,
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

		node_tmp := &cross_link_node_inf {
			x : x,
			y : y,
			view_range : 10,
		}
		clm.Crosslink_mgr_add(node_tmp)
		//fmt.Println(clm)
		clm.Check()
		fmt.Print(".")
		if i % 100 == 0 {
			fmt.Println()
		}
	}

	node0 := &cross_link_node_inf {
		x : 0,
		y : 0,
		view_range : 10,
	}
	node0_id := clm.Crosslink_mgr_add(node0)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node1 := &cross_link_node_inf {
		x : 1,
		y : 1,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node1)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node2 := &cross_link_node_inf {
		x : 11,
		y : 11,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node2)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node3 := &cross_link_node_inf {
		x : 16,
		y : 16,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node3)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node4 := &cross_link_node_inf {
		x : 21,
		y : 21,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node4)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node5 := &cross_link_node_inf {
		x : -6,
		y : -6,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node5)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node6 := &cross_link_node_inf {
		x : -11,
		y : -11,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node6)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node7 := &cross_link_node_inf {
		x : -16,
		y : -16,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node7)
	//fmt.Println(clm.Crosslink_mgr_dump())
	clm.Check()

	node8 := &cross_link_node_inf {
		x : -21,
		y : -21,
		view_range : 10,
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

		node_tmp := &cross_link_node_inf {
			x : x,
			y : y,
			view_range : 10,
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
