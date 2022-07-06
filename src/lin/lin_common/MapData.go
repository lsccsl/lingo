package lin_common

import (
	//"container/list"
	"math"
	"sort"
)

const (
	WEIGHT_slash = 14
	WEIGHT_straight = 10
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

type Coord2d struct {
	X int
	Y int
}
type MapData struct {
	wid int
	hei int

	mapBit []uint8

	openNodeMgr SearchOpenNodeMgr
}

func (pthis*MapData)LoadMap(mapFile string)error{
	bmp := Bitmap{}
	err := bmp.ReadBmp(mapFile)
	if err != nil {
		return err
	}

	pthis.mapBit = bmp.BmpData
	pthis.wid = bmp.GetWidth()
	pthis.hei = bmp.GetHeight()

	return nil
}

func (pthis*MapData)GetBitBlock(x int, y int)bool{
	if x < 0 || x >= pthis.wid {
		return true
	}
	if y < 0 || y >= pthis.hei {
		return true
	}

	idx := y * pthis.wid + x
	idxByte := idx / 8
	idxBit := 7 - idx % 8
	posByte := pthis.mapBit[idxByte]
	posBit := posByte & (1 << idxBit)

	return posBit == 0
}


type SearchNode struct {
	pos Coord2d
	parent *SearchNode
	neighbor [SEARCH_NEIGHBOR_MAX]*SearchNode
	startWeight int
	endWeight int
	totalWeight int
}

type MAP_HISTORY_PATH map[Coord2d]*SearchNode
type SearchOpenNodeMgr struct {
	root *SearchNode
	nodes []*SearchNode
	mapHistoryPath MAP_HISTORY_PATH
}
func (pthis*SearchOpenNodeMgr)Len() int {
	return len(pthis.nodes)
}
func (pthis*SearchOpenNodeMgr)Less(i, j int) bool {
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
func (pthis*SearchOpenNodeMgr)Swap(i, j int) {
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
func (pthis*SearchOpenNodeMgr)addNode(node *SearchNode) {
	pthis.nodes = append(pthis.nodes, node)
	sort.Sort(pthis)
}
func (pthis*SearchOpenNodeMgr)getNearestNode() *SearchNode {
	if len(pthis.nodes) == 0 {
		return nil
	}
	lastIdx := len(pthis.nodes) - 1
	node := pthis.nodes[lastIdx]
	pthis.nodes = pthis.nodes[:lastIdx]
	return node
}


func calEndWeight(src Coord2d, dst Coord2d) int {
	return int(math.Abs(float64(src.X - dst.X)) + math.Abs(float64(src.Y - dst.Y)))
}


func (pthis*MapData)PathSearch(src Coord2d, dst Coord2d) (path []Coord2d) {

	// search around by weight,
	// weight = end_weight + start_weight,
	// end_weight = abs(x_diff) + abs(y_diff), start_weight += 14 or 10

	startNode := &SearchNode{
		pos:src,
		startWeight:0,
		endWeight:calEndWeight(src, dst),
		parent:nil,
	}

	searchMgr := &SearchOpenNodeMgr{
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
				curPos.Y += 1
				startWeight = WEIGHT_straight
			case SEARCH_NEIGHBOR_down:
				curPos.Y -= 1
				startWeight = WEIGHT_straight
			case SEARCH_NEIGHBOR_left:
				curPos.X -= 1
				startWeight = WEIGHT_straight
			case SEARCH_NEIGHBOR_right:
				curPos.X += 1
				startWeight = WEIGHT_straight

			case SEARCH_NEIGHBOR_up_left:
				curPos.X -= 1
				curPos.Y -= 1
				startWeight = WEIGHT_slash
			case SEARCH_NEIGHBOR_up_right:
				curPos.X += 1
				curPos.Y -= 1
				startWeight = WEIGHT_slash
			case SEARCH_NEIGHBOR_down_left:
				curPos.X -= 1
				curPos.Y += 1
				startWeight = WEIGHT_slash
			case SEARCH_NEIGHBOR_down_right:
				curPos.X += 1
				curPos.Y += 1
				startWeight = WEIGHT_slash
			}

			if curPos.X == dst.X && curPos.Y == dst.Y {
				LogDebug("suc find")
				bFind = true
				break
			}
			if pthis.GetBitBlock(curPos.X, curPos.Y) {
				continue
			}
			_, ok := searchMgr.mapHistoryPath[curPos]
			if ok {
				continue
			}

			nodeNeighbor := &SearchNode{
				pos:curPos,
				startWeight:startWeight + node.startWeight,
				endWeight:calEndWeight(curPos, dst),
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
		return path
	}

	return nil
}

func (pthis*MapData)DumpMap(strMapFile string, path []Coord2d, src * Coord2d , dst * Coord2d) {

	dataLen := len(pthis.mapBit)
	tmpBMP := make([]uint8, dataLen * 24)

	for idx, val := range pthis.mapBit {
		for i := 7; i >= 0 ; i -- {
			tmp := val & (1 << i)
			var clr uint8 = 0
			if tmp != 0 {
				clr = 0xff
			}
			newIdx := idx * 24 + (7 - i) * 3
			tmpBMP[newIdx + 0] = clr
			tmpBMP[newIdx + 1] = clr
			tmpBMP[newIdx + 2] = clr
		}
	}

	if path != nil {
		for _, val := range path {
			idx := (val.Y*pthis.wid + val.X) * 3
			tmpBMP[idx+0] = 0
			tmpBMP[idx+1] = 0
			tmpBMP[idx+2] = 0xff
		}
	}

	if src != nil {
		idx := (src.Y*pthis.wid + src.X) * 3
		tmpBMP[idx+0] = 0xff
		tmpBMP[idx+1] = 0
		tmpBMP[idx+2] = 0
	}
	if dst != nil {
		idx := (dst.Y*pthis.wid + dst.X) * 3
		tmpBMP[idx+0] = 0xff
		tmpBMP[idx+1] = 0
		tmpBMP[idx+2] = 0
	}

	bmp := CreateBMP(pthis.wid, pthis.hei, 24, tmpBMP)

	bmp.WriteBmp(strMapFile)
}
