package kubewatcher

import (
	client "k8s.io/kubernetes/pkg/client/unversioned"
	kwatch "k8s.io/kubernetes/pkg/watch"
	yceutils "app/backend/controller/yce/utils"
	"reflect"
	mcache	"gitlab.com/gokits/cache"
)

type Watcher interface {
	Setup(*mcache.RedisCache, int32, string, *client.Client)
	// Dispatch()
	Aggregate(chan kwatch.Event)
	Stop()
	Handle(kwatch.Event)
	Add(kwatch.Event)
	Modified(kwatch.Event)
	Delete(kwatch.Event)
	Error(kwatch.Event)
}

type KubeWatcher struct {
	ApiServer   string
	Client      *client.Client
	Deployments Watcher
	Pods        Watcher
	Services    Watcher
	EventCh     chan kwatch.Event
	StopCh      chan bool
	// ResetWatcherCh chan bool
}

func NewKubeWatcher(apiServer string) *KubeWatcher {
	kw := new(KubeWatcher)
	kw.ApiServer = apiServer
	c, err := yceutils.CreateK8sClient(apiServer)
	if err != nil {
		log.Errorf("KubeWatcher CreateK8sClient Error")
		return nil
	}
	kw.Client = c
	kw.Deployments = new(DeploymentWatcher)
	kw.Pods = new(PodWatcher)
	kw.Services = new(ServiceWatcher)

	kw.EventCh = make(chan kwatch.Event)
	kw.StopCh = make(chan bool)

	return kw
}

func (kw *KubeWatcher) Run(cache *mcache.RedisCache) {
	dcId := yceutils.GetDcIdByApiServer(kw.ApiServer)
	dc, _ := yceutils.GetDatacenterByDcId(dcId)
	dcName := dc.Name

	kw.Deployments.Setup(cache, dcId, dcName, kw.Client)
	kw.Pods.Setup(cache, dcId, dcName, kw.Client)
	kw.Services.Setup(cache, dcId, dcName, kw.Client)

	go kw.Deployments.Aggregate(kw.EventCh)
	go kw.Pods.Aggregate(kw.EventCh)
	go kw.Services.Aggregate(kw.EventCh)

	log.Infof("KubeWatcher %s Running", kw.ApiServer)
	kw.Handle()
}

func (kw *KubeWatcher) Handle() {
	for {
		select {
		case e := <- kw.EventCh : {
			kw.Dispatch(e)
			log.Tracef("KubeWatcher Handle receive event of type %s", e.Type)
		}
		case _ = <- kw.StopCh : {
			log.Tracef("KubeWatcher Handle receive Stop event")
			kw.Deployments.Stop()
			kw.Services.Stop()
			kw.Pods.Stop()
			return
		}
		}
	}
}

func (kw *KubeWatcher) Dispatch(e kwatch.Event) {
	if e.Object == nil {
		log.Errorf("KubeWatcher Dispatch Nil Event Object type %s", e.Type)
		return
	}
	if e.Type == "" {
		log.Errorf("KubeWatcher Dispatch Empty Event Type Object %s", reflect.TypeOf(e.Object).String())
		return
	}
	log.Tracef("%s:%s", e.Type, reflect.TypeOf(e.Object).String())
	switch reflect.TypeOf(e.Object).String() {
	case "*extensions.Deployment": {
		kw.Deployments.Handle(e)
		log.Debugf("KubeWatcher event Deployment")
	}
	case "*api.Pod": {
		kw.Pods.Handle(e)
		log.Debugf("KubeWatcher event Pod")
	}
	case "*api.Service": {
		kw.Services.Handle(e)
		log.Debugf("KubeWatcher event Service")
	}
	}
}

//为了解决重启yce后，因为Redis缓存里的值仍存在导致添加新值失败的问题（甚至在重启过程中，以其他方式比如kubectl触发的事件无法被捕捉而导致的信息不一致）
func (kw *KubeWatcher) Warm() {

}
/*
func (kw *KubeWatcher) ResetWatcher() {
	for {
		time.Sleep(time.Duration(1500) * time.Second)
		kw.ResetWatcherCh <- true
	}
}
*/
