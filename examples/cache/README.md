# 探索 CPU cache 带来的性能差异

```sh
$ sysctl hw.l1dcachesize
hw.l1dcachesize: 32768

$ sysctl -a | grep cacheline
hw.cachelinesize: 64
```

`32768 / 64 = 512`


处理器如何保持缓存一致性？ 如果两个缓存行共享一些公共地址，处理器会将它们标记为共享。
如果一个线程修改了共享行，它会将两者都标记为已修改。
为了保证缓存的一致性，它需要内核之间的协调，这可能会显着降低应用程序的性能。 这个问题被称为虚假分享。

## 参考资料

- https://teivah.medium.com/go-and-cpu-caches-af5d32cc5592
- https://medium.com/software-design/why-software-developers-should-care-about-cpu-caches-8da04355bb8a
- https://itnext.io/understanding-the-lmax-disruptor-caaaa2721496
