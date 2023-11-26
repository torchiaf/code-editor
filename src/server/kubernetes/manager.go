package kubernetes

import (
	"context"

	"server/config"
	"server/utils"

	e "server/error"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var clientset = InitKubeconfig()
var c = config.GetConfig()

func scaleDeployment(name string, namespace string, scale int32) string {
	s, err := clientset.AppsV1().Deployments(namespace).GetScale(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		e.FailOnError(err, "Failed to get Deployment/Scale")
	}

	sc := *s
	sc.Spec.Replicas = scale

	res, err := clientset.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), name, &sc, metav1.UpdateOptions{})

	if err != nil {
		panic(err.Error())
	}

	return res.Name
}

func StartEditor() string {
	return scaleDeployment(c.App, c.Namespace, 1)
}

func StopEditor() string {
	return scaleDeployment(c.App, c.Namespace, 0)
}

func GetPods() string {
	pods, err := clientset.CoreV1().Pods(c.App).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return utils.ToString(pods)
}
