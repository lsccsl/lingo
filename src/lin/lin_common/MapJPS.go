package lin_common

var dirStraight = [4]Coord2d{
	{1,0},
	{-1,0},
	{0,-1},
	{0,1},
}
var dirSlash = [4]Coord2d{
	{1,1},
	{-1,1},
	{1,-1},
	{-1,-1},
}

func (pthis*MapData)PathJPS(src Coord2d, dst Coord2d) (path []Coord2d) {

	/*
		(1)节点 A 是起点、终点.
		(2)节点A 至少有一个强迫邻居. 以横向x为例,说明此时被阻挡了,并且越过这个点之后,有新的连通区域
		(3)父节点在斜方向(斜向搜索)，节点A的水平或者垂直方向上有满足 (1)、(2) 的节点 用于转向
	*/



	return nil
}
