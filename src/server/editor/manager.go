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

	pods, err := clientset.CoreV1().Pods(c.App.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: label})
	if err != nil || len(pods.Items) == 0 {
		e.FailOnError(err, "Pod not found")
		return err
	}

	req := clientset.CoreV1().RESTClient().Post().Resource("pods").Name(pods.Items[0].Name).Namespace(c.App.Namespace).SubResource("exec")

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
		Container: c.App.Name,
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
	watcher, err := clientset.CoreV1().Pods(c.App.Namespace).Watch(context.TODO(), metav1.ListOptions{
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

type EditorI interface {
	Create()
	Config()
	Destroy()
	Login()
	Store()
}

type EditorConfigKeys struct {
	status   string
	path     string
	password string
}

type Editor struct {
	id        string
	name      string
	namespace string
	keys      EditorConfigKeys
}

func New(user models.User) Editor {

	id := fmt.Sprintf("%s-%s", c.App.Name, user.Id)

	return Editor{
		id:        id,
		name:      c.App.Name,
		namespace: c.App.Namespace,
		keys: EditorConfigKeys{
			status:   fmt.Sprintf("%s_STATUS", id),
			path:     fmt.Sprintf("%s_PATH", id),
			password: fmt.Sprintf("%s_PASSWORD", id),
		},
	}
}

func (editor Editor) Store() StoreData {
	return Store.Get(editor)
}

func (editor Editor) Login(port int32, password string) (models.CodeServerSession, error) {

	var session models.CodeServerSession

	// editor login endpoint
	loginUrl := ""

	if config.Config.IsDev {
		loginUrl = fmt.Sprintf("http://localhost/code-editor/%s/login", editor.Store().Path)
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

func (editor Editor) configsCreate() {
	data := StoreData{
		Status:   Enabled,
		Path:     utils.RandomString(13),
		Password: utils.RandomString(20, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	}

	Store.Set(editor, data)
}

func (editor Editor) configsDestroy() {
	Store.Del(editor)
}

func (editor Editor) serviceCreate() *v1.Service {

	service := utils.ParseK8sResource[*v1.Service]("assets/templates/service.yaml")

	service.Name = editor.id
	service.Labels[NAME_LABEL] = editor.id
	service.Labels[INSTANCE_LABEL] = editor.name
	service.Spec.Selector["app.code-editor/path"] = editor.Store().Path

	ret, _ := clientset.CoreV1().Services(editor.namespace).Create(context.TODO(), service, metav1.CreateOptions{})

	return ret
}

func (editor Editor) serviceDestroy() {
	clientset.CoreV1().Services(editor.namespace).Delete(context.TODO(), editor.id, metav1.DeleteOptions{})
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
		Name:      c.Resources.IngressName,
	}, in)
	if err != nil {
		e.FailOnError(err, "Failed to get Editor IngressRoute")
	}

	routeUnstructured := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "traefik.containo.us/v1alpha1",
			"kind":       "Rule",
			"match":      fmt.Sprintf("Host(`localhost`) && PathPrefix(`/code-editor/%s/`)", editor.Store().Path),
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

func (editor Editor) ruleDelete() {
	cli, _ := client.New(restConfig, client.Options{})

	in := &unstructured.Unstructured{}
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Kind:    "IngressRoute",
		Version: "traefik.containo.us/v1alpha1",
	})

	err := cli.Get(context.Background(), client.ObjectKey{
		Namespace: editor.namespace,
		Name:      c.Resources.IngressName,
	}, in)
	if err != nil {
		e.FailOnError(err, "Failed to get Editor IngressRoute")
	}

	routes, _, _ := unstructured.NestedSlice(in.Object, "spec", "routes")

	for i := range routes {

		services, _, _ := unstructured.NestedSlice(routes[i].(map[string]interface{}), "services")

		name, _, _ := unstructured.NestedString(services[0].(map[string]interface{}), "name")

		if name == editor.id {
			routes = append(routes[:i], routes[i+1:]...)
			break
		}

	}

	if err := unstructured.SetNestedSlice(in.Object, routes, "spec", "routes"); err != nil {
		e.FailOnError(err, "Failed to set Editor rules")
	}

	err = cli.Update(context.TODO(), in)
	if err != nil {
		e.FailOnError(err, "Failed to get Editor IngressRoute")
	}
}

func (editor Editor) deploymentCreate(service *v1.Service) *corev1.Deployment {

	deployment := utils.ParseK8sResource[*corev1.Deployment]("assets/templates/deployment.yaml")

	deployment.Name = editor.id

	deployment.Labels[MATCH_LABEL] = editor.Store().Path
	deployment.Spec.Selector.MatchLabels[MATCH_LABEL] = editor.Store().Path
	deployment.Spec.Template.Labels[MATCH_LABEL] = editor.Store().Path

	deployment.Labels[INSTANCE_LABEL] = editor.name
	deployment.Spec.Selector.MatchLabels[INSTANCE_LABEL] = editor.name
	deployment.Spec.Template.Labels[INSTANCE_LABEL] = editor.name

	deployment.Spec.Template.Spec.Containers[0].Name = c.App.Name
	deployment.Spec.Template.Spec.ServiceAccountName = editor.name

	deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.SecretKeyRef.Name = c.Resources.ConfigName
	deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.SecretKeyRef.Key = editor.keys.password

	ret, _ := clientset.AppsV1().Deployments(editor.namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})

	return ret
}

func (editor Editor) deploymentDestroy() {
	clientset.AppsV1().Deployments(editor.namespace).Delete(context.TODO(), editor.id, metav1.DeleteOptions{})
}

func (editor Editor) Create() (int32, error) {

	editor.configsCreate()

	service := editor.serviceCreate()

	editor.ruleCreate()

	editor.deploymentCreate(service)

	label := matchLabel(editor.Store().Path)

	waitPodRunning(context.TODO(), label)

	port := service.Spec.Ports[0].Port

	return port, nil
}

func (editor Editor) Config(gitCmd string) error {
	label := matchLabel(editor.Store().Path)

	execCmdOnPod(label, gitCmd, nil, nil, nil)

	return nil
}

func (editor Editor) Destroy(user models.User) error {

	editor.deploymentDestroy()

	editor.serviceDestroy()

	editor.ruleDelete()

	editor.configsDestroy()

	return nil
}
