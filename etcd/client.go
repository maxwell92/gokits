package etcdclient

import (
	"time"
	"encoding/json"
	mlog "github.com/maxwell92/gokits/log"
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

var log = mlog.Log

type Etcd interface {
	Get(key string) string
	Put(key, value string)
	Del(key string)
	EndpointHealth() string
	EndpointStatus() string
	MemberList() string
}

type EtcdClient struct {
	client *clientv3.Client
	ctx context.Context
	cancel context.CancelFunc
}

type EtcdClientConfig struct {
	*clientv3.Config
}

func NewEtcdClientConfig(Endpoints []string, DialTimeout time.Duration) *EtcdClientConfig {
	ecc := new(EtcdClientConfig)
	ecc.Config = new(clientv3.Config)
	ecc.Endpoints = Endpoints
	ecc.DialTimeout = DialTimeout
	return ecc
}

func NewEtcdClient(cfg *EtcdClientConfig) *EtcdClient {
	ec := new(EtcdClient)
	cli, err := clientv3.New(*cfg.Config)
	if err != nil {
		log.Fatalf("NewEtcdClient Error: error=%s", err)
		return nil
	}
	ec.client = cli
	ec.ctx, ec.cancel = context.WithTimeout(context.Background(), cfg.DialTimeout)
	return ec
}

func (ec *EtcdClient)Get(key string) string {
	resp, err := ec.client.Get(ec.ctx, key)
	if err != nil {
		log.Errorf("EtcdClient Get Key Error: error=%s, key=%s", err, key)
		return ""
	}

	for _, ev := range resp.Kvs {
		log.Infof("%s : %s", ev.Key, ev.Value)
	}

	ec.cancel()
	data, err := json.Marshal(resp.Kvs)
	if err != nil {
		log.Errorf("EtcdClient Get Key Marshal Error: error=%s, key=%s", err, key)
		return ""
	}

	return string(data)
}

func (ec *EtcdClient)Put(key, value string) {

}

func (ec *EtcdClient)Del(key string) {

}

func (ec *EtcdClient)EndpointHealth() string {
	return ""
}

func (ec *EtcdClient)EndpointStatus() string {
	return ""
}

func (ec *EtcdClient)MemberList() string {
	resp, err := ec.client.MemberList(context.Background())
	if err != nil {
		log.Errorf("EtcdClient MemberList Error: error=%s", err)
		return ""
	}

	data, err := json.Marshal(resp.Members)
	if err != nil {
		log.Errorf("EtcdClient MemberList Marshal Error: error=%s", err)
		return ""
	}

	return string(data)
}

