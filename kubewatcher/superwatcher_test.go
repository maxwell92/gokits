package kubewatcher

import (
"testing"
"time"
)


func TestWatcher_Run(t *testing.T) {

	apiServers := []string{"172.21.1.11:8080", "172.17.106.52:8080"}
	sw := NewSuperWatcher(apiServers)
	go sw.Run()
	time.Sleep(time.Duration(1000) * time.Second)
}
