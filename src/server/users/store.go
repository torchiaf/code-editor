package users

import (
	"context"
	"log"

	c "github.com/torchiaf/code-editor/server/config"
	k "github.com/torchiaf/code-editor/server/kube"
	"github.com/torchiaf/code-editor/server/models"
	"github.com/torchiaf/code-editor/server/utils"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type storeI interface {
	List()
	Get()
	Set()
	Del()
}

type store struct {
}

var _store = initStore()

func initStore() map[string]models.User {
	secret, err := k.Clientset.CoreV1().Secrets(c.Config.App.Namespace).Get(context.TODO(), c.Config.Resources.UsersName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	var users []models.User
	err = yaml.Unmarshal(secret.Data["users"], &users)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return utils.Map(users, func(user models.User) string { return user.Name })
}

func (store store) List() []models.User {
	return utils.Slice[models.User](_store)
}

func (store store) Get(username string) (models.User, bool) {
	return _store[username], _store[username] != models.User{}
}

func (store store) Set(user models.User) {
	secret, err := k.Clientset.CoreV1().Secrets(c.Config.App.Namespace).Get(context.TODO(), c.Config.Resources.UsersName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	var users []models.User
	err = yaml.Unmarshal(secret.Data["users"], &users)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	users = append(users, user)

	_byte, err := yaml.Marshal(users)
	if err != nil {
		panic(err)
	}

	secret.Data = make(map[string][]byte)
	secret.Data["users"] = _byte

	_, err = k.Clientset.CoreV1().Secrets(c.Config.App.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	_store[user.Name] = user
}

func (store store) Del(username string) {
	secret, err := k.Clientset.CoreV1().Secrets(c.Config.App.Namespace).Get(context.TODO(), c.Config.Resources.UsersName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	var users []models.User
	err = yaml.Unmarshal(secret.Data["users"], &users)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for i := range users {
		if username == users[i].Name {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}

	_byte, err := yaml.Marshal(users)
	if err != nil {
		panic(err)
	}

	secret.Data = make(map[string][]byte)
	secret.Data["users"] = _byte

	_, err = k.Clientset.CoreV1().Secrets(c.Config.App.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	delete(_store, username)
}

var Store = store{}
