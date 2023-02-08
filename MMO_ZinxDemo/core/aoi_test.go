package core

import (
	"fmt"
	"testing"
)

func TestNewAoiManager(t *testing.T) {
	// 初始化AOIManager
	AOIManager := NewAOIManager(0,250,5,0,250,5)

	fmt.Println(AOIManager)
}

func TestAOIManagerSurroundGridsByGid(t *testing.T) {
	// 初始化AOIManager
	AOIManager := NewAOIManager(0,250,5,0,250,5)

	for gid,_ := range AOIManager.grids {
		// 得到当前gid的周边九宫格消息
		grids := AOIManager.GetSurroundGridsByGid(gid)
		fmt.Println("gid = ",gid,"grids len = ",len(grids))
		gids := make([]int,0,len(grids))
		for _,grid := range grids {
			gids = append(gids, grid.Gid)
		}
		fmt.Println("surrounding grid Ids are:",gids)
	}
}