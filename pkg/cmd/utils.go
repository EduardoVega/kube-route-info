package cmd

import (
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// PodsToString returns a string of pod names separated by semicolons
func PodsToString(pods *v1.PodList) (podsString string) {
	podsString = ""
	arrayLength := len(pods.Items)

	for index, pod := range pods.Items {
		podsString += pod.Name

		if (index + 1) != arrayLength {
			podsString += ","
		}
	}

	return
}

// PortsToString returns a string of ports separated by semicolons
func PortsToString(ports []v1.ServicePort) (portsString string) {
	portsString = ""
	arrayLength := len(ports)

	for index, port := range ports {

		portsString += strconv.FormatInt(int64(port.Port), 10) + " "

		portsString += PortToString(port.TargetPort.Type, port.TargetPort.StrVal, port.TargetPort.IntVal)

		if port.NodePort != 0 {
			portsString += " " + strconv.FormatInt(int64(port.NodePort), 10)
		}

		if (index + 1) != arrayLength {
			portsString += ","
		}

	}

	return
}

// PortToString returns a string port
func PortToString(portType intstr.Type, strVal string, intVal int32) (portString string) {
	switch portType {
	case 0:
		portString = strconv.FormatInt(int64(intVal), 10)
	case 1:
		portString = strVal
	}

	return
}

// ServiceTypeToString returns a string service type
func ServiceTypeToString(serviceType v1.ServiceType) (serviceTypeString string) {
	switch serviceType {
	case v1.ServiceTypeLoadBalancer:
		serviceTypeString = "LoadBalancer"
	case v1.ServiceTypeNodePort:
		serviceTypeString = "NodePort"
	case v1.ServiceTypeClusterIP:
		serviceTypeString = "ClusterIP"
	case v1.ServiceTypeExternalName:
		serviceTypeString = "ExternalName"
	}

	return
}
