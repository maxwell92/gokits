package kubewatcher

import (
	mcache "github.com/maxwell92/gokits/cache"
	"time"
	mytime "github.com/maxwell92/gokits/time"
)

func NewSuperWatcher(apiServers []string) *SuperWatcher {
	sw := &SuperWatcher{
		KubeWatcher: make([]*KubeWatcher, 0),
		Cache: mcache.RedisCacheInstance(),
	}

	for _, apiServer := range apiServers {
		kw := NewKubeWatcher(apiServer)
		sw.KubeWatcher = append(sw.KubeWatcher, kw)
	}

	return sw
}



func (sw *SuperWatcher) Run() {
	for _, kw := range sw.KubeWatcher {
		go kw.Run(sw.Cache)
	}

	sw.ResetWatcher()
}

func (sw *SuperWatcher) ResetWatcher() {
	 for {
		log.Tracef("SuperWatcher ResetWatcher Started")
		time.Sleep(time.Duration(mytime.DAYINSECONDS) * time.Second)
		log.Infof("SuperWatcher ResetWatcher Begin Reset")
		for _, kw := range sw.KubeWatcher {
			log.Tracef("SuperWatcher KubeWatcher will be reset, watcher %p ,stopCh %p", kw, kw.StopCh)
			kw.StopCh <- true
			go kw.Run(sw.Cache)
		}
		 log.Infof("SuperWatcher ResetWatcher Reset Success")
	 }
}
