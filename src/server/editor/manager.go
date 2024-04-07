package editor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	config "github.com/torchiaf/code-editor/server/config"
	e "github.com/torchiaf/code-editor/server/error"
	k "github.com/torchiaf/code-editor/server/kube"
	"github.com/torchiaf/code-editor/server/models"
	"github.com/torchiaf/code-editor/server/utils"

	corev1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/remotecommand"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
	log.Println(a)
	return len(p), nil
}

func execCmdOnPod(
	ctx context.Context,
	label string,
	command string,
	// stdin io.Reader,
	// stdout io.Writer,
	// stderr io.Writer,
) error {
	cmd := []string{
		"sh",
		"-c",
		command,
	}

	pods, err := k.Clientset.CoreV1().Pods(c.App.Namespace).List(ctx, metav1.ListOptions{LabelSelector: label})
	if err != nil || len(pods.Items) == 0 {
		e.FailOnError(err, "Pod not found")
		return err
	}

	req := k.Clientset.CoreV1().RESTClient().Post().Resource("pods").Name(pods.Items[0].Name).Namespace(c.App.Namespace).SubResource("exec")

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
	exec, err := remotecommand.NewSPDYExecutor(k.RestConfig, "POST", url)
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
	watcher, err := k.Clientset.CoreV1().Pods(c.App.Namespace).Watch(ctx, metav1.ListOptions{
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
	status         string
	path           string
	password       string
	vscodeSettings string
}

type Editor struct {
	ctx       context.Context
	id        string
	name      string
	namespace string
	keys      EditorConfigKeys
}

func New(ctx context.Context, user models.User) Editor {

	id := fmt.Sprintf("%s-%s", c.App.Name, user.Id)

	return Editor{
		ctx:       ctx,
		id:        id,
		name:      c.App.Name,
		namespace: c.App.Namespace,
		keys: EditorConfigKeys{
			status:         fmt.Sprintf("%s_STATUS", id),
			path:           fmt.Sprintf("%s_PATH", id),
			password:       fmt.Sprintf("%s_PASSWORD", id),
			vscodeSettings: fmt.Sprintf("%s_VSCODE_SETTINGS", id),
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

	if err != nil || resp.StatusCode != 302 {
		return session, errors.New("Editor, code-server login error")
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

func (editor Editor) configsCreate(enableConfig models.EnableConfig) {

	defaultSettings := utils.ReadFile("assets/templates/vscode-settings.json")

	var settingsMap map[string]interface{}

	// Merge VSCode settings
	if err := json.Unmarshal(defaultSettings, &settingsMap); err != nil {
		panic(err)
	}

	maps.Copy(settingsMap, enableConfig.VscodeSettings)
	for _, extension := range enableConfig.Extensions {
		maps.Copy(settingsMap, extension.Settings)
	}

	vscodeSettings, err := json.Marshal(settingsMap)
	if err != nil {
		panic(err)
	}

	data := map[string][]byte{
		editor.keys.status:         []byte(Enabled),
		editor.keys.path:           []byte(utils.RandomString(13)),
		editor.keys.password:       []byte(utils.RandomString(20, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")),
		editor.keys.vscodeSettings: vscodeSettings,
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

	ret, _ := k.Clientset.CoreV1().Services(editor.namespace).Create(editor.ctx, service, metav1.CreateOptions{})

	return ret
}

func (editor Editor) serviceDestroy() {
	k.Clientset.CoreV1().Services(editor.namespace).Delete(editor.ctx, editor.id, metav1.DeleteOptions{})
}

func (editor Editor) ruleCreate() error {
	cli, _ := client.New(k.RestConfig, client.Options{})

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

	err = cli.Update(editor.ctx, in)
	if err != nil {
		e.FailOnError(err, "Failed to get Editor IngressRoute")
	}

	return nil
}

func (editor Editor) ruleDelete() {
	cli, _ := client.New(k.RestConfig, client.Options{})

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

	err = cli.Update(editor.ctx, in)
	if err != nil {
		e.FailOnError(err, "Failed to get Editor IngressRoute")
	}
}

func (editor Editor) deploymentCreate(cfg models.EnableConfig) *corev1.Deployment {

	deployment := utils.ParseK8sResource[*corev1.Deployment]("assets/templates/deployment.yaml")

	deployment.Name = editor.id

	// Labels
	deployment.Labels[MATCH_LABEL] = editor.Store().Path
	deployment.Spec.Selector.MatchLabels[MATCH_LABEL] = editor.Store().Path
	deployment.Spec.Template.Labels[MATCH_LABEL] = editor.Store().Path

	deployment.Labels[INSTANCE_LABEL] = editor.name
	deployment.Spec.Selector.MatchLabels[INSTANCE_LABEL] = editor.name
	deployment.Spec.Template.Labels[INSTANCE_LABEL] = editor.name

	// Containers
	deployment.Spec.Template.Spec.Containers[0].Name = c.App.Name
	deployment.Spec.Template.Spec.ServiceAccountName = editor.name

	// Password
	deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.SecretKeyRef.Name = c.Resources.ConfigName
	deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.SecretKeyRef.Key = editor.keys.password

	// VSCode default settings volume
	for i, volume := range deployment.Spec.Template.Spec.Volumes {
		if volume.Name == "vscode-settings" {
			deployment.Spec.Template.Spec.Volumes[i].Secret.SecretName = c.Resources.ConfigName
			deployment.Spec.Template.Spec.Volumes[i].Secret.Items[0].Key = editor.keys.vscodeSettings
			break
		}
	}

	initContainers := utils.ParseJsonFile[map[string]v1.Container]("assets/templates/containers.json")

	// TODO should comes from API body
	deployment.Spec.Template.Spec.InitContainers = append(deployment.Spec.Template.Spec.InitContainers, initContainers["ssh"])

	if (cfg.Git != models.GitConfig{}) {
		ic := initContainers["gitconfig"]
		ic.Command = []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("chown -R 1000:1000 /etc && echo '[user]\n\temail = %s\n\tname = %s' > /home/coder/.gitconfig", cfg.Git.Email, cfg.Git.Name),
		}
		deployment.Spec.Template.Spec.InitContainers = append(deployment.Spec.Template.Spec.InitContainers, ic)
	}

	if len(cfg.Extensions) > 0 {
		ic := initContainers["extensions"]

		extensionCmd := ""
		for i, extension := range cfg.Extensions {
			extensionCmd += fmt.Sprintf("code-server --install-extension %s", extension.Id)
			if i < len(cfg.Extensions)-1 {
				extensionCmd += " && "
			}
		}
		ic.Command = []string{
			"/bin/sh",
			"-c",
			extensionCmd,
		}
		deployment.Spec.Template.Spec.InitContainers = append(deployment.Spec.Template.Spec.InitContainers, ic)
	}

	ret, _ := k.Clientset.AppsV1().Deployments(editor.namespace).Create(ctx, deployment, metav1.CreateOptions{})

	return ret
}

func (editor Editor) deploymentDestroy() {
	k.Clientset.AppsV1().Deployments(editor.namespace).Delete(editor.ctx, editor.id, metav1.DeleteOptions{})
}

func (editor Editor) Create(enableConfig models.EnableConfig) (int32, error) {

	editor.configsCreate(enableConfig)

	service := editor.serviceCreate()

	editor.ruleCreate()

	editor.deploymentCreate(enableConfig)

	label := matchLabel(editor.Store().Path)

	waitPodRunning(editor.ctx, label)

	port := service.Spec.Ports[0].Port

	return port, nil
}

func (editor Editor) Config(gitCmd string) error {
	label := matchLabel(editor.Store().Path)

	execCmdOnPod(editor.ctx, label, gitCmd)

	return nil
}

func (editor Editor) Destroy(user models.User) error {

	editor.deploymentDestroy()

	editor.serviceDestroy()

	editor.ruleDelete()

	editor.configsDestroy()

	return nil
}
