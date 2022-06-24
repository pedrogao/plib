# 密码学基本算法实现

## 分类

- 古典密码学
- 近代密码学
- 现代密码学
    - 散列算法、数字摘要：MD5，sha1，sha256
    - 可读性算法：base64，base58
    - 对称加密：DES，AES
    - 非对称加密：RSA

## 参考资料

- https://hulining.github.io/2020/05/05/go-study-notes-package-crypto/
- 实现哈希函数md5：
    - https://rosettacode.org/wiki/MD5#
- 实现哈希函数sha1：
  - https://en.wikipedia.org/wiki/SHA-1
- 实现可读性算法 base64：
    - https://en.wikibooks.org/wiki/Algorithm_Implementation/Miscellaneous/Base64
    - http://www.sunshine2k.de/articles/coding/base64/understanding_base64.html
- 实现对称加密 DES：
    - https://www.educative.io/edpresso/how-to-implement-the-des-algorithm-in-cpp
    - https://ccm.net/contents/134-introduction-to-encryption-with-des
    - https://github.com/vmarlier/go-DES
    - https://www.cnblogs.com/idreamo/p/9333753.html
- 实现非对称加密 RSA：
    - https://www.ruanyifeng.com/blog/2013/06/rsa_algorithm_part_one.html
    - https://www.ruanyifeng.com/blog/2013/07/rsa_algorithm_part_two.html
    - https://eli.thegreenplace.net/2019/rsa-theory-and-implementation/
- 使用md5生成数字签名