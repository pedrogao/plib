# 探索 CPU cache 带来的性能差异


```sh
$ sysctl hw.l1dcachesize
hw.l1dcachesize: 32768

$ sysctl -a | grep cacheline
hw.cachelinesize: 64
```

 `32768 / 64 = 512`

## 参考资料

- https://teivah.medium.com/go-and-cpu-caches-af5d32cc5592
- https://itnext.io/understanding-the-lmax-disruptor-caaaa2721496
