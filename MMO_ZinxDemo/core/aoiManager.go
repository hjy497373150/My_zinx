package core

import "fmt"


const (
	AOI_MIN_X int = 85
	AOI_MAX_X int = 410
	AOI_CNTSX int = 10
	AOI_MIN_Y int = 75
	AOI_MAX_Y int = 400
	AOI_CNTSY int = 20
)
/*
	AOI区域管理模块
	AOI(Area Of Interest):MMORPG游戏的核心算法之一，是服务器同步给玩家的信息区域
*/

type AOIManager struct {
	MinX  int           //区域左边界坐标
	MaxX  int           //区域右边界坐标
	CntsX int           //x方向格子的数量
	MinY  int           //区域上边界坐标
	MaxY  int           //区域下边界坐标
	CntsY int           //y方向的格子数量
	grids map[int] *Grid //当前区域中都有哪些格子，key=格子ID， value=格子对象
}

// 初始化一个AOI区域
func NewAOIManager(minX,maxX,cntsX,minY,maxY,cntsY int) *AOIManager {
	aoiManager := &AOIManager{
		MinX: minX,
		MaxX: maxX,
		CntsX: cntsX,
		MinY: minY,
		MaxY: maxY,
		CntsY: cntsY,
		grids: make(map[int] *Grid),
	}

	// 给AOI初始化区域的格子进行编号和初始化
	for x:=0;x<cntsX;x++ {
		for y:=0;y<cntsY;y++ {
			// 1.给格子进行编号
			gid := y*cntsX + x

			// 初始化格子放在grids里
			aoiManager.grids[gid] = NewGrid(gid,
				aoiManager.MinX + x * aoiManager.gridWidth(), 
				aoiManager.MinX + (x+1) * aoiManager.gridWidth(), 
				aoiManager.MinY + y * aoiManager.gridHeight(), 
				aoiManager.MinY + (y+1) * aoiManager.gridHeight())
		}
	}

	return aoiManager
}

// 得到每个格子在AOI区域X方向的宽度
func (aoiManager *AOIManager) gridWidth() int {
	return (aoiManager.MaxX - aoiManager.MinX) / aoiManager.CntsX
}

// 得到每个格子在AOI区域Y方向的高度
func (aoiManager *AOIManager) gridHeight() int {
	return (aoiManager.MaxY - aoiManager.MinY) / aoiManager.CntsY
}

//打印信息方法
func (aoiManager *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n Grids in AOI Manager:\n",
	aoiManager.MinX, aoiManager.MaxX, aoiManager.CntsX, aoiManager.MinY, aoiManager.MaxY, aoiManager.CntsY)

	for _,grid := range aoiManager.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

// 根据格子的gid得到周边九宫格格子的id集合
func (aoiManager *AOIManager) GetSurroundGridsByGid(gid int) (grids []*Grid) {
	// 1.判断gid是否在AOIManager中
	if _,ok := aoiManager.grids[gid];!ok {
		return 
	}
	// 2.将gid对应的格子加入集合
	grids = append(grids, aoiManager.grids[gid])

	// 3.根据gid得到当前格子所在X轴编号
	idx := gid % aoiManager.CntsX

	// 4.判断idx左右还有无格子
	if idx > 0 {
		grids = append(grids, aoiManager.grids[gid-1])
	}
	if idx < aoiManager.CntsX-1 {
		grids = append(grids, aoiManager.grids[gid+1])
	}

	// 5.将X轴上的所有格子取出，然后分别判断它们上下有没有格子
	gidsX := make([]int,0,len(grids))

	for _,grid := range grids {
		gidsX = append(gidsX, grid.Gid)
	}

	for _,v := range gidsX {
		// 计算该格子处于第几列
		idy := v / aoiManager.CntsX
		
		if idy > 0 {
			grids = append(grids, aoiManager.grids[v-aoiManager.CntsX])
		}
		if idy < aoiManager.CntsY - 1{
			grids = append(grids, aoiManager.grids[v+aoiManager.CntsX])
		}
	}

	return

}

// 通过横纵坐标得到格子编号
func (aoiManager *AOIManager) GetGidByPos(x,y float32) int {
	idx := (int(x) - aoiManager.MinX) / aoiManager.gridWidth()
	idy := (int(y) - aoiManager.MinY) / aoiManager.gridHeight()

	return idy * aoiManager.CntsX + idx
}

// 通过横纵坐标得到周边九宫格内全部的playerids
func (aoiManager *AOIManager) GetSurroundGridsByPos(x,y float32) (playerIds []int) {
	// 1.得到当前玩家的gid
	gid := aoiManager.GetGidByPos(x,y)
	// 2.通过gid得到周边九宫格信息
	grids := aoiManager.GetSurroundGridsByGid(gid)

	for _, grid := range grids {
		playerIds = append(playerIds, grid.FindAllPlayer()...)
		fmt.Printf("-----> gridId = %d,playerIds = %v -------",grid.Gid,grid.PlayerIds)
	} 
	return 
}

// 添加一个Player到一个格子中
func (aoiManager *AOIManager) AddPlayerToGridById(playerId,gid int) {
	aoiManager.grids[gid].AddPlayer(playerId)
}

// 移除一个格子中的某一个Player
func (aoiManager *AOIManager) RemovePlayerFromGrid(playerId,gid int) {
	aoiManager.grids[gid].RemovePlayer(playerId)
}

// 通过Gid获取全部的PlayerId
func (aoiManager *AOIManager) GetAllPlayerByGid(gid int) (playerIds []int) {
	playerIds = aoiManager.grids[gid].FindAllPlayer()

	return 
}

// 通过坐标将Player添加到一个格子中
func (aoiManager *AOIManager) AddPlayToGridByPos(playerId int,x,y float32) {
	gid := aoiManager.GetGidByPos(x,y)
	aoiManager.grids[gid].AddPlayer(playerId)
}

// 通过坐标把一个player从一个格子中删除
func (aoiManager *AOIManager) RemovePlayFromGridByPos(playerId int,x,y float32) {
	gid := aoiManager.GetGidByPos(x,y)
	aoiManager.grids[gid].RemovePlayer(playerId)
}