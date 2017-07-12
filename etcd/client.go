package etcdclient

import (
	"fmt"
	"time"	

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

type Etcd interface {
	Get(key string) string
	Put(key, value string)

}

type EtcdClient struct {

}


