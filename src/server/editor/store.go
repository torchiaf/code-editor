package editor

import (
	"context"
	"fmt"

	"server/config"
	"server/models"
	"server/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StoreData struct {
	Status   string
	Path     string
	Password string
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

	secret, err := clientset.CoreV1().Secrets(config.Config.App.Namespace).Get(context.TODO(), c.Resources.ConfigName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	users := utils.Slice[models.User](c.Users)

	for _, user := range users {
		id := fmt.Sprintf("%s-%s", c.App.Name, user.Id)

		if len(secret.Data[fmt.Sprintf("%s_STATUS", id)]) > 0 {
			data := StoreData{
				Status:   string(secret.Data[fmt.Sprintf("%s_STATUS", id)]),
				Path:     string(secret.Data[fmt.Sprintf("%s_PATH", id)]),
				Password: string(secret.Data[fmt.Sprintf("%s_PASSWORD", id)]),
			}

			store[id] = data
		}
	}

	return store
}

func (store store) Get(editor Editor) StoreData {
	return _store[editor.id]
}

func (store store) Set(editor Editor, data StoreData) {
	secret, err := clientset.CoreV1().Secrets(editor.namespace).Get(context.TODO(), c.Resources.ConfigName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	secret.StringData = map[string]string{
		editor.keys.status:   data.Status,
		editor.keys.path:     data.Path,
		editor.keys.password: data.Password,
	}

	_, err = clientset.CoreV1().Secrets(editor.namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	_store[editor.id] = data
}

func (store store) Del(editor Editor) {
	secret, err := clientset.CoreV1().Secrets(editor.namespace).Get(context.TODO(), c.Resources.ConfigName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	delete(secret.Data, editor.keys.status)
	delete(secret.Data, editor.keys.path)
	delete(secret.Data, editor.keys.password)

	_, err = clientset.CoreV1().Secrets(editor.namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	delete(_store, editor.id)
}

var Store = store{}
