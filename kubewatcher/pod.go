package kubewatcher

import (
	"fmt"
	"strings"

	mcache "github.com/maxwell92/gokits/cache"
	mydeployment "app/backend/controller/yce/deploy"
	mytime "github.com/maxwell92/gokits/time"

	"k8s.io/kubernetes/pkg/api"
	kwatch "k8s.io/kubernetes/pkg/watch"

	client "k8s.io/kubernetes/pkg/client/unversioned"
	yceutils "app/backend/controller/yce/utils"
)

type PodWatcher struct {
	Events kwatch.Interface
	StopCh chan bool
	Cache  *mcache.RedisCache
	dcId   int32
	dcName string
	Fn     map[kwatch.EventType]func(kwatch.Event)
}

func (pw *PodWatcher) Aggregate(eventCh chan kwatch.Event) {
	for {
		select {
		case e := <-pw.Events.ResultChan():
			{
				eventCh <- e
				log.Tracef("KubeWatcher PodWatcher receive event of type %s", e.Type)
			}
		case _ = <-pw.StopCh:
			{
				log.Tracef("KubeWatcher PodWatcher receive Stop")
				pw.Events.ResultChan()
				return
			}
		}
	}
}

func (pw *PodWatcher) Stop() {
	pw.Events.Stop()
	log.Tracef("KubeWatcher PodWatcher Stop")
}

func (pw *PodWatcher) Handle(e kwatch.Event) {
	log.Tracef("KubeWatcher PodWatcher handle event type %s, %p", e.Type, &e)
	pw.Fn[e.Type](e)
}

func (pw *PodWatcher) Setup(cache *mcache.RedisCache, dcId int32, dcName string, c *client.Client) {
	pod, err := c.Pods(api.NamespaceAll).Watch(api.ListOptions{Watch: true, TimeoutSeconds: yceutils.Int64Ptr(mytime.DAYINSECONDS)})
	if err != nil {
		log.Errorf("PodWatcher Get Pod Watcher Error: error=%s", err)
		return
	}

	pw.StopCh = make(chan bool)
	pw.Cache = cache
	pw.Events = pod
	pw.dcId = dcId
	pw.dcName = dcName
	pw.Fn = make(map[kwatch.EventType]func(kwatch.Event))
	pw.Fn[kwatch.Added] = pw.Add
	pw.Fn[kwatch.Modified] = pw.Modified
	pw.Fn[kwatch.Deleted] = pw.Delete
	pw.Fn[kwatch.Error] = pw.Error
}

func (pw *PodWatcher) findPodInCache(key, podName string) *api.Pod {
	cacheDp := pw.Cache.Get(key)
	log.Debugf("KubeWatcher PodWatcher Get Key %s value %s", key, cacheDp)
	if cacheDp == "" {
		log.Errorf("PodWatcher findPodInCache Error")
		return nil
	}
	log.Debugf("PodWatcher findPodInCache %s", cacheDp)
	curDp := new(mydeployment.Deployments)
	curDp = curDp.Decode(cacheDp)
	for idx, p := range curDp.PodList.Items {
		if p.Name == podName {
			log.Debugf("KubeWatcher PodWatcher findPodInCache Deployment %s has pod: %s success", curDp.DeploymentName, podName)
			return &p
		}

		if idx == (len(curDp.PodList.Items) - 1) {
			log.Errorf("KubeWatcher PodWatcher Find Deployment %s has pod: %s failed", curDp.DeploymentName, podName)
			return nil
		}
	}

	log.Errorf("KubeWatcher Find Deployment has pod %s Error", podName)
	return nil
}

func (pw *PodWatcher) addIfNotExist(pod *api.Pod, podList *api.PodList) {
	// 有可能deploymentAdd先来，但是rs还没有创建出来，导致获取不到pod
	if len(podList.Items) == 0 {
		podList.Items = append(podList.Items, *pod)
		log.Tracef("KubeWatcher PodWatcher Add first Pod %s", pod.Name)
		return
	}
	for idx, p := range podList.Items {
		if p.Name == pod.Name {
			podList.Items[idx] = *pod
			log.Tracef("KubeWatcher PodWatcher addIfNotExit substitute pod %s ", pod.Name)
			break
		}

		if idx == (len(podList.Items) - 1) {
			podList.Items = append(podList.Items, *pod)
			log.Tracef("KubeWatcher PodWatcher addIfNotExist add pod %s ", pod.Name)
			return
		}
	}
}

func (pw *PodWatcher) getDeploymentFromPod(podName string) string {
	frag := strings.Split(podName, "-")
	length := len(frag)

	if length > 2 {
		return strings.Join(frag[:length-2], "-")
	} else {
		return podName
	}
}

func (pw *PodWatcher) Add(e kwatch.Event) {
	pod := e.Object.(*api.Pod)
	log.Debugf("KubeWatcher PodWatcher %s Event", e.Type)

	dpNamespace := pod.Namespace
	dpName := pw.getDeploymentFromPod(pod.Name)
	key := fmt.Sprintf("%d:%s:%s", pw.dcId, dpNamespace, dpName)
	cacheDp := pw.Cache.Get(key)
	log.Debugf("KubeWatcher PodWatcher Get Key %s value %s", key, cacheDp)
	if cacheDp == "" {
		log.Warnf("KubeWatche PodWatche get key %s empty", key)
		return
	}
	dp := new(mydeployment.Deployments)
	dp = dp.Decode(cacheDp)
	log.Debugf("KubeWatcher PodWatcher Decode with %d pod", len(dp.PodList.Items))
	// 为了解决DeploymentAdd事件先到，现有Deployment已经获取到了该pod，然后在此处重复添加的问题
	pw.addIfNotExist(pod, dp.PodList)
	data := dp.Encode()
	log.Debugf("KubeWatcher PodWatcher Encode with %d pod", len(dp.PodList.Items))
	ok, err := pw.Cache.Set(key, data)
	if err != nil || !ok {
		log.Errorf("KubeWatcher PodWatcher Add Error: error=%s", err)
		return
	}
	log.Debugf("KubeWatcher PodWatcher Add with %d pod", len(dp.PodList.Items))

	if pw.findPodInCache(key, pod.Name) != nil {
		log.Debugf("KubeWatcher PodWatcher Add find pod %s in cache", pod.Name)
	} else {
		log.Errorf("KubeWatcher PodWatcher Add find pod %s in cache error", pod.Name)
		return
	}

	log.Debugf("KubeWatcher PodWatcher Add pod %s success", pod.Name)
}

func (pw *PodWatcher) Modified(e kwatch.Event) {
	pod := e.Object.(*api.Pod)
	log.Debugf("KubeWatcher PodWatcher %s Event", e.Type)

	dpNamespace := pod.Namespace
	dpName := pw.getDeploymentFromPod(pod.Name)
	key := fmt.Sprintf("%d:%s:%s", pw.dcId, dpNamespace, dpName)
	if pw.findPodInCache(key, pod.Name) != nil {
		log.Debugf("KubeWatcher PodWatcher Modfied find pod %s in cache", pod.Name)
	} else {
		log.Errorf("KubeWatcher PodWatcher Modified find pod %s in cache error", pod.Name)
		return
	}

	cacheDp := pw.Cache.Get(key)
	log.Debugf("KubeWatcher PodWatcher Get Key %s value %s", key, cacheDp)
	if cacheDp == "" {
		log.Warnf("KubeWatcher PodWatcher Get key %s empty", key)
		return
	}
	dp := new(mydeployment.Deployments)
	dp = dp.Decode(cacheDp)

	for idx, p := range dp.PodList.Items {
		if p.Name == pod.Name {
			dp.PodList.Items[idx] = *pod
			log.Debugf("KubeWatcher PodWatcher Find pod %s ", pod.Name)
			break
		}

		if idx == (len(dp.PodList.Items) - 1) {
			log.Errorf("KubeWatcher PodWatcher Modified Error: pod %s not found", pod.Name)
			return
		}
	}

	data := dp.Encode()
	ok, err := pw.Cache.Set(key, data)
	if err != nil || !ok {
		log.Errorf("KubeWatcher PodWatcher Modified Error: error=%s", err)
		return
	}
	log.Debugf("KubeWatcher PodWatcher Modified pod %s success ", pod.Name)
}

func (pw *PodWatcher) Delete(e kwatch.Event) {
	pod := e.Object.(*api.Pod)
	log.Debugf("KubeWatcher PodWatcher %s Event", e.Type)

	dpNamespace := pod.Namespace
	dpName := pw.getDeploymentFromPod(pod.Name)
	key := fmt.Sprintf("%d:%s:%s", pw.dcId, dpNamespace, dpName)
	cacheDp := pw.Cache.Get(key)
	log.Debugf("KubeWatcher PodWatcher Get Key %s value %s", key, cacheDp)
	if cacheDp == "" {
		log.Warnf("KubeWatcher PodWatcher get key %s empty", key)
		return
	}
	dp := new(mydeployment.Deployments)
	dp = dp.Decode(cacheDp)
	for idx, p := range dp.PodList.Items {
		if p.Name == pod.Name {
			dp.PodList.Items[idx] = api.Pod{}
			start := dp.PodList.Items[:idx]
			end := dp.PodList.Items[idx+1:]
			start = append(start, end...)
			dp.PodList.Items = start
			log.Debugf("KubeWatcher PodWatcher Find pod %s", p.Name)
			break
		}

		if idx == (len(dp.PodList.Items) - 1) {
			log.Errorf("SuperWatcher PodWatcher Error: pod %s not found", pod.Name)
			return
		}
	}
	log.Debugf("Deployment %s now %d pod", dp.DeploymentName, len(dp.PodList.Items))

	data := dp.Encode()
	ok, err := pw.Cache.Set(key, data)
	if err != nil || !ok {
		log.Errorf("KubeWatcher PodWatcher Delete Error")
		return
	}
	if pw.findPodInCache(dpName, pod.Name) != nil {
		log.Debugf("KubeWatcher PodWatcher Delete find pod %s in cache", pod.Name)
	} else {
		log.Errorf("KubeWatcher PodWatcher Delete find pod %s in cache error", pod.Name)
		return
	}

	log.Debugf("KubeWatcher PodWatcher Delete Success")
}

func (pw *PodWatcher) Error(e kwatch.Event) {
	log.Errorf("KubeWatcher PodWatcher Error")
}
