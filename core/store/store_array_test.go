package store

import (
	"testing"
	"tmp-rank/core"
)

// 测试排序功能
func TestArrayAddSort(t *testing.T) {
	as := NewArraySortStore[*core.PlayerScore]()

	as.Add(&core.PlayerScore{Uid: 1, Score: 100, Ctime: 10, Level: 1, Name: "Player_1"})
	as.Add(&core.PlayerScore{Uid: 2, Score: 200, Ctime: 20, Level: 1, Name: "Player_2"})
	as.Add(&core.PlayerScore{Uid: 3, Score: 100, Ctime: 20, Level: 2, Name: "Player_1"})
	as.Add(&core.PlayerScore{Uid: 4, Score: 100, Ctime: 20, Level: 1, Name: "Player_2"})
	as.Add(&core.PlayerScore{Uid: 5, Score: 100, Ctime: 20, Level: 1, Name: "Player_3"})

	as.Sort()

	expectedOrder := []int64{2, 1, 3, 4, 5}
	for i, uid := range expectedOrder {
		if as.sort[i].value.Uid != uid {
			t.Errorf("idx %d expected %d, got %d", i, uid, as.sort[i].value.Uid)
		}
	}
}

// 测试获取范围结果
func TestArrayGetSort(t *testing.T) {
	as := &ArraySortStore[*core.PlayerScore]{
		dict: make(map[int64]*ArraySortWrapper[*core.PlayerScore]),
	}

	as.Add(&core.PlayerScore{Uid: 1, Score: 5, Ctime: 10, Level: 1, Name: "Player_1"})
	as.Add(&core.PlayerScore{Uid: 2, Score: 4, Ctime: 20, Level: 1, Name: "Player_2"})
	as.Add(&core.PlayerScore{Uid: 3, Score: 3, Ctime: 20, Level: 2, Name: "Player_1"})
	as.Add(&core.PlayerScore{Uid: 4, Score: 2, Ctime: 20, Level: 1, Name: "Player_2"})
	as.Add(&core.PlayerScore{Uid: 5, Score: 1, Ctime: 20, Level: 1, Name: "Player_3"})

	as.Sort()

	type ExceptedRes struct {
		uid          int64
		num          int
		exceptedUid  []int64
		exceptedRank int64
	}

	testSet := []ExceptedRes{
		{3, 2, []int64{1, 2, 3, 4, 5}, 1},
		{4, 2, []int64{2, 3, 4, 5}, 2},
		{1, 2, []int64{1, 2, 3}, 1},
		{5, 2, []int64{3, 4, 5}, 3},
	}

	for _, testOne := range testSet {
		testResult, startRank := as.GetSort(testOne.uid, testOne.num)
		if len(testResult) != len(testOne.exceptedUid) {
			t.Errorf("getSort expected %d got %d", len(testOne.exceptedUid), len(testResult))
		}
		for i := 0; i < len(testResult); i++ {
			if testResult[i].Uid != testOne.exceptedUid[i] {
				t.Errorf("getSort Result index %d  expected %d got %d", i, testOne.exceptedUid[i], testResult[i].Uid)
			}
		}
		if startRank != testOne.exceptedRank {
			t.Errorf("getSort startRank expected %d got %d", testOne.exceptedRank, startRank)
		}
	}
}
