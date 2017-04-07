package kubewatcher

import (
	mcache "gokits/cache"
	mydeployment "app/backend/controller/yce/deploy"
	mytime "gokits/time"

	"fmt"
	"k8s.io/kubernetes/pkg/api"
	kwatch "k8s.io/kubernetes/pkg/watch"
	"strings"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	yceutils "app/backend/controller/yce/utils"
)

type ServiceWatcher struct {
	// Events <-chan kwatch.Event
	Events kwatch.Interface
	Cache  *mcache.RedisCache
	dcId   int32
	dcName string
	StopCh chan bool

	Fn map[kwatch.EventType]func(kwatch.Event)
}

func (sw *ServiceWatcher) Aggregate(eventCh chan kwatch.Event) {
	for {
		select {
		case e := <-sw.Events.ResultChan():
			{
				eventCh <- e
				log.Tracef("KubeWatcher ServiceWatcher receive event of type %s", e.Type)
			}
		case _ = <-sw.StopCh:
			{
				log.Tracef("KubeWatcher ServiceWatcher receive Stop")
				sw.Events.ResultChan()
				return
			}
		}
	}
}

func (sw *ServiceWatcher) Stop() {
	sw.Events.Stop()
	log.Tracef("KubeWatcher ServiceWatcher Stop")
}

func (sw *ServiceWatcher) Handle(e kwatch.Event) {
	log.Tracef("KubeWatcher ServiceWatcher handle event type %s, %p", e.Type, &e)
	sw.Fn[e.Type](e)
}

func (sw *ServiceWatcher) getDeploymentFromService(svcName string) string {
	frag := strings.Split(svcName, "-")
	length := len(frag)
	if length > 1 {
		log.Debugf("ServiceWatcher getDeploymentFromService get DeploymentName from %s by serviceName %s", strings.Join(frag[:length-1], "-"), svcName)
		return strings.Join(frag[:length-1], "-")
	} else {
		log.Debugf("ServiceWatcher getDeploymentFromService get DeploymentName from %s by serviceName %s", svcName, svcName)
		return svcName
	}
}

func (sw *ServiceWatcher) Setup(cache *mcache.RedisCache, dcId int32, dcName string, c *client.Client) {
	svc, err := c.Services(api.NamespaceAll).Watch(api.ListOptions{Watch: true, TimeoutSeconds: yceutils.Int64Ptr(mytime.DAYINSECONDS)})
	if err != nil {
		log.Errorf("Service Get Service Watcher Error: error=%s", err)
		return
	}

	sw.StopCh = make(chan bool)
	sw.Cache = cache
	sw.Events = svc
	sw.dcId = dcId
	sw.dcName = dcName
	sw.Fn = make(map[kwatch.EventType]func(kwatch.Event))
	sw.Fn[kwatch.Added] = sw.Add
	sw.Fn[kwatch.Modified] = sw.Modified
	sw.Fn[kwatch.Deleted] = sw.Delete
	sw.Fn[kwatch.Error] = sw.Error
}

func (sw *ServiceWatcher) Add(e kwatch.Event) {
	svc := e.Object.(*api.Service)
	log.Debugf("KubeWatcher ServiceWatcher %s Event", e.Type)
	dpNamespace := svc.Namespace
	dpName := sw.getDeploymentFromService(svc.Name)
	key := fmt.Sprintf("%d:%s:%s", sw.dcId, dpNamespace, dpName)
	cacheDp := sw.Cache.Get(key)
	log.Debugf("KubeWatcher ServiceWatcher Get Key %s value %s", key, cacheDp)
	if cacheDp == "" {
		log.Warnf("KubeWatcher ServiceWatcher get key %s empty", key)
		return
	}
	dp := new(mydeployment.Deployments)
	dp = dp.Decode(cacheDp)
	dp.Service = svc
	data := dp.Encode()
	ok, err := sw.Cache.Set(key, data)
	if err != nil || !ok {
		log.Errorf("KubeWatcher ServiceWatcher Add Error: error=%s", err)
		return
	}

	log.Debugf("KubeWatcher ServiceWatcher Service Add %s success", svc.Name)
}

func (sw *ServiceWatcher) Modified(e kwatch.Event) {
	svc := e.Object.(*api.Service)
	log.Debugf("KubeWatcher ServiceWatcher %s Event", e.Type)

	dpNamespace := svc.Namespace
	dpName := sw.getDeploymentFromService(svc.Name)
	key := fmt.Sprintf("%d:%s:%s", sw.dcId, dpNamespace, dpName)
	cacheDp := sw.Cache.Get(key)
	log.Debugf("KubeWatcher ServiceWatcher Get Key %s value %s", key, cacheDp)
	if cacheDp == "" {
		log.Warnf("KubeWatcher ServiceWatcher get key %s empty", key)
		return
	}
	dp := new(mydeployment.Deployments)
	dp = dp.Decode(cacheDp)
	dp.Service = svc
	data := dp.Encode()
	ok, err := sw.Cache.Set(key ,data)
	if err != nil || !ok {
		log.Errorf("KubeWatcher ServiceWatcher Update Error: error=%s", err)
		return
	}

	log.Debugf("KubeWatcher ServiceWatcher Update %s success", svc.Name)
}

func (sw *ServiceWatcher) Delete(e kwatch.Event) {
	svc := e.Object.(*api.Service)
	log.Debugf("KubeWatcher ServiceWatcher %s Event", e.Type)

	dpNamespace := svc.Namespace
	dpName := sw.getDeploymentFromService(svc.Name)
	key := fmt.Sprintf("%d:%s:%s", sw.dcId, dpNamespace, dpName)
	cacheDp := sw.Cache.Get(key)
	log.Debugf("KubeWatcher ServiceWatcher Get Key %s value %s", key, cacheDp)
	if cacheDp == "" {
		log.Warnf("KubeWatcher ServiceWatcher get key %s empty", key)
		return
	}
	dp := new(mydeployment.Deployments)
	dp = dp.Decode(cacheDp)
	dp.Service = nil
	data := dp.Encode()
	ok, err := sw.Cache.Set(key, data)
	if err != nil || !ok {
		log.Errorf("KubeWatcher ServiceWatcher Delete Error: error=%s", err)
		return
	}

	log.Debugf("KubeWatcher ServiceWatcher Delete %s success", svc.Name)
}

func (sw *ServiceWatcher) Error(e kwatch.Event) {
	log.Errorf("KubeWatcher ServiceWatcher Error")
}
