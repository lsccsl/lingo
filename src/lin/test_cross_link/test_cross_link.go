package main

import (
	"fmt"
	"lin/lin_common"
)


type cross_link_inf struct {

}

func (pthis*cross_link_inf)Ntf_node_in_view(node_id int, node_id_in_view int) {
	fmt.Println(node_id, " view ", node_id_in_view)
}
func (pthis*cross_link_inf)Ntf_node_out_view(node_id int, node_id_out_view int) {
	fmt.Println(node_id, " out view ", node_id_out_view)
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

	node0 := &cross_link_node_inf {
		x : 0,
		y : 0,
		view_range : 10,
	}
	node0_id := clm.Crosslink_mgr_add(node0)
	//fmt.Println(clm.Crosslink_mgr_dump())

	node1 := &cross_link_node_inf {
		x : 1,
		y : 1,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node1)
	//fmt.Println(clm.Crosslink_mgr_dump())

	node2 := &cross_link_node_inf {
		x : 11,
		y : 11,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node2)
	//fmt.Println(clm.Crosslink_mgr_dump())

	node3 := &cross_link_node_inf {
		x : 16,
		y : 16,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node3)
	//fmt.Println(clm.Crosslink_mgr_dump())

	node4 := &cross_link_node_inf {
		x : 21,
		y : 21,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node4)
	//fmt.Println(clm.Crosslink_mgr_dump())

	node5 := &cross_link_node_inf {
		x : -6,
		y : -6,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node5)
	//fmt.Println(clm.Crosslink_mgr_dump())

	node6 := &cross_link_node_inf {
		x : -11,
		y : -11,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node6)
	//fmt.Println(clm.Crosslink_mgr_dump())

	node7 := &cross_link_node_inf {
		x : -16,
		y : -16,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node7)
	//fmt.Println(clm.Crosslink_mgr_dump())

	node8 := &cross_link_node_inf {
		x : -21,
		y : -21,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node8)
	//fmt.Println(clm)
	clm.Check()

	fmt.Println("\n\n\n update pos\n")
	clm.Crosslink_mgr_update_pos(node0_id, 11, 11)
	//fmt.Println(clm)
	clm.Check()

	fmt.Println("\n\n\n delete\n")
	clm.Crosslink_mgr_del(node0_id)
	fmt.Println(clm)
	clm.Check()

	return
}
