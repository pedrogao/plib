# Go 数组(array)、切片(slice)、映射(map)详解

## array

Go的数组是值类型，无论是赋值还是传参都是值拷贝。

## slice

```
ptr -> 内部数组
len
cap
```

slice 内部容器是数组，append或者切片操作都会自动处理数组的扩容。

切片拷贝使用 `copy` 函数。

## map

## 参考资料

- [数组、切片和映射](https://gfw.go101.org/article/container.html)
- [数组](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-array/)
- [切片](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-array-and-slice/)
- [哈希表](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-hashmap/)
- [深入解析 Go 中 Slice 底层实现](https://halfrost.com/go_slice/)