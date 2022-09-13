package main

import "fmt"

type person struct {
	name string
}

func handle1(persons []*person) {
	for _, p := range persons {
		p.name += ".0.."
	}
}

func handle2(persons []person) {
	for _, p := range persons {
		p.name += ".0.."
	}
}

func dPrint(persons []*person) {
	for _, p := range persons {
		fmt.Printf("%+v\n", p)
	}
}

func dPrint1(persons []person) {
	for _, p := range persons {
		fmt.Printf("%+v\n", p)
	}
}

func main() {
	// 切片类型本身也是一个值类型，包括 map 也是一个值类型
	// 但为啥将切片传入函数可以得到更改？
	// 因为切片里面的元素是指针类型，因此拷贝切片就是拷贝了一堆指针
	// 而指针可以直接修改对应的元素值
	// 但如果切片的元素也是值类型，那么传入函数中的参数将是一个全新的切片
	// 切片、元素都是值拷贝，因此即使在元素中修改了，也不会影响到原来的元素值

	p1 := person{name: "pedro"}
	persons := []*person{&p1}
	dPrint(persons)
	handle1(persons)
	dPrint(persons)

	i := []person{{name: "pedro"}}
	dPrint1(i)
	handle2(i)
	dPrint1(i)
}
