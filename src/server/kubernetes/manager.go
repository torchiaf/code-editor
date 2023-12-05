package kubernetes

import (
	"context"

	config "server/config"
	e "server/error"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var clientset = InitKubeconfig()
var c = config.Config

func scaleDeployment(name string, namespace string, scale int32) (string, error) {
	s, err := clientset.AppsV1().Deployments(namespace).GetScale(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		e.FailOnError(err, "Failed to get Deployment/Scale")
	}

	sc := *s
	sc.Spec.Replicas = scale

	res, err := clientset.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), name, &sc, metav1.UpdateOptions{})

	if err != nil {
		return "", err
	}

	return res.Name, nil
}

func StartEditor() (string, error) {
	return scaleDeployment(c.App, c.Namespace, 1)
}

func StopEditor() (string, error) {
	return scaleDeployment(c.App, c.Namespace, 0)
}
