# 分布式缓存

设计一个分布式缓存系统，需要考虑资源控制、淘汰策略、并发、分布式节点通信等各个方面的问题。而且针对不同的应用场景，还需要在不同的特性之间权衡。

`GeeCache` 基本上模仿了 `groupcache` 的实现。支持特性有：

- 单级缓存和基于 `HTTP` 的分布式缓存

- `LRU` 缓存策略

- 使用 `Go` 锁机制防止缓存击穿

- 使用一致性哈希选择节点，实现简单负载均衡

- 使用 `protobuf` 优化节点间二进制通信

代码结构的雏形

```
geecache/
    |--lru/
        |--lru.go             // lru 缓存淘汰策略
    |--byteview.go            // 缓存值的抽象与封装
    |--cache.go               // 并发控制
    |--geecache.go            // 负责与外部交互，控制缓存存储和获取的主流程
    |--http.go                // 提供被其他节点访问的能力（基于 http）
    |--peers.go               // 节点选择的抽象
    |--consistenthash/
        |--consistenthash.go  // 一致性哈希算法
    |--geecachepb/
        |--geecachepb.proto   // protobuf 的使用
        |--gecachepb.pb.go    // .proto 文件转换为 Go 代码
```