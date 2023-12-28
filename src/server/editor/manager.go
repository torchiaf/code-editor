package editor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var clientset, restConfig = initKubeconfig()
var c = config.Config

type logStreamer struct {
	b bytes.Buffer
}

func (l *logStreamer) String() string {
	return l.b.String()
}

func (l *logStreamer) Write(p []byte) (n int, err error) {
	a := strings.TrimSpace(string(p))
	l.b.WriteString(a)
	log.Printf(a)
	return len(p), nil
}

func execCmdOnPod(label string, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
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
	l := &logStreamer{}

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

func (editor Editor) Login(port int32, password string) (models.CodeServerSession, error) {

	var session models.CodeServerSession

	// editor login endpoint
	loginUrl := ""

	if config.Config.IsDev {
		loginUrl = fmt.Sprintf("http://localhost/code-editor/%s/login", editor.path)
	} else {
		loginUrl = fmt.Sprintf("http://%s:%d/login", editor.id, port)
	}

	// JSON body
	data := url.Values{}
	data.Set("password", password)

	// Create a HTTP post request
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return session, errors.New("Editor, login request creation error")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Host", "localhost")

	resp, err := client.Do(req)

	if err != nil {
		return session, errors.New("Editor, login response error")
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()

	if len(cookies) == 0 {
		return session, errors.New("Login failed, invalid User or Password")
	}

	cookie := cookies[0]

	session.Name = cookie.Name
	session.Value = cookie.Value

	return session, nil
}

func deleteDeployment(user models.User, name string) error {
	clientset.AppsV1().Deployments(c.Namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

	return nil
}

type EditorI interface {
	Create()
	Config()
	Destroy()
	Login()
	credentialsCreate()
	serviceCreate()
	ruleCreate()
	deploymentCreate()
}

type Editor struct {
	id        string
	namespace string
	label     string
	path      string
}

func (editor Editor) credentialsCreate() *v1.Secret {
	secret := utils.ParseK8sResource[*v1.Secret]("assets/templates/secret.yaml")

	secret.Name = editor.id
	secret.Namespace = editor.namespace

	secret.StringData = map[string]string{
		"PASSWORD": utils.RandomString(20, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	}

	clientset.CoreV1().Secrets(c.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})

	return secret
}

func (editor Editor) serviceCreate() *v1.Service {

	service := utils.ParseK8sResource[*v1.Service]("assets/templates/service.yaml")

	service.Name = editor.id
	service.Labels["app.kubernetes.io/name"] = editor.id
	service.Spec.Selector["app.code-editor/path"] = editor.path

	clientset.CoreV1().Services(editor.namespace).Create(context.TODO(), service, metav1.CreateOptions{})

	return service
}

func (editor Editor) ruleCreate() error {
	cli, _ := client.New(restConfig, client.Options{})

	in := &unstructured.Unstructured{}
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Kind:    "IngressRoute",
		Version: "traefik.containo.us/v1alpha1",
	})

	err := cli.Get(context.Background(), client.ObjectKey{
		Namespace: editor.namespace,
		Name:      "code-editor-ui",
	}, in)
	if err != nil {
		e.FailOnError(err, "Failed to get Editor IngressRoute")
	}

	routeUnstructured := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "traefik.containo.us/v1alpha1",
			"kind":       "Rule",
			"match":      fmt.Sprintf("Host(`localhost`) && PathPrefix(`/code-editor/%s/`)", editor.path),
			"middlewares": []interface{}{
				map[string]interface{}{
					"name": "strip-prefix",
				},
			},
			"services": []interface{}{
				map[string]interface{}{
					"name": editor.id,
					"port": "http",
				},
			},
		},
	}

	routes, _, _ := unstructured.NestedSlice(in.Object, "spec", "routes")

	routes = append(routes, routeUnstructured.Object)

	if err := unstructured.SetNestedSlice(in.Object, routes, "spec", "routes"); err != nil {
		e.FailOnError(err, "Failed to set Editor rules")
	}

	err = cli.Update(context.TODO(), in)
	if err != nil {
		e.FailOnError(err, "Failed to get Editor IngressRoute")
	}

	return nil
}

func (editor Editor) deploymentCreate(service *v1.Service) *corev1.Deployment {

	deployment := utils.ParseK8sResource[*corev1.Deployment]("assets/templates/deployment.yaml")

	label := "app.code-editor/path"

	deployment.Name = editor.id
	deployment.Labels[label] = editor.path
	deployment.Spec.Selector.MatchLabels[label] = editor.path
	deployment.Spec.Template.Labels[label] = editor.path

	deployment.Spec.Template.Spec.Containers[0].EnvFrom[0].SecretRef = &v1.SecretEnvSource{
		LocalObjectReference: v1.LocalObjectReference{
			Name: editor.id,
		},
	}

	clientset.AppsV1().Deployments(c.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})

	return deployment
}

func New(user models.User) Editor {
	return Editor{
		id:        fmt.Sprintf("%s-%s", c.App, user.Name),
		namespace: c.Namespace,
		label:     fmt.Sprintf("app.code-editor/path=%s", user.Path),
		path:      user.Path,
	}
}

func (editor Editor) Create() (int32, string, error) {

	service := editor.serviceCreate()

	editor.ruleCreate()

	secret := editor.credentialsCreate()

	editor.deploymentCreate(service)

	waitPodRunning(context.TODO(), editor.label)

	port := service.Spec.Ports[0].Port
	password := secret.StringData["PASSWORD"]

	return port, password, nil
}

func (editor Editor) Config(gitCmd string) error {
	execCmdOnPod(editor.label, gitCmd, nil, nil, nil)

	return nil
}

func (editor Editor) Destroy(user models.User) error {
	name := fmt.Sprintf("%s-%s", c.App, user.Name)

	return deleteDeployment(user, name)
}
