package etcdclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"encoding/json"
	"fmt"
)

type MyStruct struct {
		Key []byte `json:"key"`
		Value []byte `json:"value"`
}

type MyStructList struct {
	List []*MyStruct `json:"list"`
}

func TestGet(t *testing.T) {
	cfg := NewEtcdClientConfig([]string{"10.151.32.27:30156"}, time.Second)
	cli := NewEtcdClient(cfg)
	dataSet := []string{
		// Existed
		"magic/mushroom",
	}
	data := cli.Get(dataSet[0])
	fmt.Println(data)

	// obj := new(MyStructList)
	obj := make([]MyStruct, 10)
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		assert.Error(t, err, "Error")
		log.Errorf("error=%s", err)
	}

	// fmt.Println(string(obj[0].Value))
	fmt.Printf("%v\n", string(obj[0].Value))

	// assert.Equal(t, "maxwell",  data,"Should Be Equal")
	// assert.Equal(t, "maxwell",  string(obj[0].Value),"Should Be Equal")
}

func TestMemberList(t *testing.T) {
	cfg := NewEtcdClientConfig([]string{"10.151.32.27:30156"}, time.Second)
	cli := NewEtcdClient(cfg)
	data := cli.MemberList()
	assert.Equal(t, "--",  data,"Should Be Equal")
}