package kubernetes

import (
	"context"

	"server/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var clientset = InitKubeconfig()

func CreateEditor() string {
	return ""
}

func GetPods() string {
	pods, err := clientset.CoreV1().Pods("code-editor").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return utils.ToString(pods)
}
