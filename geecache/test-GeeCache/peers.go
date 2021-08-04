package test_GeeCache

import pb "Golang/Learning/geecache/test-GeeCache/geecachepb"

// 之前实现了 geeCache 的取得缓存值的流程 1 和流程 3，此处实现流程 2
/*
使用一致性哈希选择节点          是                                    是
      ｜-----> 是否是远程节点 -----> HTTP 客户端访问远程节点 --> 成功？-----> 服务端返回返回值
      ｜ 否                                               ↓  否
      ｜-----------------------------------------> 回退到本地节点处理。
*/

// PeerPicker 根据传入的 key 选择相应节点 PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 的 Get() 方法用于从对应 group 查找缓存值。PeerGetter 对应上述流程中的 HTTP 客户端
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
