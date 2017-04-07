package kubewatcher

import (
	mycache "gokits/cache"
	mydeployment "app/backend/controller/yce/deploy"
	"fmt"
	"testing"
	/*
	yceutils "app/backend/controller/yce/utils"
	"k8s.io/kubernetes/pkg/api"
	"time"
	kwatch "k8s.io/kubernetes/pkg/watch"
	"strconv"
	*/
)

/*
func TestWatcherTImeout(t *testing.T) {
	cli, _ := yceutils.CreateK8sClient("172.21.1.11:8080")
	watcher1, _ := cli.Pods(api.NamespaceAll).Watch(api.ListOptions{})
	watcher2, _ := cli.Pods(api.NamespaceAll).Watch(api.ListOptions{Watch: true})

	eventCh1 := watcher1.ResultChan()
	eventCh2 := watcher2.ResultChan()

}
*/
func TestKubeWatcher_Run(t *testing.T) {

	cache := mycache.BuntCacheInstance()
	cache.Create("DeploymentCache")
	cache.Index([]string{mycache.UserName, mycache.DcName, mycache.DcId, mycache.DeploymentName, mycache.OrgName})

	cache.Update("1:ops:testhealthz","{\"userName\":\"liyao.miao\",\"dcName\":\"办公网2\",\"dcId\":66,\"orgName\":\"ops\",\"deploymentName\":\"testhealthz\",\"updateTime\":\"2017-03-23 23:57:51 +0800 CST\",\"deployment\":{\"metadata\":{\"name\":\"testhealthz\",\"namespace\":\"ops\",\"selfLink\":\"/apis/extensions/v1beta1/namespaces/ops/deployments/testhealthz\",\"uid\":\"783d879e-0fe1-11e7-9bb2-005056915fa1\",\"resourceVersion\":\"18689564\",\"generation\":23,\"creationTimestamp\":\"2017-03-23T15:57:51Z\",\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"},\"annotations\":{\"deployment.kubernetes.io/revision\":\"3\"}},\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"template\":{\"metadata\":{\"creationTimestamp\":null,\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"spec\":{\"volumes\":null,\"containers\":[{\"name\":\"testhealthz\",\"image\":\"img.reg.3g:15000/healthz:2.0\",\"ports\":[{\"name\":\"port1\",\"containerPort\":8080,\"protocol\":\"TCP\"}],\"resources\":{\"limits\":{\"cpu\":\"1\",\"memory\":\"2G\"},\"requests\":{\"cpu\":\"500m\",\"memory\":\"1G\"}},\"livenessProbe\":{\"httpGet\":{\"path\":\"/healthz\",\"port\":8080,\"scheme\":\"HTTP\"},\"initialDelaySeconds\":5,\"timeoutSeconds\":1,\"periodSeconds\":5,\"successThreshold\":1,\"failureThreshold\":3},\"terminationMessagePath\":\"/dev/termination-log\",\"imagePullPolicy\":\"IfNotPresent\"}],\"restartPolicy\":\"Always\",\"terminationGracePeriodSeconds\":30,\"dnsPolicy\":\"ClusterFirst\",\"serviceAccountName\":\"\",\"securityContext\":{}}},\"strategy\":{\"type\":\"RollingUpdate\",\"rollingUpdate\":{\"maxUnavailable\":2,\"maxSurge\":2}}},\"status\":{\"observedGeneration\":23,\"replicas\":1,\"updatedReplicas\":1,\"unavailableReplicas\":1}},\"podList\":{\"metadata\":{},\"items\":[]},\"service\":{\"metadata\":{\"creationTimestamp\":null},\"spec\":{\"ports\":null,\"selector\":null},\"status\":{\"loadBalancer\":{}}}}")
	cache.Update("1:ops:testhealthz1","{\"userName\":\"liyao.miao\",\"dcName\":\"办公网2\",\"dcId\":66,\"orgName\":\"ops\",\"deploymentName\":\"testhealthz1\",\"updateTime\":\"2017-03-23 23:57:51 +0800 CST\",\"deployment\":{\"metadata\":{\"name\":\"testhealthz\",\"namespace\":\"ops\",\"selfLink\":\"/apis/extensions/v1beta1/namespaces/ops/deployments/testhealthz\",\"uid\":\"783d879e-0fe1-11e7-9bb2-005056915fa1\",\"resourceVersion\":\"18689564\",\"generation\":23,\"creationTimestamp\":\"2017-03-23T15:57:51Z\",\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"},\"annotations\":{\"deployment.kubernetes.io/revision\":\"3\"}},\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"template\":{\"metadata\":{\"creationTimestamp\":null,\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"spec\":{\"volumes\":null,\"containers\":[{\"name\":\"testhealthz\",\"image\":\"img.reg.3g:15000/healthz:2.0\",\"ports\":[{\"name\":\"port1\",\"containerPort\":8080,\"protocol\":\"TCP\"}],\"resources\":{\"limits\":{\"cpu\":\"1\",\"memory\":\"2G\"},\"requests\":{\"cpu\":\"500m\",\"memory\":\"1G\"}},\"livenessProbe\":{\"httpGet\":{\"path\":\"/healthz\",\"port\":8080,\"scheme\":\"HTTP\"},\"initialDelaySeconds\":5,\"timeoutSeconds\":1,\"periodSeconds\":5,\"successThreshold\":1,\"failureThreshold\":3},\"terminationMessagePath\":\"/dev/termination-log\",\"imagePullPolicy\":\"IfNotPresent\"}],\"restartPolicy\":\"Always\",\"terminationGracePeriodSeconds\":30,\"dnsPolicy\":\"ClusterFirst\",\"serviceAccountName\":\"\",\"securityContext\":{}}},\"strategy\":{\"type\":\"RollingUpdate\",\"rollingUpdate\":{\"maxUnavailable\":2,\"maxSurge\":2}}},\"status\":{\"observedGeneration\":23,\"replicas\":1,\"updatedReplicas\":1,\"unavailableReplicas\":1}},\"podList\":{\"metadata\":{},\"items\":[]},\"service\":{\"metadata\":{\"creationTimestamp\":null},\"spec\":{\"ports\":null,\"selector\":null},\"status\":{\"loadBalancer\":{}}}}")
	cache.Update("1:ops:testhealthz2","{\"userName\":\"liyao.miao\",\"dcName\":\"办公网2\",\"dcId\":66,\"orgName:\"ops\",\"deploymentName\":\"testhealthz2\",\"updateTime\":\"2017-03-23 23:57:51 +0800 CST\",\"deployment\":{\"metadata\":{\"name\":\"testhealthz\",\"namespace\":\"ops\",\"selfLink\":\"/apis/extensions/v1beta1/namespaces/ops/deployments/testhealthz\",\"uid\":\"783d879e-0fe1-11e7-9bb2-005056915fa1\",\"resourceVersion\":\"18689564\",\"generation\":23,\"creationTimestamp\":\"2017-03-23T15:57:51Z\",\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"},\"annotations\":{\"deployment.kubernetes.io/revision\":\"3\"}},\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"template\":{\"metadata\":{\"creationTimestamp\":null,\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"spec\":{\"volumes\":null,\"containers\":[{\"name\":\"testhealthz\",\"image\":\"img.reg.3g:15000/healthz:2.0\",\"ports\":[{\"name\":\"port1\",\"containerPort\":8080,\"protocol\":\"TCP\"}],\"resources\":{\"limits\":{\"cpu\":\"1\",\"memory\":\"2G\"},\"requests\":{\"cpu\":\"500m\",\"memory\":\"1G\"}},\"livenessProbe\":{\"httpGet\":{\"path\":\"/healthz\",\"port\":8080,\"scheme\":\"HTTP\"},\"initialDelaySeconds\":5,\"timeoutSeconds\":1,\"periodSeconds\":5,\"successThreshold\":1,\"failureThreshold\":3},\"terminationMessagePath\":\"/dev/termination-log\",\"imagePullPolicy\":\"IfNotPresent\"}],\"restartPolicy\":\"Always\",\"terminationGracePeriodSeconds\":30,\"dnsPolicy\":\"ClusterFirst\",\"serviceAccountName\":\"\",\"securityContext\":{}}},\"strategy\":{\"type\":\"RollingUpdate\",\"rollingUpdate\":{\"maxUnavailable\":2,\"maxSurge\":2}}},\"status\":{\"observedGeneration\":23,\"replicas\":1,\"updatedReplicas\":1,\"unavailableReplicas\":1}},\"podList\":{\"metadata\":{},\"items\":[]},\"service\":{\"metadata\":{\"creationTimestamp\":null},\"spec\":{\"ports\":null,\"selector\":null},\"status\":{\"loadBalancer\":{}}}}")
	cache.Update("1:ops:testhealthz3","{\"userName\":\"liyao.miao\",\"dcName\":\"办公网2\",\"dcId\":66,\"orgName:\"ops\",\"deploymentName\":\"testhealthz3\",\"updateTime\":\"2017-03-23 23:57:51 +0800 CST\",\"deployment\":{\"metadata\":{\"name\":\"testhealthz\",\"namespace\":\"ops\",\"selfLink\":\"/apis/extensions/v1beta1/namespaces/ops/deployments/testhealthz\",\"uid\":\"783d879e-0fe1-11e7-9bb2-005056915fa1\",\"resourceVersion\":\"18689564\",\"generation\":23,\"creationTimestamp\":\"2017-03-23T15:57:51Z\",\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"},\"annotations\":{\"deployment.kubernetes.io/revision\":\"3\"}},\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"template\":{\"metadata\":{\"creationTimestamp\":null,\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"spec\":{\"volumes\":null,\"containers\":[{\"name\":\"testhealthz\",\"image\":\"img.reg.3g:15000/healthz:2.0\",\"ports\":[{\"name\":\"port1\",\"containerPort\":8080,\"protocol\":\"TCP\"}],\"resources\":{\"limits\":{\"cpu\":\"1\",\"memory\":\"2G\"},\"requests\":{\"cpu\":\"500m\",\"memory\":\"1G\"}},\"livenessProbe\":{\"httpGet\":{\"path\":\"/healthz\",\"port\":8080,\"scheme\":\"HTTP\"},\"initialDelaySeconds\":5,\"timeoutSeconds\":1,\"periodSeconds\":5,\"successThreshold\":1,\"failureThreshold\":3},\"terminationMessagePath\":\"/dev/termination-log\",\"imagePullPolicy\":\"IfNotPresent\"}],\"restartPolicy\":\"Always\",\"terminationGracePeriodSeconds\":30,\"dnsPolicy\":\"ClusterFirst\",\"serviceAccountName\":\"\",\"securityContext\":{}}},\"strategy\":{\"type\":\"RollingUpdate\",\"rollingUpdate\":{\"maxUnavailable\":2,\"maxSurge\":2}}},\"status\":{\"observedGeneration\":23,\"replicas\":1,\"updatedReplicas\":1,\"unavailableReplicas\":1}},\"podList\":{\"metadata\":{},\"items\":[]},\"service\":{\"metadata\":{\"creationTimestamp\":null},\"spec\":{\"ports\":null,\"selector\":null},\"status\":{\"loadBalancer\":{}}}}")
	cache.Update("66:ops:testhealthz","{\"userName\":\"liyao.miao\",\"dcName\":\"办公网2\",\"dcId\":66,\"orgName:\"ops\",\"deploymentName\":\"testhealthz\",\"updateTime\":\"2017-03-23 23:57:51 +0800 CST\",\"deployment\":{\"metadata\":{\"name\":\"testhealthz\",\"namespace\":\"ops\",\"selfLink\":\"/apis/extensions/v1beta1/namespaces/ops/deployments/testhealthz\",\"uid\":\"783d879e-0fe1-11e7-9bb2-005056915fa1\",\"resourceVersion\":\"18689564\",\"generation\":23,\"creationTimestamp\":\"2017-03-23T15:57:51Z\",\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"},\"annotations\":{\"deployment.kubernetes.io/revision\":\"3\"}},\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"template\":{\"metadata\":{\"creationTimestamp\":null,\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"spec\":{\"volumes\":null,\"containers\":[{\"name\":\"testhealthz\",\"image\":\"img.reg.3g:15000/healthz:2.0\",\"ports\":[{\"name\":\"port1\",\"containerPort\":8080,\"protocol\":\"TCP\"}],\"resources\":{\"limits\":{\"cpu\":\"1\",\"memory\":\"2G\"},\"requests\":{\"cpu\":\"500m\",\"memory\":\"1G\"}},\"livenessProbe\":{\"httpGet\":{\"path\":\"/healthz\",\"port\":8080,\"scheme\":\"HTTP\"},\"initialDelaySeconds\":5,\"timeoutSeconds\":1,\"periodSeconds\":5,\"successThreshold\":1,\"failureThreshold\":3},\"terminationMessagePath\":\"/dev/termination-log\",\"imagePullPolicy\":\"IfNotPresent\"}],\"restartPolicy\":\"Always\",\"terminationGracePeriodSeconds\":30,\"dnsPolicy\":\"ClusterFirst\",\"serviceAccountName\":\"\",\"securityContext\":{}}},\"strategy\":{\"type\":\"RollingUpdate\",\"rollingUpdate\":{\"maxUnavailable\":2,\"maxSurge\":2}}},\"status\":{\"observedGeneration\":23,\"replicas\":1,\"updatedReplicas\":1,\"unavailableReplicas\":1}},\"podList\":{\"metadata\":{},\"items\":[]},\"service\":{\"metadata\":{\"creationTimestamp\":null},\"spec\":{\"ports\":null,\"selector\":null},\"status\":{\"loadBalancer\":{}}}}")
	cache.Update("66:ops:testhealthz1","{\"userName\":\"liyao.miao\",\"dcName\":\"办公网2\",\"dcId\":66,\"orgName:\"ops\",\"deploymentName\":\"testhealthz1\",\"updateTime\":\"2017-03-23 23:57:51 +0800 CST\",\"deployment\":{\"metadata\":{\"name\":\"testhealthz\",\"namespace\":\"ops\",\"selfLink\":\"/apis/extensions/v1beta1/namespaces/ops/deployments/testhealthz\",\"uid\":\"783d879e-0fe1-11e7-9bb2-005056915fa1\",\"resourceVersion\":\"18689564\",\"generation\":23,\"creationTimestamp\":\"2017-03-23T15:57:51Z\",\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"},\"annotations\":{\"deployment.kubernetes.io/revision\":\"3\"}},\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"template\":{\"metadata\":{\"creationTimestamp\":null,\"labels\":{\"author\":\"liyao.miao\",\"name\":\"testhealthz\",\"version\":\"2.0\"}},\"spec\":{\"volumes\":null,\"containers\":[{\"name\":\"testhealthz\",\"image\":\"img.reg.3g:15000/healthz:2.0\",\"ports\":[{\"name\":\"port1\",\"containerPort\":8080,\"protocol\":\"TCP\"}],\"resources\":{\"limits\":{\"cpu\":\"1\",\"memory\":\"2G\"},\"requests\":{\"cpu\":\"500m\",\"memory\":\"1G\"}},\"livenessProbe\":{\"httpGet\":{\"path\":\"/healthz\",\"port\":8080,\"scheme\":\"HTTP\"},\"initialDelaySeconds\":5,\"timeoutSeconds\":1,\"periodSeconds\":5,\"successThreshold\":1,\"failureThreshold\":3},\"terminationMessagePath\":\"/dev/termination-log\",\"imagePullPolicy\":\"IfNotPresent\"}],\"restartPolicy\":\"Always\",\"terminationGracePeriodSeconds\":30,\"dnsPolicy\":\"ClusterFirst\",\"serviceAccountName\":\"\",\"securityContext\":{}}},\"strategy\":{\"type\":\"RollingUpdate\",\"rollingUpdate\":{\"maxUnavailable\":2,\"maxSurge\":2}}},\"status\":{\"observedGeneration\":23,\"replicas\":1,\"updatedReplicas\":1,\"unavailableReplicas\":1}},\"podList\":{\"metadata\":{},\"items\":[]},\"service\":{\"metadata\":{\"creationTimestamp\":null},\"spec\":{\"ports\":null,\"selector\":null},\"status\":{\"loadBalancer\":{}}}}")

	cache.Update("26:ops:app-1", "{\"userName\":\"xueyan.xu\",\"dcName\":\"电信\",\"dcId\":26,\"orgName\":\"ops\",\"deploymentName\":\"app-1\",\"updateTime\":\"2017-02-2217: 48: 41+0800CST\"}")
	cache.Update("1:dev:a1", "{\"userName\":\"liyao.miao\",\"dcName\":\"办公网\",\"dcId\":1,\"orgName\":\"dev\",\"deploymentName\":\"a1\",\"updateTime\":\"2017-02-2217: 45: 30+0800CST\"}")
	cache.Update("1:dev:a6","{\"userName\":\"liyao.miao\",\"dcName\":\"办公网\",\"dcId\":1,\"orgName\":\"dev\",\"deploymentName\":\"a6\",\"updateTime\":\"2017-02-2217: 49: 19+0800CST\"}")
	cache.Update("1:ops:a2", "{\"userName\":\"liyao.miao\",\"dcName\":\"办公网\",\"dcId\":1,\"orgName\":\"ops\",\"deploymentName\":\"a2\",\"updateTime\":\"2017-02-2217: 45: 56+0800CST\"}")
	cache.Update("26:dev:a1", "{\"userName\":\"liyao.miao\",\"dcName\":\"电信\",\"dcId\":26,\"orgName\":\"dev\",\"deploymentName\":\"a1\",\"updateTime\":\"2017-02-2217: 45: 27+0800CST\"}")
	cache.Update("26:dev:a4", "{\"userName\":\"liyao.miao\",\"dcName\":\"电信\",\"dcId\":26,\"orgName\":\"dev\",\"deploymentName\":\"a4\",\"updateTime\":\"2017-02-2217: 47: 08+0800CST\"}")
	cache.Update("26:ops:a6", "{\"userName\":\"liyao.miao\",\"dcName\":\"电信\",\"dcId\":26,\"orgName\":\"ops\",\"deploymentName\":\"a6\",\"updateTime\":\"2017-02-2217: 49: 15+0800CST\"}")


	multi := make([]mycache.Indexer, 1)
	multi[0] = mycache.Indexer{mycache.OrgName, "ops"}
	//multi[0] = mycache.Indexer{mycache.DeploymentName, "a1"}

	results, err := cache.Find(multi...)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println( len(*results))
		for _, s := range *results {
			oldDp := new(mydeployment.Deployments)
			oldDp = oldDp.Decode(s)
			fmt.Println(oldDp.DeploymentName)
		}
	}

	/*
	kw := NewKubeWatcher("172.21.1.11:8080")
	dp, err := kw.Client.Deployments("ops").Watch(api.ListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	*/

	/*
	cache := mycache.BuntCacheInstance()
	cache.Create("DeploymentCache")
	cache.Index([]string{mycache.UserName, mycache.DcName, mycache.DcId, mycache.DeploymentName, mycache.OrgName})

	multi := make([]mycache.Indexer, 3)
	multi[0] = mycache.Indexer{mycache.DcId, "1"}
	multi[1] = mycache.Indexer{mycache.OrgName, "ops"}
	multi[2] = mycache.Indexer{mycache.DeploymentName, "testhealthz"}

	// results, err := mycache.BuntCacheInstance().Find(multi...)
	results, err := cache.Find(multi...)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(len(*results))
	}
	*/

	/*
	deploy := dp.ResultChan()
	fmt.Printf("%p\n", deploy)

	kw.Deployments.Setup(deploy, cache, 1, "")

	go kw.Deployments.Dispatch()

	watch := func(event <-chan kwatch.Event) {
		for {
			select {
			case e := <-event:
				{
					fmt.Println(e.Type)
				}
			}
		}
	}
	go watch(deploy)

	time.Sleep(time.Duration(1000) * time.Second)
	*/
}
