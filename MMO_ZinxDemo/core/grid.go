package core

import (
	"fmt"
	"sync"
)

/*
	一个AOI地图中的格子类型
*/
type Grid struct {
	// 格子ID
	Gid int
	// 格子的左边边界坐标
	MinX int
	// 格子的右边边界坐标
	MaxX int
	// 格子的上边界坐标
	MinY int
	// 格子的下边边界坐标
	MaxY int
	// 当前格子内玩家或物体的ID集合
	PlayerIds map[int] bool
	// 保护当前集合的锁
	GLock sync.RWMutex
}

// 初始化当前格子的方法
func NewGrid(gid,minx,maxx,miny,maxy int) *Grid {
	return &Grid{
		Gid: gid,
		MinX: minx,
		MaxX: maxx,
		MinY: miny,
		MaxY: maxy,
		PlayerIds: make(map[int] bool),
	}
}

// 向当前格子添加玩家
func (g *Grid)AddPlayer(playerId int) {
	// 使用写锁保护
	g.GLock.Lock()
	defer g.GLock.Unlock()
	
	fmt.Println("playerId = ",playerId,"is added to GridId = ",g.Gid)
	g.PlayerIds[playerId] = true
	
}

// 从当前格子中删除玩家
func (g *Grid)RemovePlayer(playerId int) {
	// 使用写锁保护
	g.GLock.Lock()
	defer g.GLock.Unlock()

	// 如果集合中有这个玩家，就删除他
	if _,ok := g.PlayerIds[playerId];ok {
		fmt.Println("playerId = ",playerId,"is removed from GridId = ",g.Gid)
		delete(g.PlayerIds, playerId)
	} else {
		fmt.Println("Not Found playerId = ",playerId,"in GridId = ",g.Gid)
	}
}

// 得到格子中的全部玩家
func (g *Grid) FindAllPlayer() (playerIds []int) {
	// 使用读锁保护
	g.GLock.RLock()
	defer g.GLock.RUnlock()

	for playerId, _ := range g.PlayerIds {
		playerIds = append(playerIds, playerId)
	}
	return 
}

//打印信息方法
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.Gid, g.MinX, g.MaxX, g.MinY, g.MaxY, g.PlayerIds)
}