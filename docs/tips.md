# 小技巧

运行 1.18beta：

```sh
$ export GOROOT=~/go/go1.18beta1
$ ~/go/go1.18beta1/bin/go test ./pkg/collection
```

一键替换 interface{} 为 any：

```sh
gofmt -w -r 'interface{} -> any' ./
```