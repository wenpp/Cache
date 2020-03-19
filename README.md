# 简单的伪分布式缓存Demo
## 简单说明
### 架构图
了解架构请参看如下链接：https://www.processon.com/view/link/5e70ab41e4b092510f5fc672
### 项目结构
Cache-
   |-cache
     |-Cache.go                 核心数据结构，构建了一个LRU的缓存结构。通过字典表和双向链表实现
     |-CacheNameSpace.go        缓存命名空间，用于把缓存进行分类存储
     |-CacheNode.go             基于HashRing的缓存节点
     |-CacheView.go             构建了不可变的返回值的缓存数据
     |-Consts.go                常量
     |-HashRing.go              哈希环
   |-db                         模拟了一个简单的数据源
   |-test
     |-cache_test.go            核心数据测试
     |-cachenamespace_test.go   缓存控件测试
## RoadMap
- [x] 基本的缓存结构；
- [x] 一致性Hash；
- [x] 多节点服务；
- [ ] 把Http通信改造成基于protobuf的通信；
- [ ] 基于raft协议进行分布式改造；