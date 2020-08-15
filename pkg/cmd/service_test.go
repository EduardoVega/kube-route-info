package cmd

import (
	"bytes"
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	v1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
)

// Mock client struct
type ServiceMockClient struct{}

func NewServiceMockClient() *ServiceMockClient {
	return &ServiceMockClient{}
}

func (c *ServiceMockClient) GetPodsByLabels(labels map[string]string) (*v1.PodList, error) {

	podList := []v1.Pod{
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-foo-1",
				Namespace: "default",
				Labels: map[string]string{
					"app":     "foo",
					"version": "v1",
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-foo-2",
				Namespace: "default",
				Labels: map[string]string{
					"app":     "foo",
					"version": "v1",
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-foo-3",
				Namespace: "default",
				Labels: map[string]string{
					"app":     "foo",
					"version": "v1",
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-foo-4",
				Namespace: "default",
				Labels: map[string]string{
					"app":     "foo",
					"version": "v2",
				},
			},
		},
	}

	matchedPodList := v1.PodList{
		Items: []v1.Pod{},
	}

	for _, pod := range podList {
		if reflect.DeepEqual(pod.ObjectMeta.Labels, labels) {
			matchedPodList.Items = append(matchedPodList.Items, pod)
		}
	}

	return &matchedPodList, nil
}

func (c *ServiceMockClient) GetServiceByName(name string) (*v1.Service, error) {
	serviceList := []v1.Service{
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-clusterip",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeClusterIP,
				Ports: []v1.ServicePort{
					{Port: 80, TargetPort: intstr.IntOrString{Type: 0, IntVal: 80}},
					{Port: 443, TargetPort: intstr.IntOrString{Type: 1, StrVal: "https"}},
				},
				Selector: map[string]string{
					"app":     "foo",
					"version": "v1",
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-clusterip-no-pods",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeClusterIP,
				Ports: []v1.ServicePort{
					{Port: 80, TargetPort: intstr.IntOrString{Type: 0, IntVal: 80}},
					{Port: 443, TargetPort: intstr.IntOrString{Type: 1, StrVal: "https"}},
				},
				Selector: map[string]string{
					"app":     "bar",
					"version": "v1",
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-nodeport",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeNodePort,
				Ports: []v1.ServicePort{
					{Port: 80, TargetPort: intstr.IntOrString{Type: 1, StrVal: "http"}, NodePort: 1234},
				},
				Selector: map[string]string{
					"app":     "foo",
					"version": "v2",
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-loadbalancer",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeLoadBalancer,
				Ports: []v1.ServicePort{
					{Port: 80, TargetPort: intstr.IntOrString{Type: 0, IntVal: 80}, NodePort: 1234},
				},
				Selector: map[string]string{
					"app":     "foo",
					"version": "v1",
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-externalname",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type:         v1.ServiceTypeExternalName,
				ExternalName: "my.external.app.com",
			},
		},
	}

	for _, service := range serviceList {
		if service.Name == name {
			return &service, nil
		}
	}

	return nil, nil
}

func (c *ServiceMockClient) GetIngressByName(name string) (*v1beta1.Ingress, error) { return nil, nil }

func TestServicePrintGraphSuccessful(t *testing.T) {

	tests := []struct {
		serviceName   string
		expectedGraph string
	}{
		{
			"service-clusterip",
			"[Service]  service-clusterip\n├── [Pod]  pod-foo-1\n├── [Pod]  pod-foo-2\n└── [Pod]  pod-foo-3\n",
		},
		{
			"service-clusterip-no-pods",
			"[Service]  service-clusterip-no-pods\n",
		},
		{
			"service-nodeport",
			"[Service]  service-nodeport\n└── [Pod]  pod-foo-4\n",
		},
		{
			"service-loadbalancer",
			"[Service]  service-loadbalancer\n├── [Pod]  pod-foo-1\n├── [Pod]  pod-foo-2\n└── [Pod]  pod-foo-3\n",
		},
		{
			"service-externalname",
			"[Service]  service-externalname\n└── [Hostname]  my.external.app.com\n",
		},
	}

	for _, test := range tests {

		mockClient := NewServiceMockClient()

		service := NewService(mockClient, "default")

		buf := &bytes.Buffer{}

		service.PrintGraph(test.serviceName, buf)

		if buf.String() != test.expectedGraph {
			t.Errorf("Returned tree graph was incorrect,\ngot:\n%s\nwant:\n%s", buf.String(), test.expectedGraph)
		}
	}
}

func TestServicePrintTableSuccessful(t *testing.T) {

	tests := []struct {
		serviceName   string
		expectedTable string
	}{
		{
			"service-clusterip",
			"NAME                TYPE        PORT(S)           POD(S)\nservice-clusterip   ClusterIP   80 80,443 https   pod-foo-1,pod-foo-2,pod-foo-3\n",
		},
		{
			"service-clusterip-no-pods",
			"NAME                        TYPE        PORT(S)           POD(S)\nservice-clusterip-no-pods   ClusterIP   80 80,443 https   \n",
		},
		{
			"service-nodeport",
			"NAME               TYPE       PORT(S)        POD(S)\nservice-nodeport   NodePort   80 http 1234   pod-foo-4\n",
		},
		{
			"service-loadbalancer",
			"NAME                   TYPE           PORT(S)      POD(S)\nservice-loadbalancer   LoadBalancer   80 80 1234   pod-foo-1,pod-foo-2,pod-foo-3\n",
		},
		{
			"service-externalname",
			"NAME                   TYPE           PORT(S)   HOSTNAME\nservice-externalname   ExternalName             my.external.app.com\n",
		},
	}

	for _, test := range tests {

		mockClient := NewServiceMockClient()

		service := NewService(mockClient, "default")

		buf := &bytes.Buffer{}

		service.PrintTable(test.serviceName, buf)

		if buf.String() != test.expectedTable {
			t.Errorf("Returned table was incorrect,\ngot:\n%swant:\n%s", buf.String(), test.expectedTable)
		}
	}
}
