package kubewatcher

import (
	mylog "gokits/log"
	mcache "gokits/cache"
)

var log = mylog.Log


type SuperWatcher struct {
	KubeWatcher []*KubeWatcher
	Cache *mcache.RedisCache
}
