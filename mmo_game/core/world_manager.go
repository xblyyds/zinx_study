package core

import "sync"

/*
	当前游戏世界的总管理模块
*/
type WorldManager struct {
	AoiMgr  *AOIManager       // 当前世界地图的AOI规划管理器
	Players map[int32]*Player // 当前在线玩家的集合
	pLock   sync.RWMutex      // 保护Players的互斥读写锁
}

// 提供一个对外的世界管理模块句柄
var WorldMgrObj *WorldManager

// 提供WorldManager 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		Players: make(map[int32]*Player),
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
	}
}

// 提供添加一个玩家的功能, 将玩家添加进玩家信息表Players
func (wm *WorldManager) AddPlayer(player *Player) {
	// 将player添加到 世界管理器中
	wm.pLock.Lock()
	wm.Players[player.PID] = player
	wm.pLock.Unlock()

	// 将player添加到AOI网络规划中
	wm.AoiMgr.AddPIDByPos(int(player.PID), player.X, player.Z)
}

// 从玩家信息表中移出一个玩家
func (wm *WorldManager) RemovePlayerByPid(pID int32) {
	/* 感觉这部分应该不用, 因为玩家下线后, 再次上线应该还是相同位置
	// 得到当前玩家
	player := wm.Players[pID]
	// 将玩家从AOIManager中删除
	wm.AoiMgr.RemoveFromGrIDByPos(int(pID), player.X, player.Z)
	*/
	// 将玩家从世界管理中删除
	wm.pLock.Lock()
	delete(wm.Players, pID)
	wm.pLock.Unlock()
}

// 提供玩家ID 获取对应玩家信息
func (wm *WorldManager) GetPlayerByPID(pID int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pID]
}

// 获取所有玩家信息
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	// 创建返回的player集合切片
	players := make([]*Player, 0)

	// 添加切片
	for _, player := range wm.Players {
		players = append(players, player)
	}

	// 返回
	return players
}
