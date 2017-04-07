package kubewatcher

import (
	"fmt"

	mydeployment "app/backend/controller/yce/deploy"
	mcache "github.com/maxwell92/gokits/cache"
	mytime "github.com/maxwell92/gokits/time"

	yceutils "app/backend/controller/yce/utils"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	kwatch "k8s.io/kubernetes/pkg/watch"
)

type DeploymentWatcher struct {
	// Events <-chan kwatch.Event
	Events kwatch.Interface
	StopCh chan bool
	Cache  *mcache.RedisCache
	dcId   int32
	dcName string
	client *client.Client

	Fn map[kwatch.EventType]func(kwatch.Event)
}

func (dw *DeploymentWatcher) Aggregate(eventCh chan kwatch.Event) {
	for {
		select {
		case e := <-dw.Events.ResultChan():
			{
				eventCh <- e
				log.Tracef("KubeWatcher DeploymentWatcher receive event of type %s", e.Type)
			}
		case _ = <-dw.StopCh:
			{
				log.Tracef("KubeWatcher DeploymentWatcher receive Stop")
				dw.Events.Stop()
				return
			}
		}
	}
}

func (dw *DeploymentWatcher) Stop() {
	dw.Events.Stop()
	log.Infof("KubeWatcher DeploymentWatcher Stop")
}

func (dw *DeploymentWatcher) Handle(e kwatch.Event) {
	log.Debugf("KubeWatcher DeploymentWatcher handle event type %s, %p", e.Type, &e)
	dw.Fn[e.Type](e)
}

func (dw *DeploymentWatcher) Setup(cache *mcache.RedisCache, dcId int32, dcName string, c *client.Client) {
	dp, err := c.Deployments(api.NamespaceAll).Watch(api.ListOptions{Watch: true, TimeoutSeconds: yceutils.Int64Ptr(mytime.DAYINSECONDS)})
	if err != nil {
		log.Errorf("KubeWatcher Get Deployment Watcher Error: error=%s", err)
		return
	}
	dw.Cache = cache
	dw.StopCh = make(chan bool)
	dw.Events = dp
	dw.dcId = dcId
	dw.dcName = dcName
	dw.Fn = make(map[kwatch.EventType]func(kwatch.Event))
	dw.Fn[kwatch.Added] = dw.Add
	dw.Fn[kwatch.Modified] = dw.Modified
	dw.Fn[kwatch.Deleted] = dw.Delete
	dw.Fn[kwatch.Error] = dw.Error
	dw.client = c
}

func (dw *DeploymentWatcher) Add(e kwatch.Event) {
	dp := e.Object.(*extensions.Deployment)
	log.Debugf("KubeWatcher DeploymentWatcher %s Event", e.Type)

	// 为了解决PodAdd事件或ServiceAdd事件先到找不到Deployment而无法加入的问题
	// 所以在每个deployment首次加入时，通过yceutils对podList和Service进行初始化
	podList, _ := yceutils.GetPodListByDeployment(dw.client, dp)
	svc, _ := yceutils.GetServiceByDeployment(dw.client, dp)

	newDp := &mydeployment.Deployments{
		UserName:       dp.Annotations["lastModified"],
		DcName:         dw.dcName,
		DcId:           dw.dcId,
		OrgName:        dp.Namespace,
		DeploymentName: dp.Name,
		UpdateTime:     dp.ObjectMeta.CreationTimestamp.String(),
		Deployment:     dp,
		PodList:        podList,
		Service:        svc,
	}

	key := fmt.Sprintf("%d:%s:%s", dw.dcId, newDp.OrgName, newDp.DeploymentName)
	data := newDp.Encode()
	ok, err := dw.Cache.Set(key, data)
	if !ok {
		log.Errorf("KubeWatcher DeploymentWatcher Add Error: error=%s", err)
		return
	}
	log.Debugf("KubeWatcher DeploymentWatcher Set key %s success", key)

	dpNameKey := fmt.Sprintf("%d:%s", newDp.DcId, newDp.OrgName)
	log.Debugf("KubeWatcher DeploymentWatcher dpNameKey %s", dpNameKey)

	length, err := dw.Cache.Scard(dpNameKey)
	if err != nil {
		log.Errorf("KubeWatcher DeploymentWatcher Scard Key %s Error: error=%s", dpNameKey, err)
	}
	log.Debugf("KubeWatcher DeploymentWatcher Scard Key %d results", length)

	dpNameList, err := dw.Cache.Smember(dpNameKey)
	if err != nil {
		log.Errorf("KubeWatcher DeploymentWatcher Smember Key %s Error: error=%s", dpNameKey, err)
	}
	log.Debugf("KubeWatcher DeploymentWatcher Smember Key %s results %v", dpNameKey, dpNameList)

	ok, err = dw.Cache.Sadd(dpNameKey, newDp.DeploymentName)
	if err != nil {
		log.Errorf("KubeWatcher DeploymentWatcher Sadd key %s with %s Error: error=%s", dpNameKey, err, newDp.DeploymentName)
		return
	}

	if !ok {
		log.Warnf("KubeWatcher DeploymentWatcher Sadd key %s with %s Error", dpNameKey, newDp.DeploymentName)
		return
	}

	log.Debugf("KubeWatcher DeploymentWatcher Sadd key %s success", dpNameKey)

	log.Infof("KubeWatcher DeploymentWatcher Smember Deployment of %s total %v", dp.Namespace, dpNameList)
}

func (dw *DeploymentWatcher) Modified(e kwatch.Event) {
	dp := e.Object.(*extensions.Deployment)
	log.Tracef("KubeWatcher DeploymentWatcher %s Event", e.Type)

	key := fmt.Sprintf("%d:%s:%s", dw.dcId, dp.Namespace, dp.Name)
	cacheDp := dw.Cache.Get(key)
	log.Debugf("KubeWatcher DeploymentWatcher Get Key %s value %s", key, cacheDp)
	oldDp := new(mydeployment.Deployments)
	oldDp = oldDp.Decode(cacheDp)
	log.Debugf("KubeWatcher DeploymentWatcher Decoded with %d pod", len(oldDp.PodList.Items))

	oldDp.Deployment = dp

	data := oldDp.Encode()
	log.Debugf("KubeWatcher DeploymentWatcher Encoded with %d pod", len(oldDp.PodList.Items))
	ok, err := dw.Cache.Set(key, data)
	log.Debugf("KubeWatcher DeploymentWatcher Update key %s deployment %s ", key, oldDp.DeploymentName)
	if err != nil {
		log.Errorf("KubeWatcher DeploymentWatcher Modified Error: error=%s", err)
		return
	}

	if !ok {
		log.Warnf("KubeWatcher DeploymentWatcher Modified Error")
		return
	}

	log.Infof("KubeWatcher DeploymentWatcher Modified Deployment Success")
}

func (dw *DeploymentWatcher) Delete(e kwatch.Event) {
	dp := e.Object.(*extensions.Deployment)
	log.Debugf("KubeWatcher DeploymentWatcher %s Event", e.Type)
	key := fmt.Sprintf("%d:%s:%s", dw.dcId, dp.Namespace, dp.Name)
	dpNameKey := fmt.Sprintf("%d:%s", dw.dcId, dp.Namespace)
	ok, err := dw.Cache.Delete(key)
	if err != nil || !ok {
		log.Errorf("KubeWatcher DeploymentWatcher Delete deployment %s type %s Error: error=%s", dp.Name, e.Type, err)
		return
	}

	ok, err = dw.Cache.Srem(dpNameKey, dp.Name)
	if err != nil || !ok {
		log.Errorf("KubeWatcher DeploymentWatcher Delete Error: error=%s")
		return
	}
	log.Infof("KubeWatcher DeploymentWatcher Delete deployment %s type %s success", dp.Name, e.Type)
}

func (dw *DeploymentWatcher) Error(e kwatch.Event) {
	log.Errorf("KubeWatcher DeploymentWatcher Error")
}
