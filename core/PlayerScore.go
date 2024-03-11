package core

type PlayerScore struct {
	Uid   int64
	Score int
	Ctime int
	Level int
	Name  string
}

// Compare 比较两个玩家分数 入参可传不同类型 多种策略
func (ps *PlayerScore) Compare(it interface{}) int {
	b := ps
	a, ok := it.(*PlayerScore)
	if !ok {
		return 0
	}
	if a.Score > b.Score {
		return 1
	}
	if a.Score < b.Score {
		return -1
	}
	if a.Ctime < b.Ctime {
		return 1
	}
	if a.Ctime > b.Ctime {
		return -1
	}
	if a.Level > b.Level {
		return 1
	}
	if a.Level < b.Level {
		return -1
	}
	if a.Name < b.Name {
		return 1
	}
	if a.Name > b.Name {
		return -1
	}
	return 0
}

func (ps *PlayerScore) GetKey() int64 {
	return ps.Uid
}
