package core

import "sync"

/*
	当前游戏的世界管理模块
*/
type WorldManager struct {
	// 当前世界地图的AOI管理模块
	AoiManager *AOIManager
	// 当前全部在线的playerId的集合
	Players map[int32] *Player
	// 保护player的锁
	PLock sync.RWMutex
}

// 提供一个对外的世界管理模块的句柄 单例
var WorldMgrObj *WorldManager

// 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		AoiManager: NewAOIManager(AOI_MIN_X,AOI_MAX_X,AOI_CNTSX,AOI_MIN_Y,AOI_MAX_Y,AOI_CNTSY) ,
		Players: make(map[int32]*Player),
	}
}

// 添加玩家
func (wm *WorldManager)AddPlayer(player *Player) {
	wm.PLock.Lock()
	wm.Players[player.PlayerId] = player
	wm.PLock.Unlock()

	// 将player加入到AOIManager中
	wm.AoiManager.AddPlayToGridByPos(int(player.PlayerId),player.X,player.Z)
}

// 删除玩家
func (wm *WorldManager)RemovePlayer(playerId int32) {

	// 得到当前玩家并从AOIManager中删除
	player := wm.Players[playerId]
	wm.AoiManager.RemovePlayFromGridByPos(int(playerId),player.X,player.Z)

	wm.PLock.Lock()
	delete(wm.Players,playerId)
	wm.PLock.Unlock()

}

// 通过玩家ID查询Player
func (wm *WorldManager)GetPlayerByPid(PlayerId int32) *Player {
	wm.PLock.RLock()
	defer wm.PLock.RUnlock()

	return wm.Players[PlayerId]
}

// 获取全部的在线玩家
func (wm *WorldManager)GetAllPlayers() []*Player {
	wm.PLock.RLock()
	defer wm.PLock.RUnlock()

	players := make([]*Player,0)
	for _,player := range wm.Players {
		players = append(players, player)
	}
	return players
}