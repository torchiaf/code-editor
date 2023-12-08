package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	config "server/config"
	e "server/error"
	"server/models"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
)

var clientset, restConfig = InitKubeconfig()
var c = config.Config

type LogStreamer struct {
	b bytes.Buffer
}

func (l *LogStreamer) String() string {
	return l.b.String()
}

func (l *LogStreamer) Write(p []byte) (n int, err error) {
	a := strings.TrimSpace(string(p))
	l.b.WriteString(a)
	log.Printf(a)
	return len(p), nil
}

func ExecCmdOnPod(label string, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	cmd := []string{
		"sh",
		"-c",
		command,
	}

	pods, err := clientset.CoreV1().Pods(c.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: label})
	if err != nil || len(pods.Items) == 0 {
		e.FailOnError(err, "Pod not found")
		return err
	}

	req := clientset.CoreV1().RESTClient().Post().Resource("pods").Name(pods.Items[0].Name).Namespace(c.Namespace).SubResource("exec")

	scheme := runtime.NewScheme()
	if err := v1.AddToScheme(scheme); err != nil {
		return err
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&v1.PodExecOptions{
		Stdin:     true,
		Stdout:    true,
		Stderr:    false,
		TTY:       true,
		Container: c.App,
		Command:   cmd,
	}, parameterCodec)
	url := req.URL()
	exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", url)
	if err != nil {
		return err
	}
	var streamErr error
	l := &LogStreamer{}

	streamErr = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: l,
		Stderr: nil,
		Tty:    true,
	})

	if streamErr != nil {
		return streamErr
	}

	return nil
}

func waitPodRunning(ctx context.Context, label string) error {
	watcher, err := clientset.CoreV1().Pods(c.Namespace).Watch(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
	})

	if err != nil {
		return err
	}

	defer watcher.Stop()

	for {
		select {
		case event := <-watcher.ResultChan():
			pod := event.Object.(*v1.Pod)

			if pod.Status.Phase == v1.PodRunning {
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// func (k *k8sClient) waitPodDeleted(ctx context.Context, resName string) error {
// 	watcher, err := k.createPodWatcher(ctx, resName)
// 	if err != nil {
// 		return err
// 	}

// 	defer watcher.Stop()

// 	for {
// 		select {
// 		case event := <-watcher.ResultChan():

// 			if event.Type == watch.Deleted {
// 				k.logger.Debugf("The POD \"%s\" is deleted", resName)

// 				return nil
// 			}

// 		case <-ctx.Done():
// 			k.logger.Debugf("Exit from waitPodDeleted for POD \"%s\" because the context is done", resName)
// 			return nil
// 		}
// 	}
// }

func ScaleCodeServer(user models.User, replicas int) (string, error) {

	num := int32(replicas)
	namespace := c.Namespace
	name := fmt.Sprintf("%s-%s", c.App, user.Name)

	deployment := clientset.AppsV1().Deployments(namespace)

	s, err := deployment.GetScale(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		e.FailOnError(err, "Failed to get Code-server Deployment/Scale")
	}

	sc := *s

	if sc.Spec.Replicas == num {
		log.Printf("Deployment is already in desired status")
		return "", nil
	}

	sc.Spec.Replicas = num

	_, err = deployment.UpdateScale(context.TODO(), name, &sc, metav1.UpdateOptions{})

	if err != nil {
		return "", err
	}

	label := fmt.Sprintf("app.code-editor/path=%s", user.Path)

	if num > 0 {
		waitPodRunning(context.TODO(), label)
	}

	return "", nil
}
