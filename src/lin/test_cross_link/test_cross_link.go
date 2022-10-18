package main

import (
	"fmt"
	"lin/lin_common"
)


type cross_link_inf struct {

}

func (pthis*cross_link_inf)Ntf_node_in_view(node_id int, node_id_in_viewed int) {

}
func (pthis*cross_link_inf)Ntf_node_out_view(node_id int, node_id_out_view int) {

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

	fmt.Println(clm.Crosslink_mgr_dump())

	node1 := &cross_link_node_inf {
		x : 1,
		y : 1,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node1)

	node2 := &cross_link_node_inf {
		x : 11,
		y : 11,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node2)

	node3 := &cross_link_node_inf {
		x : 16,
		y : 16,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node3)

	node4 := &cross_link_node_inf {
		x : 21,
		y : 21,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node4)

	node5 := &cross_link_node_inf {
		x : -6,
		y : -6,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node5)

	node6 := &cross_link_node_inf {
		x : -11,
		y : -11,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node6)

	node7 := &cross_link_node_inf {
		x : -16,
		y : -16,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node7)

	node8 := &cross_link_node_inf {
		x : -21,
		y : -21,
		view_range : 10,
	}
	clm.Crosslink_mgr_add(node8)


	clm.Crosslink_mgr_update_pos(node0_id, 12, 12)

	return
}
