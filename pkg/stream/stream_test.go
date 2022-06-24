package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream_Distinct(t *testing.T) {
	// 1 2 3 4 5
	Just(1, 2, 3, 3, 4, 5, 5).Distinct(func(item int) int {
		return item
	}).ForEach(func(item int) {
		t.Log(item)
	})

	// 1 2 3 4
	Just(1, 2, 3, 3, 4, 5, 5).Distinct(func(item int) int {
		// 对大于4的item进行特殊去重逻辑,最终只保留一个>3的item
		if item > 3 {
			return 4
		}
		return item
	}).ForEach(func(item int) {
		t.Log(item)
	})
}

func TestStream_Filter(t *testing.T) {
	// 保留偶数 2,4
	channel := Just(1, 2, 3, 4, 5).Filter(func(item int) bool {
		return item%2 == 0
	}).channel()
	for item := range channel {
		t.Log(item)
	}
}

func TestStream_Walk(t *testing.T) {
	// 返回 300,100,200
	Just(1, 2, 3).Walk(func(item int, pip chan<- int) {
		pip <- item * 100
	}, WithWorkers(3)).ForEach(func(item int) {
		t.Log(item)
	})
}

func TestStream_Head(t *testing.T) {
	// 返回 1,2
	channel := Just(1, 2, 3, 4, 5).Head(2).channel()
	for item := range channel {
		t.Log(item)
	}
}

func TestStream_Tail(t *testing.T) {
	// 4,5
	channel := Just(1, 2, 3, 4, 5).Tail(2).channel()
	for item := range channel {
		t.Log(item)
	}
	// 1,2,3,4,5
	channel2 := Just(1, 2, 3, 4, 5).Tail(6).channel()
	for item := range channel2 {
		t.Log(item)
	}
}

func TestStream_Map(t *testing.T) {
	// 返回 10,20,30,40,60,20....，并发执行的情况下是无序的
	channel := Just(1, 2, 3, 4, 5, 2, 2, 2, 2, 2, 2).
		Map(func(item int) int {
			return item * 10
		}, UnlimitedWorkers()).channel()
	for item := range channel {
		t.Log(item)
	}
}

func TestStream_Reverse(t *testing.T) {
	// 5,4,3,2,1
	channel := Just(1, 2, 3, 4, 5).Reverse().channel()
	for item := range channel {
		t.Log(item)
	}
}

func TestStream_Sort(t *testing.T) {
	// 5,4,3,2,1
	channel := Just(1, 2, 3, 4, 5).Sort(func(a, b int) bool {
		return a > b
	}).channel()
	for item := range channel {
		t.Log(item)
	}
}

func TestStream_Concat(t *testing.T) {
	// 10,9,6,5,4,3,2,1
	channel := Just(1, 2, 3, 4, 5).
		Concat(Just(6, 9, 10)).
		Sort(func(a, b int) bool {
			return a > b
		}).channel()
	for item := range channel {
		t.Log(item)
	}
}

func TestStream_AllMatch(t *testing.T) {
	assert := assert.New(t)
	// 10,9,6,5,4,3,2,1
	ok := Just(1, 2, 3, 4, 5).
		Concat(Just(6, 9, 10)).
		Sort(func(a, b int) bool {
			return a > b
		}).
		AllMatch(func(item int) bool {
			return item > 5
		})
	assert.False(ok)
}
