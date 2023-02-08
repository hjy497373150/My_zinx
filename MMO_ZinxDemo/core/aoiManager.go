package core

import "fmt"

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