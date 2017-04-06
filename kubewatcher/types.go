package kubewatcher

import (
	mylog "app/backend/common/util/log"
	mcache "app/backend/common/cache"
)

var log = mylog.Log


type SuperWatcher struct {
	KubeWatcher []*KubeWatcher
	Cache *mcache.RedisCache
}
