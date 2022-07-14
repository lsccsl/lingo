package lin_common

import (
	"fmt"
	"math"
	"sort"
)

type JPS_DIR int
const (
	JPS_DIR_up    JPS_DIR = 0
	JPS_DIR_down  JPS_DIR = 1
	JPS_DIR_left  JPS_DIR = 2
	JPS_DIR_right JPS_DIR = 3

	JPS_DIR_up_left    JPS_DIR = 4
	JPS_DIR_up_right   JPS_DIR = 5
	JPS_DIR_down_left  JPS_DIR = 6
	JPS_DIR_down_right JPS_DIR = 7

	JPS_DIR_MAX  JPS_DIR = 8
)

var array_dirJPS = [JPS_DIR_MAX]Coord2d{
	{0,-1},
	{0,1},
	{-1,0},
	{1,0},
	{-1,-1},
	{1,-1},
	{-1,1},
	{1,1},
}

type JPSNode struct {
	parent *JPSNode
	subNode []*JPSNode

	totalWeight int

	pos Coord2d
	startWeight int
	endWeight int
	forceNeighbor []Coord2d
}
type MAP_JSP_HISTORY_PATH map[Coord2d]*JPSNode
type JSPMgr struct {
	pMapData *MapData
	root *JPSNode
	nodes []*JPSNode
	mapHistoryPath MAP_JSP_HISTORY_PATH

	src Coord2d
	dst Coord2d
}
func (pthis*JSPMgr)Len() int {
	return len(pthis.nodes)
}
func (pthis*JSPMgr)Less(i, j int) bool {
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
func (pthis*JSPMgr)Swap(i, j int) {
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
func (pthis*JSPMgr)addNode(node *JPSNode, parent *JPSNode) {
	if node == nil {
		return
	}

	if pthis.isInHistory(node.pos) {
		LogErr("already in history")
		return
	}

	node.totalWeight = node.startWeight + node.endWeight
	pthis.nodes = append(pthis.nodes, node)
	sort.Sort(pthis)
	pthis.addHistory(node)

	if node == pthis.root {
		return
	}

	if pthis.root == nil {
		pthis.root = node
	} else {
		if parent == nil {
			node.parent = pthis.root
		} else {
			node.parent = parent
		}
	}

	if node.parent != nil {
		node.parent.subNode = append(node.parent.subNode, node)
	}
}
func (pthis*JSPMgr)getNearestNode() *JPSNode {
	if len(pthis.nodes) == 0 {
		return nil
	}
	lastIdx := len(pthis.nodes) - 1
	node := pthis.nodes[lastIdx]
	pthis.nodes = pthis.nodes[:lastIdx]
	return node
}
func (pthis*JSPMgr)isInHistory(pos Coord2d) bool {
	_, ok := pthis.mapHistoryPath[pos]
	return ok
}
func (pthis*JSPMgr)addHistory(node *JPSNode) {
	pthis.mapHistoryPath[node.pos] = node
}

func getDirVector(dir JPS_DIR) *Coord2d{
	return &array_dirJPS[dir]
}

func (pthis*MapData)hasForceNeighbor(jpsMgr *JSPMgr, searchPos Coord2d, dir JPS_DIR) (bFindForceNeighbor bool, forceNeighbor []Coord2d) {
	bFindForceNeighbor = false
	forceNeighbor = nil

	if jpsMgr == nil {
		return
	}

	if jpsMgr.dst.isNear(&searchPos) {
		bFindForceNeighbor = true
		return
	}

	switch dir {
	case JPS_DIR_up:
		posLeft  := Coord2d{searchPos.X - 1, searchPos.Y}
		posRight := Coord2d{searchPos.X + 1, searchPos.Y}
		if pthis.CoordIsBlock(posLeft) {
			posN := Coord2d{posLeft.X, posLeft.Y - 1}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}
		if pthis.CoordIsBlock(posRight) {
			posN := Coord2d{posRight.X, posRight.Y - 1}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}

	case JPS_DIR_down:
		posLeft  := Coord2d{searchPos.X - 1, searchPos.Y}
		posRight := Coord2d{searchPos.X + 1, searchPos.Y}
		if pthis.CoordIsBlock(posLeft) {
			posN := Coord2d{posLeft.X, posLeft.Y + 1}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}
		if pthis.CoordIsBlock(posRight) {
			posN := Coord2d{posRight.X, posRight.Y + 1}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}

	case JPS_DIR_left:
		posUp   := Coord2d{searchPos.X, searchPos.Y - 1}
		posDown := Coord2d{searchPos.X, searchPos.Y + 1}
		if pthis.CoordIsBlock(posUp) {
			posN := Coord2d{posUp.X - 1, posUp.Y}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}
		if pthis.CoordIsBlock(posDown) {
			posN := Coord2d{posDown.X - 1, posDown.Y}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}

	case JPS_DIR_right:
		posUp   := Coord2d{searchPos.X, searchPos.Y - 1}
		posDown := Coord2d{searchPos.X, searchPos.Y + 1}

		if pthis.CoordIsBlock(posUp) {
			posN := Coord2d{posUp.X + 1, posUp.Y}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}
		if pthis.CoordIsBlock(posDown) {
			posN := Coord2d{posDown.X + 1, posDown.Y}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}

	case JPS_DIR_up_left:
		/*
           _ _ N
           _ @ *
		   N * \
		*/
		posRight := Coord2d{searchPos.X + 1, searchPos.Y}
		posDown  := Coord2d{searchPos.X, searchPos.Y + 1}
		if pthis.CoordIsBlock(posRight) {
			posN := Coord2d{posRight.X, posRight.Y - 1}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}
		if pthis.CoordIsBlock(posDown) {
			posN := Coord2d{posDown.X - 1, posDown.Y}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}

	case JPS_DIR_up_right:
		/*
            N _ _
            * @ _
		    / * N
		*/
		posLeft := Coord2d{searchPos.X - 1, searchPos.Y}
		posDown := Coord2d{searchPos.X, searchPos.Y + 1}
		if pthis.CoordIsBlock(posLeft) {
			posN := Coord2d{posLeft.X, posLeft.Y - 1}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}
		if pthis.CoordIsBlock(posDown) {
			posN := Coord2d{posDown.X + 1, posDown.Y}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}

	case JPS_DIR_down_left:
		/*
          N * /
          _ @ *
          _ _ N
		*/
		posUp    := Coord2d{searchPos.X, searchPos.Y - 1}
		posRight := Coord2d{searchPos.X + 1, searchPos.Y}
		if pthis.CoordIsBlock(posUp) {
			posN := Coord2d{posUp.X - 1, posUp.Y}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}
		if pthis.CoordIsBlock(posRight) {
			posN := Coord2d{posRight.X, posRight.Y + 1}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}

	case JPS_DIR_down_right:
		/*
            \ * N
            * @ _
            N _ _
		*/
		posUp   := Coord2d{searchPos.X, searchPos.Y - 1}
		posLeft := Coord2d{searchPos.X - 1, searchPos.Y}
		if pthis.CoordIsBlock(posUp) {
			posN := Coord2d{posUp.X + 1, posUp.Y}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}
		if pthis.CoordIsBlock(posLeft) {
			posN := Coord2d{posLeft.X, posLeft.Y + 1}
			if !pthis.CoordIsBlock(posN) {
				bFindForceNeighbor = true
				forceNeighbor = append(forceNeighbor, posN)
			}
		}
	}

	return
}

func (pthis*MapData)searchHorVer(jpsMgr *JSPMgr, curNode *JPSNode, curPos Coord2d, enDir JPS_DIR, bAdd bool) (findJump bool) {
	findJump = false
	dir := getDirVector(enDir)
	if dir == nil {
		return
	}
	searchPos := curPos
	curWeightAdd := 0
	for {
		if pthis.IsBlock(searchPos.X, searchPos.Y) {
			break
		}

		bFindForceNeighbor, forceNeighbor := pthis.hasForceNeighbor(jpsMgr, searchPos, enDir)
		if bFindForceNeighbor {
			findJump = true
			if bAdd && !jpsMgr.isInHistory(searchPos){
				jp := &JPSNode{
					pos:searchPos,
					endWeight:calEndWeight(searchPos, jpsMgr.dst),
					startWeight:curNode.startWeight + curWeightAdd,
					forceNeighbor:forceNeighbor,
				}

				// add pos to open list
				jpsMgr.addNode(jp, curNode)
				pthis.DumpJPSMap("../resource/jumppath.bmp", nil, jpsMgr)
			}
			break
		}
		searchPos = Coord2d{searchPos.X + dir.X, searchPos.Y + dir.Y}
		curWeightAdd += WEIGHT_straight
	}

	return
}

func (pthis*MapData)searchSlash(jpsMgr *JSPMgr, curNode *JPSNode, curPos Coord2d, enDir JPS_DIR) {

	bFindForceNeighbor, forceNeighbor := pthis.hasForceNeighbor(jpsMgr, curPos, enDir)
	if bFindForceNeighbor {
		curNode.forceNeighbor = append(curNode.forceNeighbor, forceNeighbor...)
	}

	dir := getDirVector(enDir)
	newPos := Coord2d{curPos.X, curPos.Y}
	curWeightAdd := 0
	for {
		if pthis.CoordIsBlock(newPos) {
			break
		}

		newPos = Coord2d{newPos.X + dir.X, newPos.Y + dir.Y}
		curWeightAdd += WEIGHT_slash

		findJump := false
		// 横向 纵向搜索
		switch enDir {
		case JPS_DIR_up_left:
			findJump = pthis.searchHorVer(jpsMgr, curNode, newPos, JPS_DIR_up, false)
			if !findJump {
				findJump = pthis.searchHorVer(jpsMgr, curNode, newPos, JPS_DIR_left, false)
			}

		case JPS_DIR_up_right:
			findJump = pthis.searchHorVer(jpsMgr, curNode, newPos, JPS_DIR_up, false)
			if !findJump {
				findJump = pthis.searchHorVer(jpsMgr, curNode, newPos, JPS_DIR_right, false)
			}

		case JPS_DIR_down_left:
			findJump = pthis.searchHorVer(jpsMgr, curNode, newPos, JPS_DIR_down, false)
			if !findJump {
				findJump = pthis.searchHorVer(jpsMgr, curNode, newPos, JPS_DIR_left, false)
			}

		case JPS_DIR_down_right:
			findJump = pthis.searchHorVer(jpsMgr, curNode, newPos, JPS_DIR_down, false)
			if !findJump {
				findJump = pthis.searchHorVer(jpsMgr, curNode, newPos, JPS_DIR_right, false)
			}
		}

		if findJump && !jpsMgr.isInHistory(newPos){
			jp := &JPSNode{
				pos:newPos,
				endWeight:calEndWeight(newPos, jpsMgr.dst),
				startWeight:curNode.startWeight + curWeightAdd,
				forceNeighbor:nil,
			}

			jpsMgr.addNode(jp, curNode)
			pthis.DumpJPSMap("../resource/jumppath.bmp", nil, jpsMgr)
		}
	}
}

func (pthis*MapData)PathJPS(src Coord2d, dst Coord2d) (path []Coord2d) {

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	/*
		(1)节点 A 是起点、终点.
		(2)节点A 至少有一个强迫邻居. 以横向x为例,说明此时被阻挡了,并且越过这个点之后,有新的连通区域
		(3)父节点在斜方向(斜向搜索)，节点A的水平或者垂直方向上有满足 (1)、(2) 的节点 用于转向

	(1).如果节点A没有父方向P(起点)
	则直线方向按照 (上下左右)四个方向， dirList = {上、下、左、右}
	斜方向按照(左上、右上、右下、左下)四个方向搜索 dirList = {左上、右上、右下、左下}
	(2).如果节点A有父方向P
	则 PA=（X,Y）
	将PA分解为水平 horizontalDir=(X,0)，垂直 verticalDir=(0,Y)

	还是先考虑水平和垂直的搜索方向,dirList = {}
	如果 horizontalDir=(X,0) != (0, 0) 即 X ！= 0 则 将 horizontalDir 加入到 dirList
	如果 verticalDir=(0,Y) !=(0, 0) 即 Y ！= 0 则 将 verticalDir 加入到 dirList
	直线方向搜索 dirList 中的方向

	然后是斜向 dirList = {}
	如果 PA=(X,Y)，X ！= 0 且 Y ！= 0， 则将 PA方向加入到 dirList
	如果 A有强迫邻居 {N1, N2, N3…},则将 AN1，AN2，AN3，。。。都加入到 dirList

	强近邻居是斜的,所以要分解

	————————————————
	版权声明：本文为CSDN博主「[奋斗不止]」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
	原文链接：https://blog.csdn.net/LIQIANGEASTSUN/article/details/118766080
	http://qiao.github.io/PathFinding.js/visual/
	*/

	startNode := &JPSNode{
		parent: nil,
		pos:src,
		endWeight:calEndWeight(src, dst),
	}
	startNode.totalWeight = startNode.endWeight

	jpsMgr := &JSPMgr{
		pMapData:pthis,
		root:startNode,
		mapHistoryPath: make(MAP_JSP_HISTORY_PATH),
		src:src,
		dst:dst,
	}
	jpsMgr.addNode(startNode, nil)

	node := startNode

	for {
		curNode := jpsMgr.getNearestNode()
		if curNode == nil {
			node = curNode
			pthis.DumpJPSMap("../resource/jumppath.bmp", nil, jpsMgr)
			break
		}
		curPos := curNode.pos

		if curPos.isNear(&dst) {
			pthis.DumpJPSMap("../resource/jumppath.bmp", nil, jpsMgr)
			break
		}

		var straightDir []JPS_DIR
		var slashDir []JPS_DIR

		if curNode.parent != nil {
			// 根据父节点相对位置,以及强迫邻居相对位置分解方向搜索
			var vecRelativeDir []Coord2d

			relativeParentDir := curNode.pos.dec(&curNode.parent.pos)
			vecRelativeDir = append(vecRelativeDir, relativeParentDir)

			for _, val := range curNode.forceNeighbor {
				relativeDir := val.dec(&curNode.pos)
				vecRelativeDir = append(vecRelativeDir, relativeDir)
			}

			for _, relativeDir := range vecRelativeDir {
				if relativeDir.X > 0 {
					straightDir = append(straightDir, JPS_DIR_right)
				} else if relativeDir.X > 0 {
					straightDir = append(straightDir, JPS_DIR_left)
				}
				if relativeDir.Y > 0 {
					straightDir = append(straightDir, JPS_DIR_down)
				} else if relativeDir.Y > 0 {
					straightDir = append(straightDir, JPS_DIR_up)
				}

				if relativeDir.X > 0 && relativeDir.Y > 0 {
					slashDir = append(slashDir, JPS_DIR_down_right)
				}
				if relativeDir.X > 0 && relativeDir.Y < 0 {
					slashDir = append(slashDir, JPS_DIR_up_right)
				}
				if relativeDir.X < 0 && relativeDir.Y > 0 {
					slashDir = append(slashDir, JPS_DIR_down_left)
				}
				if relativeDir.X < 0 && relativeDir.Y < 0 {
					slashDir = append(slashDir, JPS_DIR_up_left)
				}
			}
		} else {
			for i := JPS_DIR_up; i <= JPS_DIR_right; i ++ {
				straightDir = append(straightDir, i)
			}

			for  i := JPS_DIR_up_left; i <= JPS_DIR_down_right; i ++ {
				slashDir = append(slashDir, i)
			}
		}

		for _, val := range straightDir {
			pthis.searchHorVer(jpsMgr, curNode, curPos, val, true)
		}

		for _, val := range slashDir {
			pthis.searchSlash(jpsMgr, curNode, curPos, val)
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

func (pthis*MapData)DumpNodeSub(searchMgr *JSPMgr, node *JPSNode, bmp * Bitmap) {
	tmpBMP := bmp.BmpData
	widBytePitch := bmp.widBytePitch
	for _, val := range node.subNode {
		coordDiff := val.pos.dec(&node.pos)
		if coordDiff.X != 0 {
			coordDiff.X = coordDiff.X / int(math.Abs(float64(coordDiff.X)))
		}
		if coordDiff.Y != 0 {
			coordDiff.Y = coordDiff.Y / int(math.Abs(float64(coordDiff.Y)))
		}

		pos := node.pos.add(&coordDiff)
		for {
			idx := pos.Y * widBytePitch + pos.X * 3
			tmpBMP[idx + 0] = 0xff
			tmpBMP[idx + 1] = 0
			tmpBMP[idx + 2] = 0xff

			pos = pos.add(&coordDiff)
			if pos.isNear(&val.pos) {
				break
			}
		}

		pthis.DumpNodeSub(searchMgr, val, bmp)
	}
}

func (pthis*MapData)DumpJPSMap(strMapFile string, path []Coord2d, searchMgr *JSPMgr) {

	widBytePitch := CalWidBytePitch(pthis.widReal, 24)
	dataLen := widBytePitch * pthis.hei
	tmpBMP := make([]uint8, dataLen)
	bmp := CreateBMP(pthis.widReal, pthis.hei, 24, nil)
	bmp.BmpData = tmpBMP

	for j := 0; j < pthis.hei; j ++ {
		for i := 0; i < pthis.widReal; i ++ {
			var clr uint8 = 0xff
			if pthis.IsBlock(i, j) {
				clr = 0
			}
			newIdx := j * widBytePitch + i * 3
			tmpBMP[newIdx + 0] = clr
			tmpBMP[newIdx + 1] = clr
			tmpBMP[newIdx + 2] = clr
		}
	}

	pthis.DumpNodeSub(searchMgr, searchMgr.root, bmp)

	if searchMgr != nil {
		for key, node := range searchMgr.mapHistoryPath {
			idx := key.Y*widBytePitch + key.X * 3
			tmpBMP[idx+0] = 0
			tmpBMP[idx+1] = 0xff
			tmpBMP[idx+2] = 0

			for _, val := range node.forceNeighbor {
				idx := val.Y*widBytePitch + val.X * 3
				tmpBMP[idx+0] = 0
				tmpBMP[idx+1] = 0
				tmpBMP[idx+2] = 0xff
			}
		}
	}

	if path != nil {
		for _, val := range path {
			idx := val.Y*widBytePitch + val.X * 3
			tmpBMP[idx+0] = 0
			tmpBMP[idx+1] = 0
			tmpBMP[idx+2] = 0xff
		}
	}

	{
		idx := searchMgr.src.Y*widBytePitch + searchMgr.src.X * 3
		tmpBMP[idx+0] = 0xff
		tmpBMP[idx+1] = 0
		tmpBMP[idx+2] = 0
	}
	{
		idx := searchMgr.dst.Y*widBytePitch + searchMgr.dst.X * 3
		tmpBMP[idx+0] = 0xff
		tmpBMP[idx+1] = 0
		tmpBMP[idx+2] = 0
	}

	bmp.WriteBmp(strMapFile)
}