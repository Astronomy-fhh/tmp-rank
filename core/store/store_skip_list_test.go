package store

import (
	"testing"
	"tmp-rank/core"
)

// 测试排序功能
func TestAddSkipListSort(t *testing.T) {

	sls := NewSkipListStore[*core.PlayerScore]()

	sls.Add(&core.PlayerScore{Uid: 4, Score: 2, Ctime: 20, Level: 1, Name: "Player_2"})
	sls.Add(&core.PlayerScore{Uid: 3, Score: 3, Ctime: 20, Level: 2, Name: "Player_1"})
	sls.Add(&core.PlayerScore{Uid: 1, Score: 5, Ctime: 10, Level: 1, Name: "Player_1"})
	sls.Add(&core.PlayerScore{Uid: 2, Score: 4, Ctime: 20, Level: 1, Name: "Player_2"})
	sls.Add(&core.PlayerScore{Uid: 5, Score: 1, Ctime: 20, Level: 1, Name: "Player_3"})

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

	sls.skipList.Look()
	println()
	sls.skipList.LookTail()

	for _, testOne := range testSet {
		testResult, startRank := sls.GetSort(testOne.uid, testOne.num)
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
