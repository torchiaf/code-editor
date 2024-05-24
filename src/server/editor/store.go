package editor

import (
	"context"
	"fmt"
	"maps"

	k "github.com/torchiaf/code-editor/server/kube"

	"github.com/torchiaf/code-editor/server/users"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ctx = context.Background()

type StoreData struct {
	ViewName       string
	Status         string
	Path           string
	Query          string
	Password       string
	VScodeSettings string
	GitAuth        string
	Session        string
	RepoType       string
	Repo           string
}

type store struct {
}

var _store = initStore()

func initStore() map[string]StoreData {

	store := map[string]StoreData{}

	secret, err := k.Clientset.CoreV1().Secrets(c.App.Namespace).Get(ctx, c.Resources.ConfigName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	for _, user := range users.Store.List() {
		id := fmt.Sprintf("%s-%s", c.App.Name, user.Id)

		if len(secret.Data[fmt.Sprintf("%s_STATUS", id)]) > 0 {
			dataStore := StoreData{
				ViewName:       string(secret.Data[fmt.Sprintf("%s_VIEW_NAME", id)]),
				Status:         string(secret.Data[fmt.Sprintf("%s_STATUS", id)]),
				Path:           string(secret.Data[fmt.Sprintf("%s_PATH", id)]),
				Query:          string(secret.Data[fmt.Sprintf("%s_QUERY", id)]),
				Password:       string(secret.Data[fmt.Sprintf("%s_PASSWORD", id)]),
				VScodeSettings: string(secret.Data[fmt.Sprintf("%s_VSCODE_SETTINGS", id)]),
				GitAuth:        string(secret.Data[fmt.Sprintf("%s_GIT_AUTH", id)]),
				Session:        string(secret.Data[fmt.Sprintf("%s_SESSION", id)]),
				RepoType:       string(secret.Data[fmt.Sprintf("%s_REPO_TYPE", id)]),
				Repo:           string(secret.Data[fmt.Sprintf("%s_REPO", id)]),
			}

			store[id] = dataStore
		}
	}

	return store
}

func (store store) Get(editor Editor) StoreData {
	return _store[editor.Id]
}

func (store store) Set(editor Editor, data map[string][]byte) {
	secret, err := k.Clientset.CoreV1().Secrets(editor.namespace).Get(ctx, c.Resources.ConfigName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	if secret.Data == nil {
		secret.Data = data
	} else {
		maps.Copy(secret.Data, data)
	}

	_, err = k.Clientset.CoreV1().Secrets(editor.namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	dataStore := StoreData{
		ViewName:       string(data[editor.keys.viewName]),
		Status:         string(data[editor.keys.status]),
		Path:           string(data[editor.keys.path]),
		Query:          string(data[editor.keys.query]),
		Password:       string(data[editor.keys.password]),
		VScodeSettings: string(data[editor.keys.vscodeSettings]),
		GitAuth:        string(data[editor.keys.gitAuth]),
		Session:        string(data[editor.keys.session]),
		RepoType:       string(data[editor.keys.repoType]),
		Repo:           string(data[editor.keys.repo]),
	}

	_store[editor.Id] = dataStore
}

func (store store) Upd(editor Editor, session string, repoType string, repo string, query string) {

	data := store.Get(editor)
	m := make(map[string][]byte)

	m[editor.keys.viewName] = []byte(data.ViewName)
	m[editor.keys.status] = []byte(data.Status)
	m[editor.keys.path] = []byte(data.Path)
	m[editor.keys.password] = []byte(data.Password)
	m[editor.keys.vscodeSettings] = []byte(data.VScodeSettings)
	m[editor.keys.gitAuth] = []byte(data.GitAuth)

	sessionValue := data.Session
	if session != "" {
		sessionValue = session
	}

	repoTypeValue := data.RepoType
	if repoType != "" {
		repoTypeValue = repoType
	}

	repoValue := data.Repo
	if repo != "" {
		repoValue = repo
	}

	queryValue := data.Query
	if query != "" {
		queryValue = query
	}

	m[editor.keys.session] = []byte(sessionValue)
	m[editor.keys.repoType] = []byte(repoTypeValue)
	m[editor.keys.repo] = []byte(repoValue)
	m[editor.keys.query] = []byte(queryValue)

	store.Set(editor, m)
}

func (store store) Del(editor Editor) {
	secret, err := k.Clientset.CoreV1().Secrets(editor.namespace).Get(ctx, c.Resources.ConfigName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	delete(secret.Data, editor.keys.viewName)
	delete(secret.Data, editor.keys.status)
	delete(secret.Data, editor.keys.path)
	delete(secret.Data, editor.keys.query)
	delete(secret.Data, editor.keys.password)
	delete(secret.Data, editor.keys.vscodeSettings)
	delete(secret.Data, editor.keys.gitAuth)
	delete(secret.Data, editor.keys.session)
	delete(secret.Data, editor.keys.repoType)
	delete(secret.Data, editor.keys.repo)

	_, err = k.Clientset.CoreV1().Secrets(editor.namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	delete(_store, editor.Id)
}

var Store = store{}
