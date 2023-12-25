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
	"server/utils"

	corev1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/remotecommand"

	traefikv1alpha1 "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func createService(name string, path string) {
	service := utils.ParseK8sResource[*v1.Service]("routes/service.yaml")

	service.Name = name
	service.Labels["app.kubernetes.io/name"] = name
	service.Spec.Selector["app.code-editor/path"] = path

	clientset.CoreV1().Services(c.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
}

func createRoute(name string, path string) string {
	cli, _ := client.New(restConfig, client.Options{})

	in := &unstructured.Unstructured{}
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Kind:    "IngressRoute",
		Version: "traefik.containo.us/v1alpha1",
	})

	err := cli.Get(context.Background(), client.ObjectKey{
		Namespace: c.Namespace,
		Name:      "code-editor-ui",
	}, in)
	if err != nil {
		e.FailOnError(err, "Failed to get Code-server IngressRoute")
	}

	ingressRoute := &traefikv1alpha1.IngressRoute{}

	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(in.UnstructuredContent(), &ingressRoute)

	route := utils.ParseFile[*traefikv1alpha1.Route]("routes/traefik-route.yaml")

	route.Match = fmt.Sprintf("Host(`localhost`) && PathPrefix(`/code-editor/%s/`)", path)
	route.Services[0].Name = name

	// var y interface{}
	// y = *route

	// var routes []interface{}

	// routes = append(routes, y)

	// if err := unstructured.SetNestedMap(in.Object, routes, "spec", "routes[0]"); err != nil {
	// 	e.FailOnError(err, "Failed to get Code-server IngressRoute")
	// }

	ingressRoute.Spec.Routes = append(ingressRoute.Spec.Routes, *route)

	routeUnstructured := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "traefik.containo.us/v1alpha1",
			"kind":       "Rule",
			"match":      fmt.Sprintf("Host(`localhost`) && PathPrefix(`/code-editor/%s/`)", path),
			"middlewares": []interface{}{
				map[string]interface{}{
					"name": "strip-prefix",
				},
			},
			"services": []interface{}{
				map[string]interface{}{
					"name": name,
					"port": "http",
				},
			},
		},
	}

	routes, _, _ := unstructured.NestedSlice(in.Object, "spec", "routes")

	routes = append(routes, routeUnstructured.Object)

	if err := unstructured.SetNestedSlice(in.Object, routes, "spec", "routes"); err != nil {
		e.FailOnError(err, "Failed to get Code-server IngressRoute")
	}

	// res, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(ingressRoute)

	// cli.Scheme().AddKnownTypeWithName(ingressRoute.GroupVersionKind(), ingressRoute)

	// traefikv1alpha1.SchemeBuilder.AddToScheme(cli.Scheme())

	err = cli.Update(context.TODO(), in)
	if err != nil {
		e.FailOnError(err, "Failed to get Code-server IngressRoute")
	}

	return ""
}

func createDeployment(user models.User, name string) models.CodeServerConfig {

	secret := utils.ParseK8sResource[*v1.Secret]("routes/secret.yaml")

	secret.Name = name
	secret.Namespace = c.Namespace

	password := utils.RandomString(20, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	secret.StringData = map[string]string{
		"PASSWORD": password,
	}

	clientset.CoreV1().Secrets(c.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})

	deployment := utils.ParseK8sResource[*corev1.Deployment]("routes/deployment.yaml")

	label := "app.code-editor/path"

	deployment.Name = name
	deployment.Labels[label] = user.Path
	deployment.Spec.Selector.MatchLabels[label] = user.Path
	deployment.Spec.Template.Labels[label] = user.Path

	deployment.Spec.Template.Spec.Containers[0].EnvFrom[0].SecretRef = &v1.SecretEnvSource{
		LocalObjectReference: v1.LocalObjectReference{
			Name: name,
		},
	}

	clientset.AppsV1().Deployments(c.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})

	return models.CodeServerConfig{
		Password: password,
	}
}

func ScaleCodeServer(user models.User, replicas int) (models.CodeServerConfig, error) {

	num := int32(replicas)
	name := fmt.Sprintf("%s-%s", c.App, user.Name)

	createService(name, user.Path)

	createRoute(name, user.Path)

	config := createDeployment(user, name)

	// deployment := clientset.AppsV1().Deployments(c.Namespace)

	// s, err := deployment.GetScale(context.TODO(), name, metav1.GetOptions{})
	// if err != nil {
	// 	e.FailOnError(err, "Failed to get Code-server Deployment/Scale")
	// }

	// sc := *s

	// if sc.Spec.Replicas == num {
	// 	log.Printf("Deployment is already in desired status")
	// 	return "", nil
	// }

	// sc.Spec.Replicas = num

	// _, err = deployment.UpdateScale(context.TODO(), name, &sc, metav1.UpdateOptions{})

	// if err != nil {
	// 	return "", err
	// }

	label := fmt.Sprintf("app.code-editor/path=%s", user.Path)

	if num > 0 {
		waitPodRunning(context.TODO(), label)
	}

	return config, nil
}
