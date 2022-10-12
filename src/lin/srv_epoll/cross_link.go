package main

import "lin/lin_common"

type CROSS_LINK_NODE_TYPE int
const (
	CROSS_LINK_NODE_TYPE_node        CROSS_LINK_NODE_TYPE = 1
	CROSS_LINK_NODE_TYPE_front_guard CROSS_LINK_NODE_TYPE = 2
	CROSS_LINK_NODE_TYPE_back_guard  CROSS_LINK_NODE_TYPE = 3
)

type cross_link_interface interface {
	ntf_node_in_view(node_id int, node_id_in_viewed int)
	ntf_node_out_view(node_id int, node_id_out_view int)
}

type cross_link_node_if interface {
	get_node_x()float32
	get_node_y()float32
	get_node_data()interface{}
	get_view_range()float32
}

type cross_linker interface {
	get_prev()cross_linker
	set_prev(prev cross_linker)

	get_next()cross_linker
	set_next(next cross_linker)

	get_coord()float32 // x or y
	get_view_range()float32

	get_node_type()CROSS_LINK_NODE_TYPE

	get_node_id()int
}

type MAP_NODE_ID map[int]int
type cross_linker_node struct {
	prev_ cross_linker
	next_ cross_linker

	front_ * cross_linker_front_guard
	back_  * cross_linker_back_guard

	coord_ float32
	view_range_ float32

	node_type_ CROSS_LINK_NODE_TYPE

	node_id_ int
}
type cross_node struct {
	x_node_ *cross_linker_node
	y_node_ *cross_linker_node

	node_id_ int
	node_if_ cross_link_node_if

	map_view_by_ MAP_NODE_ID // view by
	map_view_ MAP_NODE_ID // view
}
func (pthis*cross_linker_node)get_prev()cross_linker {
	return pthis.prev_
}
func (pthis*cross_linker_node)set_prev(prev cross_linker) {
	pthis.prev_ = prev
}

func (pthis*cross_linker_node)get_next()cross_linker {
	return pthis.next_
}
func (pthis*cross_linker_node)set_next(next cross_linker) {
	pthis.next_ = next
}

func (pthis*cross_linker_node)get_coord()float32 {
	return pthis.coord_
}
func (pthis*cross_linker_node)get_view_range()float32 {
	return pthis.view_range_
}

func (pthis*cross_linker_node)get_node_type()CROSS_LINK_NODE_TYPE {
	return pthis.node_type_
}

func (pthis*cross_linker_node)get_node_id()int {
	return pthis.node_id_
}


type cross_linker_front_guard struct {
	prev_ cross_linker
	next_ cross_linker

	coord_ float32

	node_type_ CROSS_LINK_NODE_TYPE
}
func (pthis*cross_linker_front_guard)get_prev()cross_linker {
	return pthis.prev_
}
func (pthis*cross_linker_front_guard)set_prev(prev cross_linker) {
	pthis.prev_ = prev
}

func (pthis*cross_linker_front_guard)get_next()cross_linker {
	return pthis.next_
}
func (pthis*cross_linker_front_guard)set_next(next cross_linker) {
	pthis.next_ = next
}

func (pthis*cross_linker_front_guard)get_coord()float32 {
	return pthis.coord_
}
func (pthis*cross_linker_front_guard)get_view_range()float32 {
	return 0
}

func (pthis*cross_linker_front_guard)get_node_type()CROSS_LINK_NODE_TYPE {
	return pthis.node_type_
}

func (pthis*cross_linker_front_guard)get_node_id()int {
	return -1
}


type cross_linker_back_guard struct {
	prev_ cross_linker
	next_ cross_linker

	coord_ float32

	node_type_ CROSS_LINK_NODE_TYPE
}
func (pthis*cross_linker_back_guard)get_prev()cross_linker {
	return pthis.prev_
}
func (pthis*cross_linker_back_guard)set_prev(prev cross_linker) {
	pthis.prev_ = prev
}

func (pthis*cross_linker_back_guard)get_next()cross_linker {
	return pthis.next_
}
func (pthis*cross_linker_back_guard)set_next(next cross_linker) {
	pthis.next_ = next
}

func (pthis*cross_linker_back_guard)get_coord()float32 {
	return pthis.coord_
}
func (pthis*cross_linker_back_guard)get_view_range()float32 {
	return 0
}

func (pthis*cross_linker_back_guard)get_node_type()CROSS_LINK_NODE_TYPE {
	return pthis.node_type_
}

func (pthis*cross_linker_back_guard)get_node_id()int {
	return -1
}


type MAP_CROSS_LINK_NODE map[int]*cross_node

type cross_link_lst struct {
	head_ cross_linker
}
type cross_link_mgr struct {
	link_x_ cross_link_lst
	link_y_ cross_link_lst

	cur_node_id_ int
	map_node_ MAP_CROSS_LINK_NODE

	crs_lnk_if_ cross_link_interface
}


func Cross_link_mgr_constructor(crs_lnk_if cross_link_interface) *cross_link_mgr {
	crs_lnk := &cross_link_mgr{
		link_x_ : cross_link_lst{head_:nil},
		link_y_ : cross_link_lst{head_:nil},
		map_node_ : make(MAP_CROSS_LINK_NODE),
		cur_node_id_ : 0,
		crs_lnk_if_:crs_lnk_if,
	}

	return crs_lnk
}

func (pthis*cross_link_lst)_inter_cross_link_add_before(new_node cross_linker, node_pos cross_linker) {
	if node_pos == nil {
		return
	}
	new_node.set_next(node_pos)
	new_node.set_prev(node_pos.get_prev())
	if node_pos.get_prev() != nil {
		node_pos.get_prev().set_next(new_node)
	}
	node_pos.set_prev(new_node)
}
func (pthis*cross_link_lst)_inter_cross_link_add_after(new_node cross_linker, node_pos cross_linker) {
	if node_pos == nil {
		return
	}
	new_node.set_prev(node_pos)
	new_node.set_next(node_pos.get_next())
	if node_pos.get_next() != nil {
		node_pos.get_next().set_prev(new_node)
	}
	node_pos.set_next(new_node)
}

func (pthis*cross_link_lst)_inter_cross_link_del(node cross_linker) {
	if node.get_prev() != nil {
		node.get_prev().set_next(node.get_next())
	}

	if node.get_next() != nil {
		node.get_next().set_prev(node.get_prev())
	}
}

func (pthis*cross_link_mgr)Cross_link_mgr_add(node_if cross_link_node_if) int {
	if node_if == nil {
		return -1
	}

	pthis.cur_node_id_++
	new_node_id := pthis.cur_node_id_

	x_new_node := &cross_linker_node{prev_:nil, next_:nil,
		front_ : nil, back_ : nil,
		coord_ : node_if.get_node_x(),
		view_range_ : node_if.get_view_range(),
		node_type_:CROSS_LINK_NODE_TYPE_node,
		node_id_: new_node_id,
	}
	y_new_node := &cross_linker_node{prev_:nil, next_:nil,
		front_ : nil, back_ : nil,
		coord_ : node_if.get_node_y(),
		view_range_ : node_if.get_view_range(),
		node_type_:CROSS_LINK_NODE_TYPE_node,
		node_id_: new_node_id,
	}

	new_node := &cross_node{
		x_node_ : x_new_node,
		y_node_ : y_new_node,
		node_id_ : new_node_id,
		node_if_ : node_if,
	}

	pthis.map_node_[new_node_id] = new_node

	// sort by xy increase
	map_x_view := make(map[int]int)
	// link x
	if pthis.link_x_.head_ == nil {
		pthis.link_x_.head_ = x_new_node
	} else {
		cur_node := pthis.link_x_.head_
		cur_node_pre := cur_node.get_prev()
		for ;cur_node != nil; {
			if cur_node.get_coord() <= x_new_node.coord_ {
				cur_node_pre = cur_node
				cur_node = cur_node.get_next()
				switch cur_node.get_node_type() {
				case CROSS_LINK_NODE_TYPE_front_guard:
					// in cur_node x view
					map_x_view[cur_node.get_node_id()] = cur_node.get_node_id()
				case CROSS_LINK_NODE_TYPE_back_guard:
					// out cur_node x view
					delete(map_x_view, cur_node.get_node_id())
				}
				continue
			} else {
				pthis.link_x_._inter_cross_link_add_before(x_new_node, cur_node)
				break
			}
		}
		if cur_node == nil {
			if cur_node_pre != nil {
				lin_common.LogDebug("link y tail", new_node_id)
				pthis.link_x_._inter_cross_link_add_after(x_new_node, cur_node_pre)
			} else {
				lin_common.LogErr("can't link to y", new_node_id)
			}
		}
	}


	// link y
	map_y_view := make(map[int]int)
	if pthis.link_y_.head_ == nil {
		pthis.link_y_.head_ = y_new_node
	} else {
		cur_node := pthis.link_y_.head_
		cur_node_pre := cur_node.get_prev()
		for ;cur_node != nil; {
			if cur_node.get_coord() <= y_new_node.coord_ {
				cur_node_pre = cur_node
				cur_node = cur_node.get_next()
				switch cur_node.get_node_type() {
				case CROSS_LINK_NODE_TYPE_front_guard:
					// in cur_node y view
					map_y_view[cur_node.get_node_id()] = cur_node.get_node_id()
				case CROSS_LINK_NODE_TYPE_back_guard:
					// out cur_node y view
					delete(map_y_view, cur_node.get_node_id())
				}
				continue
			} else {
				pthis.link_y_._inter_cross_link_add_before(y_new_node, cur_node)
				break
			}
		}
		if cur_node == nil {
			if cur_node_pre != nil {
				lin_common.LogDebug("link y tail", new_node_id)
				pthis.link_x_._inter_cross_link_add_after(x_new_node, cur_node_pre)
			} else {
				lin_common.LogErr("can't link to y", new_node_id)
			}
		}
	}


	// check map_x_view map_y_view
	new_node.map_view_by_ = make(map[int]int)
	for k, _ := range map_x_view {
		_,ok := map_y_view[k]
		if ok {
			new_node.map_view_by_[k] = k
		}
	}

	// notify the map_xy_view, the new_node in the view
	for k, _ := range new_node.map_view_by_ {
		pthis.crs_lnk_if_.ntf_node_in_view(k, new_node_id)
	}

	pthis._inter_cross_link_link_guard(new_node)

	return new_node_id
}

func(pthis*cross_link_mgr)_inter_cross_link_link_guard(new_node *cross_node) {

	if pthis.link_x_.head_ == nil {
		lin_common.LogErr("err, x link head is nil")
	}
	if pthis.link_y_.head_ == nil {
		lin_common.LogErr("err, y link head is nil")
	}

	x_new_node := new_node.x_node_
	y_new_node := new_node.y_node_
	node_if := new_node.node_if_

	// add trigger node
	map_x_view := make(map[int]int)
	if pthis.link_x_.head_ != nil {
		// add front and back guard to link
		x_new_node.front_ = &cross_linker_front_guard{
			prev_ : nil,
			next_ : x_new_node,
			coord_ : x_new_node.coord_ - node_if.get_view_range(),
			node_type_ : CROSS_LINK_NODE_TYPE_front_guard,
		}
		cur_node := x_new_node.get_prev()
		for ;cur_node != nil; {
			if x_new_node.front_.coord_ < cur_node.get_coord() {
				cur_node = cur_node.get_prev()
				continue
			} else {
				pthis.link_x_._inter_cross_link_add_after(x_new_node.front_, cur_node)
				break
			}
		}

		x_new_node.back_ = &cross_linker_back_guard{
			prev_ : x_new_node,
			next_ : nil,
			coord_ : x_new_node.coord_ + node_if.get_view_range(),
			node_type_ : CROSS_LINK_NODE_TYPE_back_guard,
		}
		cur_node = x_new_node.get_next()
		for ;cur_node != nil; {
			if x_new_node.back_.coord_ > cur_node.get_coord() {
				cur_node = cur_node.get_next()
				continue
			} else {
				pthis.link_x_._inter_cross_link_add_before(x_new_node.back_, cur_node)
				break
			}
		}

		// gather node between front and back guard
		cur_node = x_new_node.front_.get_next()
		for ;cur_node!=nil && cur_node != x_new_node.back_; {
			map_x_view[cur_node.get_node_id()] = cur_node.get_node_id()
			cur_node = cur_node.get_next()
		}
	}

	map_y_view := make(map[int]int)
	if pthis.link_y_.head_ != nil {
		// add front and back guard to link
		y_new_node.front_ = &cross_linker_front_guard{
			prev_ : nil,
			next_ : y_new_node,
			coord_ : y_new_node.coord_ - node_if.get_view_range(),
			node_type_ : CROSS_LINK_NODE_TYPE_front_guard,
		}
		cur_node := y_new_node.get_prev()
		for ;cur_node != nil; {
			if y_new_node.front_.coord_ < cur_node.get_coord() {
				cur_node = cur_node.get_prev()
				continue
			} else {
				pthis.link_y_._inter_cross_link_add_before(y_new_node.front_, cur_node)
				break
			}
		}

		y_new_node.back_ = &cross_linker_back_guard{
			prev_ : y_new_node,
			next_ : nil,
			coord_ : y_new_node.coord_ + node_if.get_view_range(),
			node_type_ : CROSS_LINK_NODE_TYPE_back_guard,
		}
		cur_node = y_new_node.get_next()
		for ;cur_node != nil; {
			if y_new_node.back_.coord_ > cur_node.get_coord() {
				cur_node = cur_node.get_next()
				continue
			} else {
				pthis.link_y_._inter_cross_link_add_before(y_new_node.back_, cur_node)
				break
			}
		}

		// gather node between front and back guard
		cur_node = y_new_node.front_.get_next()
		for ;cur_node!=nil && cur_node != y_new_node.back_; {
			map_y_view[cur_node.get_node_id()] = cur_node.get_node_id()
			cur_node = cur_node.get_next()
		}
	}

	// cal intersection, get the new node view
	new_node.map_view_ = make(map[int]int)
	for k, _ := range map_x_view {
		_,ok := map_y_view[k]
		if ok {
			new_node.map_view_[k] = k
		}
	}

	for k, _ := range new_node.map_view_ {
		pthis.crs_lnk_if_.ntf_node_in_view(new_node.node_id_, k)
	}
}

func (pthis*cross_link_mgr)Cross_link_mgr_del(node_id int) {
	node, _ := pthis.map_node_[node_id]
	if node == nil {
		return
	}
	delete(pthis.map_node_, node_id)

	for k, _ := range node.map_view_ {
		pthis.crs_lnk_if_.ntf_node_out_view(node.node_id_, k)
	}
	for k, _ := range node.map_view_by_ {
		pthis.crs_lnk_if_.ntf_node_out_view(k, node.node_id_)
		node_view_by, _ := pthis.map_node_[k]
		if node_view_by == nil {
			continue
		}
		delete(node_view_by.map_view_by_, node.node_id_)
	}

	pthis.link_y_._inter_cross_link_del(node.x_node_)
	pthis.link_y_._inter_cross_link_del(node.y_node_)
}

func (pthis*cross_link_mgr)Cross_link_mgr_update_pos(node_id int, coord_x float32, coord_y float32) {
	node, _ := pthis.map_node_[node_id]
	if node == nil {
		return
	}

	old_x := node.x_node_.coord_
	old_y := node.y_node_.coord_

	node.x_node_.coord_ = coord_x
	node.y_node_.coord_ = coord_y


	// x move
	x_node := node.x_node_
	x_map := make(map[int]int)
	for k,v := range node.map_view_by_ {
		x_map[k] = v
	}
	if coord_x < old_x {
		x_cur_node := x_node.prev_
		var x_cur_node_last cross_linker = nil
		pthis.link_y_._inter_cross_link_del(node.x_node_)
		for ;x_cur_node!=nil; {
			if x_node.coord_ < x_cur_node.get_coord() {
				x_cur_node_last = x_cur_node
				x_cur_node = x_cur_node.get_prev()
				switch x_cur_node.get_node_type() {
				case CROSS_LINK_NODE_TYPE_front_guard:
					delete(x_map, x_cur_node.get_node_id())
				case CROSS_LINK_NODE_TYPE_back_guard:
					x_map[x_cur_node.get_node_id()] = x_cur_node.get_node_id()
				}
			} else {
				pthis.link_x_._inter_cross_link_add_after(x_node, x_cur_node)
				break
			}
		}
		if x_cur_node == nil {
			if x_cur_node_last != nil {
				lin_common.LogDebug("link x head", node_id)
				pthis.link_x_._inter_cross_link_add_before(node.x_node_, x_cur_node_last)
			} else {
				lin_common.LogErr("can't link to x", node_id)
			}
		}
	} else if coord_x > old_x {
		x_cur_node := x_node.next_
		var x_cur_node_last cross_linker = nil
		pthis.link_y_._inter_cross_link_del(node.x_node_)
		for ;x_cur_node!=nil; {
			if x_node.coord_ > x_cur_node.get_coord() {
				x_cur_node_last = x_cur_node
				x_cur_node = x_cur_node.get_next()
				switch x_cur_node.get_node_type() {
				case CROSS_LINK_NODE_TYPE_front_guard:
					x_map[x_cur_node.get_node_id()] = x_cur_node.get_node_id()
				case CROSS_LINK_NODE_TYPE_back_guard:
					delete(x_map, x_cur_node.get_node_id())
				}
			} else {
				pthis.link_x_._inter_cross_link_add_before(x_node, x_cur_node)
				break
			}
		}
		if x_cur_node == nil {
			if x_cur_node_last != nil {
				lin_common.LogDebug("link x tail", node_id)
				pthis.link_x_._inter_cross_link_add_after(node.x_node_, x_cur_node_last)
			} else {
				lin_common.LogErr("can't link to x", node_id)
			}
		}
	}


	// y move
	y_node := node.y_node_
	y_map := make(map[int]int)
	for k,v := range node.map_view_by_ {
		y_map[k] = v
	}
	if coord_y < old_y {
		y_cur_node := y_node.prev_
		pthis.link_y_._inter_cross_link_del(node.y_node_)
		for ;y_cur_node!=nil; {
			if y_cur_node.get_coord() <= x_node.coord_ {
				pthis.link_y_._inter_cross_link_add_after(y_node, y_cur_node)
				break
			} else {
				switch y_cur_node.get_node_type() {
				case CROSS_LINK_NODE_TYPE_front_guard:
					delete(y_map, y_cur_node.get_node_id())
				case CROSS_LINK_NODE_TYPE_back_guard:
					y_map[y_cur_node.get_node_id()] = y_cur_node.get_node_id()
				}
			}
		}
		// todo : if y_cur_node = nil
	} else if coord_y > old_y {
		y_cur_node := y_node.next_
		pthis.link_y_._inter_cross_link_del(node.y_node_)
		for ;y_cur_node!=nil; {
			if y_cur_node.get_coord() >= x_node.coord_ {
				pthis.link_y_._inter_cross_link_add_before(y_node, y_cur_node)
				break
			} else {
				switch y_cur_node.get_node_type() {
				case CROSS_LINK_NODE_TYPE_front_guard:
					y_map[y_cur_node.get_node_id()] = y_cur_node.get_node_id()
				case CROSS_LINK_NODE_TYPE_back_guard:
					delete(y_map, y_cur_node.get_node_id())
				}
			}
		}
		// todo : if y_cur_node = nil
	}

	// final view by
	node.map_view_by_ = make(map[int]int)
	for k, _ := range x_map {
		_,ok := y_map[k]
		if ok {
			node.map_view_by_[k] = k
		}
	}


	// process front x guard move
	x_front_guard := node.x_node_.front_
	old_front_x := x_front_guard.coord_
	x_front_guard.coord_ = node.x_node_.coord_ - node.node_if_.get_view_range()
	if x_front_guard.coord_ <  old_front_x {
		cur_node := x_front_guard.prev_
		pthis.link_y_._inter_cross_link_del(x_front_guard)
		for ;cur_node != nil; {
			if cur_node.get_coord() <=  x_front_guard.coord_ {
				pthis.link_x_._inter_cross_link_add_after(x_front_guard, cur_node)
				break
			}
		}
	} else if x_front_guard.coord_ >  old_front_x {
		cur_node := x_front_guard.next_
		pthis.link_y_._inter_cross_link_del(x_front_guard)
		for ;cur_node != nil; {
			if cur_node.get_coord() >=  x_front_guard.coord_ {
				pthis.link_x_._inter_cross_link_add_before(x_front_guard, cur_node)
				break
			}
		}
	}

	// todo: process back x move
	x_back_guard := node.x_node_.back_
	old_back_x := x_back_guard.coord_
	x_back_guard.coord_ = node.x_node_.coord_ + node.node_if_.get_view_range()
	if x_back_guard.coord_ < old_back_x {

	} else if x_back_guard.coord_ > old_back_x {

	}

	// todo: process front y move
	// todo: process back y move
}
