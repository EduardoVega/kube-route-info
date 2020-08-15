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
type IngressMockClient struct{}

func NewIngressMockClient() *IngressMockClient {
	return &IngressMockClient{}
}

func (c *IngressMockClient) GetPodsByLabels(labels map[string]string) (*v1.PodList, error) {
	podList := []v1.Pod{
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-foo-1",
				Namespace: "default",
				Labels: map[string]string{
					"app": "foo",
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-foo-2",
				Namespace: "default",
				Labels: map[string]string{
					"app": "foo",
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-bar-1",
				Namespace: "default",
				Labels: map[string]string{
					"app": "bar",
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

func (c *IngressMockClient) GetServiceByName(name string) (*v1.Service, error) {
	serviceList := []v1.Service{
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-foo",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeClusterIP,
				Ports: []v1.ServicePort{
					{Port: 80, TargetPort: intstr.IntOrString{Type: 0, IntVal: 80}},
					{Port: 443, TargetPort: intstr.IntOrString{Type: 1, StrVal: "https"}},
				},
				Selector: map[string]string{
					"app": "foo",
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-bar",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeClusterIP,
				Ports: []v1.ServicePort{
					{Port: 80, TargetPort: intstr.IntOrString{Type: 1, StrVal: "http"}},
				},
				Selector: map[string]string{
					"app": "bar",
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

func (c *IngressMockClient) GetIngressByName(name string) (*v1beta1.Ingress, error) {
	ingressList := []v1beta1.Ingress{
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "networking.k8s.io/v1beta1", Kind: "Ingress"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ingress-1-backend",
				Namespace: "default",
			},
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: "v1.ingress.com",
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: "",
										Backend: v1beta1.IngressBackend{
											ServiceName: "service-foo",
											ServicePort: intstr.IntOrString{
												Type:   0,
												IntVal: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{APIVersion: "networking.k8s.io/v1beta1", Kind: "Ingress"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ingress-2-backends",
				Namespace: "default",
			},
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: "",
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: "/foo",
										Backend: v1beta1.IngressBackend{
											ServiceName: "service-foo",
											ServicePort: intstr.IntOrString{
												Type:   0,
												IntVal: 80,
											},
										},
									},
									{
										Path: "/bar",
										Backend: v1beta1.IngressBackend{
											ServiceName: "service-bar",
											ServicePort: intstr.IntOrString{
												Type:   0,
												IntVal: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, {
			TypeMeta: metav1.TypeMeta{APIVersion: "networking.k8s.io/v1beta1", Kind: "Ingress"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ingress-2-backends-2-rules",
				Namespace: "default",
			},
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: "1.rule.com",
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: "/foo",
										Backend: v1beta1.IngressBackend{
											ServiceName: "service-foo",
											ServicePort: intstr.IntOrString{
												Type:   0,
												IntVal: 80,
											},
										},
									},
									{
										Path: "/bar",
										Backend: v1beta1.IngressBackend{
											ServiceName: "service-bar",
											ServicePort: intstr.IntOrString{
												Type:   0,
												IntVal: 80,
											},
										},
									},
								},
							},
						},
					},
					{
						Host: "2.rule.com",
						IngressRuleValue: v1beta1.IngressRuleValue{
							HTTP: &v1beta1.HTTPIngressRuleValue{
								Paths: []v1beta1.HTTPIngressPath{
									{
										Path: "/externalname",
										Backend: v1beta1.IngressBackend{
											ServiceName: "service-externalname",
											ServicePort: intstr.IntOrString{
												Type:   0,
												IntVal: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, ingress := range ingressList {
		if ingress.Name == name {
			return &ingress, nil
		}
	}

	return nil, nil
}

func TestIngressPrintGraphSuccessful(t *testing.T) {

	tests := []struct {
		ingressName   string
		expectedGraph string
	}{
		{
			"ingress-1-backend",
			"[Ingress]  ingress-1-backend\n└── v1.ingress.com\n    └── \n        └── [Service]  service-foo\n            ├── [Pod]  pod-foo-1\n            └── [Pod]  pod-foo-2\n",
		},
		{
			"ingress-2-backends",
			"[Ingress]  ingress-2-backends\n└── \n    ├── /foo\n    │   └── [Service]  service-foo\n    │       ├── [Pod]  pod-foo-1\n    │       └── [Pod]  pod-foo-2\n    └── /bar\n        └── [Service]  service-bar\n            └── [Pod]  pod-bar-1\n",
		},
		{
			"ingress-2-backends-2-rules",
			"[Ingress]  ingress-2-backends-2-rules\n├── 1.rule.com\n│   ├── /foo\n│   │   └── [Service]  service-foo\n│   │       ├── [Pod]  pod-foo-1\n│   │       └── [Pod]  pod-foo-2\n│   └── /bar\n│       └── [Service]  service-bar\n│           └── [Pod]  pod-bar-1\n└── 2.rule.com\n    └── /externalname\n        └── [Service]  service-externalname\n            └── [Hostname]  my.external.app.com\n",
		},
	}

	for _, test := range tests {

		mockClient := NewIngressMockClient()

		ingress := NewIngress(mockClient, "default")

		buf := &bytes.Buffer{}

		ingress.PrintGraph(test.ingressName, buf)

		if buf.String() != test.expectedGraph {
			t.Errorf("Returned tree graph was incorrect,\ngot:\n%s\nwant:\n%s", buf.String(), test.expectedGraph)
		}
	}
}

func TestIngressPrintTableSuccessful(t *testing.T) {

	tests := []struct {
		ingressName   string
		expectedTable string
	}{
		{
			"ingress-1-backend",
			"NAME                HOST             PATH   PORT   SERVICE       TYPE        SERVICE PORT(S)   POD(S)\ningress-1-backend   v1.ingress.com          80     service-foo   ClusterIP   80 80,443 https   pod-foo-1,pod-foo-2\n",
		},
		{
			"ingress-2-backends",
			"NAME                 HOST   PATH   PORT   SERVICE       TYPE        SERVICE PORT(S)   POD(S)\ningress-2-backends          /foo   80     service-foo   ClusterIP   80 80,443 https   pod-foo-1,pod-foo-2\ningress-2-backends          /bar   80     service-bar   ClusterIP   80 http           pod-bar-1\n",
		},
		{
			"ingress-2-backends-2-rules",
			"NAME                         HOST         PATH            PORT   SERVICE                TYPE           SERVICE PORT(S)   POD(S)/HOSTNAME\ningress-2-backends-2-rules   1.rule.com   /foo            80     service-foo            ClusterIP      80 80,443 https   pod-foo-1,pod-foo-2\ningress-2-backends-2-rules   1.rule.com   /bar            80     service-bar            ClusterIP      80 http           pod-bar-1\ningress-2-backends-2-rules   2.rule.com   /externalname   80     service-externalname   ExternalName                     my.external.app.com\n",
		},
	}

	for _, test := range tests {

		mockClient := NewIngressMockClient()

		ingress := NewIngress(mockClient, "default")

		buf := &bytes.Buffer{}

		ingress.PrintTable(test.ingressName, buf)

		if buf.String() != test.expectedTable {
			t.Errorf("Returned table was incorrect,\ngot:\n%swant:\n%s", buf.String(), test.expectedTable)
		}
	}
}
