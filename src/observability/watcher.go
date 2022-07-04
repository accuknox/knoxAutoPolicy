package observability

import (
	"encoding/json"
	"io"

	"github.com/accuknox/auto-policy-discovery/src/cluster"
	"github.com/accuknox/auto-policy-discovery/src/types"
)

func addPodToList(pod types.Pod) {
	podExist := false

	for _, locpod := range Pods {
		if locpod.IP == pod.IP && locpod.PodName == pod.PodName && locpod.Namespace == pod.Namespace {
			podExist = true
			break
		}
	}
	if !podExist {
		Pods = append(Pods, types.Pod{
			Namespace: pod.Namespace,
			PodName:   pod.PodName,
			IP:        pod.IP,
		})
	}
}

func updatePodList(pod types.Pod) {
	for index, locpod := range Pods {
		if locpod.PodName == pod.PodName && locpod.Namespace == pod.Namespace {
			Pods[index].IP = pod.IP
			break
		}
	}
}

// WatchK8sPods Function
func WatchK8sPods() {
	for {
		if resp := cluster.WatchK8sPods(); resp != nil {
			defer resp.Body.Close()

			decoder := json.NewDecoder(resp.Body)
			for {
				event := types.K8sPodEvent{}
				if err := decoder.Decode(&event); err == io.EOF {
					break
				} else if err != nil {
					break
				}

				if event.Type != "ADDED" && event.Type != "MODIFIED" && event.Type != "DELETED" {
					continue
				}

				pod := types.Pod{
					Namespace: event.Object.ObjectMeta.Namespace,
					PodName:   event.Object.ObjectMeta.Name,
					IP:        event.Object.Status.PodIP,
				}

				if event.Type == "ADDED" {
					addPodToList(pod)
				} else if event.Type == "MODIFIED" {
					updatePodList(pod)
				}
			}
		}
	}
}

func GetPodName(ip string) string {
	for _, pod := range Pods {
		if pod.IP == ip {
			return pod.PodName
		}
	}
	return ip
}
