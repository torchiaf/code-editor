package editor

import (
	"context"
	"fmt"

	k "github.com/torchiaf/code-editor/server/kube"

	"github.com/torchiaf/code-editor/server/users"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StoreData struct {
	Status         string
	Path           string
	Password       string
	VScodeSettings string
}

type storeI interface {
	Get()
	Set()
	Del()
}

type store struct {
}

var _store = initStore()

func initStore() map[string]StoreData {

	store := map[string]StoreData{}

	secret, err := k.Clientset.CoreV1().Secrets(c.App.Namespace).Get(context.TODO(), c.Resources.ConfigName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	for _, user := range users.Store.List() {
		id := fmt.Sprintf("%s-%s", c.App.Name, user.Id)

		if len(secret.Data[fmt.Sprintf("%s_STATUS", id)]) > 0 {
			dataStore := StoreData{
				Status:         string(secret.Data[fmt.Sprintf("%s_STATUS", id)]),
				Path:           string(secret.Data[fmt.Sprintf("%s_PATH", id)]),
				Password:       string(secret.Data[fmt.Sprintf("%s_PASSWORD", id)]),
				VScodeSettings: string(secret.Data[fmt.Sprintf("%s_VSCODE_SETTINGS", id)]),
			}

			store[id] = dataStore
		}
	}

	return store
}

func (store store) Get(editor Editor) StoreData {
	return _store[editor.id]
}

func (store store) Set(editor Editor, data map[string][]byte) {
	secret, err := k.Clientset.CoreV1().Secrets(editor.namespace).Get(context.TODO(), c.Resources.ConfigName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	secret.Data = data

	_, err = k.Clientset.CoreV1().Secrets(editor.namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	dataStore := StoreData{
		Status:         string(data[editor.keys.status]),
		Path:           string(data[editor.keys.path]),
		Password:       string(data[editor.keys.password]),
		VScodeSettings: string(data[editor.keys.vscodeSettings]),
	}

	_store[editor.id] = dataStore
}

func (store store) Del(editor Editor) {
	secret, err := k.Clientset.CoreV1().Secrets(editor.namespace).Get(context.TODO(), c.Resources.ConfigName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	delete(secret.Data, editor.keys.status)
	delete(secret.Data, editor.keys.path)
	delete(secret.Data, editor.keys.password)
	delete(secret.Data, editor.keys.vscodeSettings)

	_, err = k.Clientset.CoreV1().Secrets(editor.namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	delete(_store, editor.id)
}

var Store = store{}
