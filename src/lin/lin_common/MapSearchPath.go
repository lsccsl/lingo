package lin_common

import (
	"math"
	"sort"
)

const (
	SEARCH_NEIGHBOR_up    = 0
	SEARCH_NEIGHBOR_down  = 1
	SEARCH_NEIGHBOR_left  = 2
	SEARCH_NEIGHBOR_right = 3

	SEARCH_NEIGHBOR_up_left    = 4
	SEARCH_NEIGHBOR_up_right   = 5
	SEARCH_NEIGHBOR_down_left  = 6
	SEARCH_NEIGHBOR_down_right = 7

	SEARCH_NEIGHBOR_MAX = 8
)



type SearchNode struct {
	pos Coord2d
	parent *SearchNode
	neighbor [SEARCH_NEIGHBOR_MAX]*SearchNode
	startWeight int
	endWeight int
	totalWeight int
}

type MAP_HISTORY_PATH map[Coord2d]*SearchNode
type SearchMgr struct {
	root *SearchNode
	nodes []*SearchNode
	mapHistoryPath MAP_HISTORY_PATH
}
func (pthis*SearchMgr)Len() int {
	return len(pthis.nodes)
}
func (pthis*SearchMgr)Less(i, j int) bool {
	if i < 0 || i >= len(pthis.nodes) {
		return false
	}
	if j < 0 || j >= len(pthis.nodes) {
		return false
	}

	node_i := pthis.nodes[i]
	node_j := pthis.nodes[j]
	if node_i.totalWeight > node_j.totalWeight {
		return true
	}
	return false
}
func (pthis*SearchMgr)Swap(i, j int) {
	if i < 0 || i >= len(pthis.nodes) {
		return
	}
	if j < 0 || j >= len(pthis.nodes) {
		return
	}
	nodeTmp := pthis.nodes[i]
	pthis.nodes[i] = pthis.nodes[j]
	pthis.nodes[j] = nodeTmp
}
func (pthis*SearchMgr)addNode(node *SearchNode) {
	pthis.nodes = append(pthis.nodes, node)
	sort.Sort(pthis)
}
func (pthis*SearchMgr)getNearestNode() *SearchNode {
	if len(pthis.nodes) == 0 {
		return nil
	}
	lastIdx := len(pthis.nodes) - 1
	node := pthis.nodes[lastIdx]
	pthis.nodes = pthis.nodes[:lastIdx]
	return node
}

const (
	ASTAR_WEIGHT_slash = 14
	ASTAR_WEIGHT_straight = 10
	ASTAR_WEIGHT_scale = 10
)

func calAStarWeight(src Coord2d, dst Coord2d) int {
	//欧式
	//return (src.X - dst.X) * (src.X - dst.X) + (src.Y - dst.Y) * (src.Y - dst.Y)

	//chebyshev
	//return int(math.Max(math.Abs(float64(src.X - dst.X)), math.Abs(float64(src.Y - dst.Y)))) * ASTAR_WEIGHT_scale

	//曼哈顿
	return int(math.Abs(float64(src.X - dst.X)) + math.Abs(float64(src.Y - dst.Y))) * JPS_WEIGHT_scale
}



func (pthis*MapData)PathSearchAStart(src Coord2d, dst Coord2d) (path []Coord2d, search * SearchMgr) {

	// search around by weight,
	// weight = end_weight + start_weight,
	// end_weight = abs(x_diff) + abs(y_diff), start_weight += 14 or 10

	startNode := &SearchNode{
		pos:src,
		startWeight:0,
		endWeight:calAStarWeight(src, dst),
		parent:nil,
	}

	searchMgr := &SearchMgr{
		root:startNode,
		mapHistoryPath: make(MAP_HISTORY_PATH),
	}
	searchMgr.mapHistoryPath[startNode.pos] = startNode

	searchMgr.addNode(startNode)

	bFind := false
	node := startNode
	for {
		node = searchMgr.getNearestNode()
		if node == nil {
			LogDebug("fail find")
			break
		}

		for i := 0; i < SEARCH_NEIGHBOR_MAX; i ++ {
			curPos := node.pos
			startWeight := 0
			switch i {
			case SEARCH_NEIGHBOR_up:
				curPos.Y -= 1
				startWeight = ASTAR_WEIGHT_straight
			case SEARCH_NEIGHBOR_down:
				curPos.Y += 1
				startWeight = ASTAR_WEIGHT_straight
			case SEARCH_NEIGHBOR_left:
				curPos.X -= 1
				startWeight = ASTAR_WEIGHT_straight
			case SEARCH_NEIGHBOR_right:
				curPos.X += 1
				startWeight = ASTAR_WEIGHT_straight

			case SEARCH_NEIGHBOR_up_left:
				curPos.X -= 1
				curPos.Y -= 1
				startWeight = ASTAR_WEIGHT_slash
			case SEARCH_NEIGHBOR_up_right:
				curPos.X += 1
				curPos.Y -= 1
				startWeight = ASTAR_WEIGHT_slash
			case SEARCH_NEIGHBOR_down_left:
				curPos.X -= 1
				curPos.Y += 1
				startWeight = ASTAR_WEIGHT_slash
			case SEARCH_NEIGHBOR_down_right:
				curPos.X += 1
				curPos.Y += 1
				startWeight = ASTAR_WEIGHT_slash
			}

			if curPos.X == dst.X && curPos.Y == dst.Y {
				LogDebug("suc find")
				bFind = true
				break
			}
			if pthis.IsBlock(curPos.X, curPos.Y) {
				continue
			}
			_, ok := searchMgr.mapHistoryPath[curPos]
			if ok {
				continue
			}

			nodeNeighbor := &SearchNode{
				pos:curPos,
				startWeight:startWeight + node.startWeight,
				endWeight:calAStarWeight(curPos, dst),
				parent:node,
			}
			nodeNeighbor.totalWeight = nodeNeighbor.startWeight + nodeNeighbor.endWeight
			searchMgr.mapHistoryPath[nodeNeighbor.pos] = nodeNeighbor

			node.neighbor[i] = nodeNeighbor
			searchMgr.addNode(nodeNeighbor)
		}

		if bFind {
			break
		}
	}

	if node != nil {
		for ;node != nil; {
			path = append(path, node.pos)
			node = node.parent
		}
		return path, searchMgr
	}

	return nil, nil
}

