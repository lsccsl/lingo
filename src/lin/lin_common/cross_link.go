package lin_common

type CROSS_LINK_NODE_TYPE int
const (
	CROSS_LINK_NODE_TYPE_node        CROSS_LINK_NODE_TYPE = 1
	CROSS_LINK_NODE_TYPE_front_guard CROSS_LINK_NODE_TYPE = 2
	CROSS_LINK_NODE_TYPE_back_guard  CROSS_LINK_NODE_TYPE = 3
)

type Cross_link_interface interface {
	Ntf_node_in_view(node_id int, node_id_in_viewed int)
	Ntf_node_out_view(node_id int, node_id_out_view int)
}

type Cross_link_node_if interface {
	Get_node_x()float32
	Get_node_y()float32
	Get_node_data()interface{}
	Get_view_range()float32
}

type cross_linker_inf interface {
	get_prev() cross_linker_inf
	set_prev(prev cross_linker_inf)

	get_next() cross_linker_inf
	set_next(next cross_linker_inf)

	get_coord()float32 // x or y
	get_view_range()float32

	get_node_type() CROSS_LINK_NODE_TYPE

	get_node_id()int
}

type MAP_NODE_ID map[int]int
type cross_linker_node struct {
	prev_ cross_linker_inf
	next_ cross_linker_inf

	front_ *cross_linker_front_guard
	back_  *cross_linker_back_guard

	coord_ float32
	view_range_ float32

	node_type_ CROSS_LINK_NODE_TYPE

	node_id_ int
}
type cross_node struct {
	x_node_ *cross_linker_node
	y_node_ *cross_linker_node

	node_id_ int
	node_if_ Cross_link_node_if

	map_view_by_ MAP_NODE_ID // view by
	map_view_    MAP_NODE_ID // view
}
func (pthis*cross_linker_node)get_prev() cross_linker_inf {
	return pthis.prev_
}
func (pthis*cross_linker_node)set_prev(prev cross_linker_inf) {
	pthis.prev_ = prev
}

func (pthis*cross_linker_node)get_next() cross_linker_inf {
	return pthis.next_
}
func (pthis*cross_linker_node)set_next(next cross_linker_inf) {
	pthis.next_ = next
}

func (pthis*cross_linker_node)get_coord()float32 {
	return pthis.coord_
}
func (pthis*cross_linker_node)get_view_range()float32 {
	return pthis.view_range_
}

func (pthis*cross_linker_node)get_node_type() CROSS_LINK_NODE_TYPE {
	return pthis.node_type_
}

func (pthis*cross_linker_node)get_node_id()int {
	return pthis.node_id_
}


type cross_linker_front_guard struct {
	prev_ cross_linker_inf
	next_ cross_linker_inf

	coord_ float32

	node_type_ CROSS_LINK_NODE_TYPE
	node_id_ int
}
func (pthis*cross_linker_front_guard)get_prev() cross_linker_inf {
	return pthis.prev_
}
func (pthis*cross_linker_front_guard)set_prev(prev cross_linker_inf) {
	pthis.prev_ = prev
}

func (pthis*cross_linker_front_guard)get_next() cross_linker_inf {
	return pthis.next_
}
func (pthis*cross_linker_front_guard)set_next(next cross_linker_inf) {
	pthis.next_ = next
}

func (pthis*cross_linker_front_guard)get_coord()float32 {
	return pthis.coord_
}
func (pthis*cross_linker_front_guard)get_view_range()float32 {
	return 0
}

func (pthis*cross_linker_front_guard)get_node_type() CROSS_LINK_NODE_TYPE {
	return pthis.node_type_
}

func (pthis*cross_linker_front_guard)get_node_id()int {
	return pthis.node_id_
}


type cross_linker_back_guard struct {
	prev_ cross_linker_inf
	next_ cross_linker_inf

	coord_ float32

	node_type_ CROSS_LINK_NODE_TYPE
	node_id_ int
}
func (pthis*cross_linker_back_guard)get_prev() cross_linker_inf {
	return pthis.prev_
}
func (pthis*cross_linker_back_guard)set_prev(prev cross_linker_inf) {
	pthis.prev_ = prev
}

func (pthis*cross_linker_back_guard)get_next() cross_linker_inf {
	return pthis.next_
}
func (pthis*cross_linker_back_guard)set_next(next cross_linker_inf) {
	pthis.next_ = next
}

func (pthis*cross_linker_back_guard)get_coord()float32 {
	return pthis.coord_
}
func (pthis*cross_linker_back_guard)get_view_range()float32 {
	return 0
}

func (pthis*cross_linker_back_guard)get_node_type() CROSS_LINK_NODE_TYPE {
	return pthis.node_type_
}

func (pthis*cross_linker_back_guard)get_node_id()int {
	return pthis.node_id_
}


type MAP_CROSS_LINK_NODE map[int]*cross_node

type cross_link_lst struct {
	head_ cross_linker_inf
	tail_ cross_linker_inf
}
type cross_link_mgr struct {
	link_x_ cross_link_lst
	link_y_ cross_link_lst

	cur_node_id_ int
	map_node_    MAP_CROSS_LINK_NODE

	crs_lnk_if_ Cross_link_interface
}

func (pthis*cross_link_lst)_inter_cross_link_add_before(new_node cross_linker_inf, node_pos cross_linker_inf) {
	if node_pos == nil {
		// add before head
		if pthis.head_ != nil {
			new_node.set_next(pthis.head_)
			new_node.set_prev(nil)
			pthis.head_.set_prev(new_node)
		}
		pthis.head_ = new_node
		if pthis.tail_ == nil {
			pthis.tail_ = new_node
		}
		return
	} else {
		new_node.set_next(node_pos)
		new_node.set_prev(node_pos.get_prev())
		if node_pos.get_prev() != nil {
			node_pos.get_prev().set_next(new_node)
		}
		node_pos.set_prev(new_node)
		if new_node.get_prev() == nil {
			pthis.head_ = new_node
		}
	}
}
func (pthis*cross_link_lst)_inter_cross_link_add_after(new_node cross_linker_inf, node_pos cross_linker_inf) {
	if node_pos == nil {
		// add after tail
		if pthis.tail_ != nil {
			new_node.set_prev(pthis.tail_)
			new_node.set_next(nil)
			pthis.tail_.set_next(new_node)
		}
		pthis.tail_ = new_node
		if pthis.head_ == nil {
			pthis.head_ = new_node
		}
		return
	} else {
		new_node.set_prev(node_pos)
		new_node.set_next(node_pos.get_next())
		if node_pos.get_next() != nil {
			node_pos.get_next().set_prev(new_node)
		}
		node_pos.set_next(new_node)
		if new_node.get_next() == nil {
			pthis.tail_ = new_node
		}
	}
}

func (pthis*cross_link_lst)_inter_cross_link_del(node cross_linker_inf) {
	if node.get_prev() != nil {
		node.get_prev().set_next(node.get_next())
	} else {
		pthis.head_ = node.get_next()
	}
	if node.get_next() != nil {
		node.get_next().set_prev(node.get_prev())
	} else {
		pthis.tail_ = node.get_prev()
	}
}
func (pthis*cross_link_lst)_inter_cross_link_empty() bool {
	if pthis.head_ == nil {
		return true
	}
	return false
}

func Cross_link_mgr_constructor(crs_lnk_if Cross_link_interface) *cross_link_mgr {
	crs_lnk := &cross_link_mgr{
		link_x_ :      cross_link_lst{head_: nil,tail_:nil},
		link_y_ :      cross_link_lst{head_: nil,tail_:nil},
		map_node_ :    make(MAP_CROSS_LINK_NODE),
		cur_node_id_ : 0,
		crs_lnk_if_:   crs_lnk_if,
	}

	return crs_lnk
}

func (pthis*cross_link_mgr)Cross_link_mgr_add(node_if Cross_link_node_if) int {
	if node_if == nil {
		return -1
	}

	pthis.cur_node_id_++
	new_node_id := pthis.cur_node_id_

	x_new_node := &cross_linker_node{prev_: nil, next_:nil,
		front_ : nil, back_ : nil,
		coord_ :      node_if.Get_node_x(),
		view_range_ : node_if.Get_view_range(),
		node_type_:   CROSS_LINK_NODE_TYPE_node,
		node_id_:     new_node_id,
	}
	y_new_node := &cross_linker_node{prev_: nil, next_:nil,
		front_ : nil, back_ : nil,
		coord_ :      node_if.Get_node_y(),
		view_range_ : node_if.Get_view_range(),
		node_type_:   CROSS_LINK_NODE_TYPE_node,
		node_id_:     new_node_id,
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
		pthis.link_x_._inter_cross_link_add_before(x_new_node,nil)
	} else {
		cur_node := pthis.link_x_.head_
		for ;cur_node != nil; {
			if cur_node.get_coord() <= x_new_node.coord_ {
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
			LogDebug("link y tail", new_node_id)
			pthis.link_x_._inter_cross_link_add_after(x_new_node,nil)
		}
	}


	// link y
	map_y_view := make(map[int]int)
	if pthis.link_y_.head_ == nil {
		pthis.link_y_._inter_cross_link_add_before(y_new_node,nil)
	} else {
		cur_node := pthis.link_y_.head_
		for ;cur_node != nil; {
			if cur_node.get_coord() <= y_new_node.coord_ {
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
			LogDebug("link y tail", new_node_id)
			pthis.link_y_._inter_cross_link_add_after(y_new_node, nil)
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
		pthis.crs_lnk_if_.Ntf_node_in_view(k, new_node_id)
	}

	pthis._inter_cross_link_link_guard(new_node)

	return new_node_id
}

func (pthis*cross_node)_inter_gen_view() map[int]int {
	// gather node between front and back guard
	map_x_view := make(map[int]int)
	x_node := pthis.x_node_
	cur_node := x_node.front_.get_next()
	for ;cur_node!=nil && cur_node != x_node.back_; {
		map_x_view[cur_node.get_node_id()] = cur_node.get_node_id()
		cur_node = cur_node.get_next()
	}

	// gather node between front and back guard
	map_y_view := make(map[int]int)
	y_node := pthis.y_node_
	cur_node = y_node.front_.get_next()
	for ;cur_node!=nil && cur_node != y_node.back_; {
		map_y_view[cur_node.get_node_id()] = cur_node.get_node_id()
		cur_node = cur_node.get_next()
	}

	// cal intersection, get the new node view
	map_view := make(map[int]int)
	for k, _ := range map_x_view {
		_,ok := map_y_view[k]
		if ok {
			map_view[k] = k
		}
	}
	return map_view
}

func(pthis*cross_link_mgr)_inter_cross_link_link_guard(new_node *cross_node) {

	if pthis.link_x_._inter_cross_link_empty() {
		LogErr("err, x link head is nil")
	}
	if pthis.link_y_._inter_cross_link_empty() {
		LogErr("err, y link head is nil")
	}

	x_new_node := new_node.x_node_
	y_new_node := new_node.y_node_
	node_if := new_node.node_if_

	// add trigger node
	if pthis.link_x_.head_ != nil {
		// add front and back guard to link
		x_new_node.front_ = &cross_linker_front_guard{
			prev_ :      nil,
			next_ :      nil,
			coord_ :     x_new_node.coord_ - node_if.Get_view_range(),
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
		if cur_node == nil {
			pthis.link_x_._inter_cross_link_add_before(x_new_node.front_, nil)
		}

		x_new_node.back_ = &cross_linker_back_guard{
			prev_ :      nil,
			next_ :      nil,
			coord_ :     x_new_node.coord_ + node_if.Get_view_range(),
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
		if cur_node == nil {
			pthis.link_x_._inter_cross_link_add_after(x_new_node.back_, nil)
		}
	}

	if pthis.link_y_.head_ != nil {
		// add front and back guard to link
		y_new_node.front_ = &cross_linker_front_guard{
			prev_ :      nil,
			next_ :      nil,
			coord_ :     y_new_node.coord_ - node_if.Get_view_range(),
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
		if cur_node == nil {
			pthis.link_y_._inter_cross_link_add_before(y_new_node.front_, nil)
		}

		y_new_node.back_ = &cross_linker_back_guard{
			prev_ :      nil,
			next_ :      nil,
			coord_ :     y_new_node.coord_ + node_if.Get_view_range(),
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
		if cur_node == nil {
			pthis.link_y_._inter_cross_link_add_after(y_new_node.back_, nil)
		}
	}

	new_node.map_view_ = new_node._inter_gen_view()

	for k, _ := range new_node.map_view_ {
		pthis.crs_lnk_if_.Ntf_node_in_view(new_node.node_id_, k)
	}
}

func (pthis*cross_link_mgr)Cross_link_mgr_del(node_id int) {
	node, _ := pthis.map_node_[node_id]
	if node == nil {
		return
	}
	delete(pthis.map_node_, node_id)

	for k, _ := range node.map_view_by_ {
		pthis.crs_lnk_if_.Ntf_node_out_view(k, node.node_id_)
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
	x_map := make(map[int]int)
	{
		x_node := node.x_node_
		for k, v := range node.map_view_by_ {
			x_map[k] = v
		}
		if coord_x < old_x {
			x_cur_node := x_node.prev_
			pthis.link_y_._inter_cross_link_del(x_node)
			for ; x_cur_node != nil; {
				if x_node.coord_ < x_cur_node.get_coord() {
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
				LogDebug("link x head", node_id)
				pthis.link_x_._inter_cross_link_add_before(x_node, nil)
			}
		} else if coord_x > old_x {
			x_cur_node := x_node.next_
			pthis.link_y_._inter_cross_link_del(x_node)
			for ; x_cur_node != nil; {
				if x_node.coord_ > x_cur_node.get_coord() {
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
				LogDebug("link x tail", node_id)
				pthis.link_x_._inter_cross_link_add_after(x_node, nil)
			}
		}
	}


	// y move
	y_map := make(map[int]int)
	{
		y_node := node.y_node_
		for k,v := range node.map_view_by_ {
			y_map[k] = v
		}
		if coord_y < old_y {
			y_cur_node := y_node.prev_
			pthis.link_y_._inter_cross_link_del(y_node)
			for ;y_cur_node!=nil; {
				if y_node.coord_ < y_cur_node.get_coord() {
					y_cur_node = y_cur_node.get_prev()
					switch y_cur_node.get_node_type() {
					case CROSS_LINK_NODE_TYPE_front_guard:
						delete(y_map, y_cur_node.get_node_id())
					case CROSS_LINK_NODE_TYPE_back_guard:
						y_map[y_cur_node.get_node_id()] = y_cur_node.get_node_id()
					}
					continue
				} else {
					pthis.link_y_._inter_cross_link_add_after(y_node, y_cur_node)
					break
				}
			}
			if y_cur_node == nil {
				LogDebug("link y head", node_id)
				pthis.link_y_._inter_cross_link_add_before(y_node, nil)
			}
		} else if coord_y > old_y {
			y_cur_node := y_node.next_
			pthis.link_y_._inter_cross_link_del(y_node)
			for ;y_cur_node!=nil; {
				if y_cur_node.get_coord() >= y_node.coord_ {
					y_cur_node = y_cur_node.get_next()
					switch y_cur_node.get_node_type() {
					case CROSS_LINK_NODE_TYPE_front_guard:
						y_map[y_cur_node.get_node_id()] = y_cur_node.get_node_id()
					case CROSS_LINK_NODE_TYPE_back_guard:
						delete(y_map, y_cur_node.get_node_id())
					}
					continue
				} else {
					pthis.link_y_._inter_cross_link_add_before(y_node, y_cur_node)
					break
				}
			}
			if y_cur_node == nil {
				LogDebug("link y tail", node_id)
				pthis.link_y_._inter_cross_link_add_after(y_node, nil)
			}
		}
	}

	// final view by
	{
		map_view_by := make(map[int]int)
		for k, _ := range x_map {
			_, ok := y_map[k]
			if ok {
				map_view_by[k] = k
			}
		}
		for k, _ := range node.map_view_by_ {
			_, ok := map_view_by[k]
			if !ok {
				pthis.crs_lnk_if_.Ntf_node_in_view(k, node_id)
			}
		}
		for k, _ := range map_view_by {
			_, ok := node.map_view_by_[k]
			if !ok {
				pthis.crs_lnk_if_.Ntf_node_in_view(k, node_id)
			}
		}
		node.map_view_by_ = map_view_by
	}


	// process front x guard move
	{
		x_front_guard := node.x_node_.front_
		old_front_x := x_front_guard.coord_
		x_front_guard.coord_ = node.x_node_.coord_ - node.node_if_.Get_view_range()
		if x_front_guard.coord_ < old_front_x {
			cur_node := x_front_guard.get_prev()
			pthis.link_x_._inter_cross_link_del(x_front_guard)
			for ; cur_node != nil; {
				if x_front_guard.coord_ < cur_node.get_coord() {
					cur_node = cur_node.get_prev()
					continue
				} else {
					pthis.link_x_._inter_cross_link_add_after(x_front_guard, cur_node)
					break
				}
			}
			if cur_node == nil {
				pthis.link_x_._inter_cross_link_add_before(x_front_guard, nil)
			}
		} else if x_front_guard.coord_ > old_front_x {
			cur_node := x_front_guard.get_next()
			pthis.link_x_._inter_cross_link_del(x_front_guard)
			for ; cur_node != nil; {
				if x_front_guard.coord_ > cur_node.get_coord() {
					cur_node = cur_node.get_next()
					continue
				} else {
					pthis.link_x_._inter_cross_link_add_before(x_front_guard, cur_node)
					break
				}
			}
			if cur_node == nil {
				pthis.link_x_._inter_cross_link_add_after(x_front_guard, nil)
			}
		}
	}

	// process back x move
	{
		x_back_guard := node.x_node_.back_
		old_back_x := x_back_guard.coord_
		x_back_guard.coord_ = node.x_node_.coord_ + node.node_if_.Get_view_range()
		if x_back_guard.coord_ < old_back_x {
			cur_node := x_back_guard.get_prev()
			pthis.link_x_._inter_cross_link_del(x_back_guard)
			for ; cur_node != nil; {
				if x_back_guard.coord_ < cur_node.get_coord() {
					cur_node = cur_node.get_prev()
					continue
				} else {
					pthis.link_x_._inter_cross_link_add_after(x_back_guard, cur_node)
					break
				}
			}
			if cur_node == nil {
				pthis.link_x_._inter_cross_link_add_before(x_back_guard, nil)
			}
		} else if x_back_guard.coord_ > old_back_x {
			cur_node := x_back_guard.get_next()
			pthis.link_x_._inter_cross_link_del(x_back_guard)
			for ; cur_node != nil; {
				if x_back_guard.coord_ > cur_node.get_coord() {
					cur_node = cur_node.get_next()
					continue
				} else {
					pthis.link_x_._inter_cross_link_add_before(x_back_guard, cur_node)
					break
				}
			}
			if cur_node == nil {
				pthis.link_x_._inter_cross_link_add_after(x_back_guard, nil)
			}
		}
	}

	// process front y move
	{
		y_front_guard := node.y_node_.front_
		old_front_y := y_front_guard.coord_
		y_front_guard.coord_ = node.x_node_.coord_ - node.node_if_.Get_view_range()
		if y_front_guard.coord_ < old_front_y {
			cur_node := y_front_guard.get_prev()
			pthis.link_y_._inter_cross_link_del(y_front_guard)
			for ; cur_node != nil; {
				if y_front_guard.coord_ < cur_node.get_coord() {
					cur_node = cur_node.get_prev()
					continue
				} else {
					pthis.link_y_._inter_cross_link_add_after(y_front_guard, cur_node)
					break
				}
			}
			if cur_node == nil {
				pthis.link_y_._inter_cross_link_add_before(y_front_guard, nil)
			}
		} else if y_front_guard.coord_ > old_front_y {
			cur_node := y_front_guard.get_next()
			pthis.link_y_._inter_cross_link_del(y_front_guard)
			for ; cur_node != nil; {
				if y_front_guard.coord_ > cur_node.get_coord() {
					cur_node = cur_node.get_next()
					continue
				} else {
					pthis.link_y_._inter_cross_link_add_before(y_front_guard, cur_node)
					break
				}
			}
			if cur_node == nil {
				pthis.link_y_._inter_cross_link_add_after(y_front_guard, nil)
			}
		}
	}

	// process back y move
	{
		y_back_guard := node.y_node_.back_
		old_back_y := y_back_guard.coord_
		y_back_guard.coord_ = node.x_node_.coord_ + node.node_if_.Get_view_range()
		if y_back_guard.coord_ < old_back_y {
			cur_node := y_back_guard.get_prev()
			pthis.link_y_._inter_cross_link_del(y_back_guard)
			for ; cur_node != nil; {
				if y_back_guard.coord_ < cur_node.get_coord() {
					cur_node = cur_node.get_prev()
					continue
				} else {
					pthis.link_y_._inter_cross_link_add_after(y_back_guard, cur_node)
					break
				}
			}
			if cur_node == nil {
				pthis.link_y_._inter_cross_link_add_before(y_back_guard, nil)
			}
		} else if y_back_guard.coord_ > old_back_y {
			cur_node := y_back_guard.get_next()
			pthis.link_y_._inter_cross_link_del(y_back_guard)
			for ; cur_node != nil; {
				if y_back_guard.coord_ > cur_node.get_coord() {
					cur_node = cur_node.get_next()
					continue
				} else {
					pthis.link_y_._inter_cross_link_add_before(y_back_guard, cur_node)
					break
				}
			}
			if cur_node == nil {
				pthis.link_y_._inter_cross_link_add_after(y_back_guard, nil)
			}
		}
	}

	{
		map_view_new := node._inter_gen_view()
		map_view_old := node.map_view_
		for k, _ := range map_view_old {
			_, ok := map_view_new[k]
			if !ok {
				pthis.crs_lnk_if_.Ntf_node_in_view(node_id, k)
			}
		}
		for k, _ := range map_view_new {
			_, ok := map_view_old[k]
			if !ok {
				pthis.crs_lnk_if_.Ntf_node_out_view(node_id, k)
			}
		}
		node.map_view_ = map_view_new
	}
}
