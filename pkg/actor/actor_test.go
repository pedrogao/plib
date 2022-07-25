package actor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type radderAdd struct {
	value int
}

type radderGet struct {
	reply chan int
}

func Test_reducerActor(t *testing.T) {
	adder := NewFromReducer(0, func(raw any, sum int) int {
		switch msg := raw.(type) {
		case radderAdd:
			return msg.value + sum
		case radderGet:
			msg.reply <- sum
			return sum
		default:
			panic(fmt.Errorf("unsupported message %T", raw))
		}
	})

	adder.Send(radderAdd{value: 1}) // +1
	adder.Send(radderAdd{value: 2}) // +2
	adder.Send(radderAdd{value: 3}) // +3
	get := radderGet{reply: make(chan int, 1)}
	adder.Send(get) // get

	actual := <-get.reply

	assert.Equal(t, 6, actual)
}

type radderAdd1 struct {
	value int
}

func (msg *radderAdd1) Apply(st int) int {
	return st + msg.value
}

func add(value int) Message[int] {
	return &radderAdd1{value: value}
}

func identity[T any](val T) T { return val }

func Test_TypedActor(t *testing.T) {
	adder := NewTyped(0)
	adder.Send(add(1))
	adder.Send(add(2))
	adder.Send(add(3))

	result := Get(adder, identity[int]) // identity 映射函数
	assert.Equal(t, 6, result)
}

func TestAddMessage(t *testing.T) {
	assert.Equal(t, 0, add(0).Apply(0))
	assert.Equal(t, 11, add(1).Apply(10))
}
