package cmd

import (
	// "fmt"
	"context"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/api/core/v1"
	//"gopkg.in/oleiade/reflections.v1"
	//"reflect"
)

type Pod struct {
	Client *kubernetes.Clientset
	Namespace string
}

type PodObj struct {
	Name string `json:"name"`
	Status v1.PodPhase `json:"status"`
	Labels map[string]string `json:"labels"`
}

func (p *Pod) GetPods(podLabels map[string]string) ([]PodObj, error){

	pods, err := p.Client.CoreV1().Pods(p.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labels.Set(podLabels).String() })
	if err != nil {
		return nil, err
	}

	podObjList := []PodObj{}

	for _, pod := range pods.Items {
		
		podObj := PodObj {
			Name: pod.Name,
			Status: pod.Status.Phase,
			Labels: pod.Labels,
		}
		
		podObjList = append(podObjList, podObj)
	}

	return podObjList, nil
}

// func CheckPodPhase(phase v1.PodPhase) string {
// 	switch phase {
// 		case v1.PodRunning:
// 			return "Running"
// 		case v1.PodPending:
// 			return "Pending" 
// 		case v1.PodSucceeded:
// 			return "Succeded"
// 		case v1.PodFailed:
// 			return "Failed"
// 		default:
// 			return "unknown"
// 	}
// }