package kubewatcher

import (
	mylog "gitlab.com/gokits/log"
	mcache "gitlab.com/gokits/cache"
)

var log = mylog.Log


type SuperWatcher struct {
	KubeWatcher []*KubeWatcher
	Cache *mcache.RedisCache
}
