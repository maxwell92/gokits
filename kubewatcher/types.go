package kubewatcher

import (
	mylog "github.com/maxwell92/gokits/log"
	mcache "github.com/maxwell92/gokits/cache"
)

var log = mylog.Log


type SuperWatcher struct {
	KubeWatcher []*KubeWatcher
	Cache *mcache.RedisCache
}
