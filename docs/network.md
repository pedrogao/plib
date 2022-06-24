# 网络(network)

```sh
$ ip tuntap add dev tun101 mode tun
$ ip link set up dev tun101
$ ip a add 192.0.2.1/24 dev tun101

$ ping 192.0.2.2
```

## references

- [Using TUN/TAP in go or how to write VPN](https://nsl.cz/using-tun-tap-in-go-or-how-to-write-vpn/)
- [kaitai](https://formats.kaitai.io/)
- [ide](https://ide.kaitai.io/)
- https://www.cnblogs.com/bakari/p/10474600.html