package util

type Item struct {
	ID    int
	Count int
}

type User struct {
	Account    string //账号
	ComStorage []Item //普通仓库物品
	VipStorage []Item //vip仓库物品
}

const (
	NONE = iota
	COMMON
	VIP
)

const (
	_ = iota
	TinsertCoreStat
)
