package lin_common

const (
	JSP_WEIGHT_slash = 14
	JSP_WEIGHT_straight = 10
)

const (
	JPS_up    = 0
	JPS_down  = 1
	JPS_left  = 2
	JPS_right = 3

	JPS_up_left    = 4
	JPS_up_right   = 5
	JPS_down_left  = 6
	JPS_down_right = 7
)

var dirStraight = [4]Coord2d{
	{0,-1},
	{0,1},
	{-1,0},
	{1,0},
}
var dirSlash = [4]Coord2d{
	{-1,-1},
	{1,-1},
	{-1,1},
	{1,1},
}

type JPSNode struct {
	parent *JPSNode
	pos Coord2d

	startWeight int
	endWeight int
	totalWeight int
}

func (pthis*MapData)PathJPS(src Coord2d, dst Coord2d) (path []Coord2d) {

	/*
		(1)节点 A 是起点、终点.
		(2)节点A 至少有一个强迫邻居. 以横向x为例,说明此时被阻挡了,并且越过这个点之后,有新的连通区域
		(3)父节点在斜方向(斜向搜索)，节点A的水平或者垂直方向上有满足 (1)、(2) 的节点 用于转向
	*/

	startNode := &JPSNode{
		parent: nil,
		pos:src,
		endWeight:calEndWeight(src, dst),
	}
	startNode.totalWeight = startNode.endWeight

	//if find jp/block/outbound map, stop
	curPos := src
	curNode := startNode

	for i := 0; i < 4; i ++ {
		dir := dirStraight[i]
		for {
			searchPos := Coord2d{curPos.X + dir.X, curPos.Y + dir.Y}
			if pthis.IsBlock(searchPos.X, searchPos.Y) {
				break
			}

			bStop := false

			switch i {
			case JPS_up:

			case JPS_down:

			case JPS_left:

			case JPS_right:
				posUp   := Coord2d{searchPos.X, searchPos.Y - 1}
				posDown := Coord2d{searchPos.X, searchPos.Y + 1}

				if pthis.CoordIsBlock(posUp) {
					posN := Coord2d{posUp.X + 1, posUp.Y}
					if !pthis.CoordIsBlock(posN) {
						bStop = true

						// todo add pos to open list
						jp := JPSNode{
							parent: curNode,
							pos:searchPos,
							endWeight:calEndWeight(searchPos, dst),
							startWeight:curNode.startWeight + JSP_WEIGHT_straight,
						}
						jp.totalWeight = jp.startWeight + jp.endWeight
					}
				}
				if pthis.CoordIsBlock(posDown) {
					posN := Coord2d{posDown.X + 1, posDown.Y}
					if !pthis.CoordIsBlock(posN) {
						bStop = true

						// todo add pos to open list
					}
				}
			}

			if bStop {
				break
			}
		}
	}

	return nil
}
