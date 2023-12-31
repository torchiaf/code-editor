package users

import (
	"context"
	"encoding/json"
	"log"

	c "server/config"
	k "server/kube"
	"server/models"
	"server/utils"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type storeI interface {
	List()
	Get()
	Set()
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
	err = yaml.Unmarshal([]byte(secret.Data["users"]), &users)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return utils.Map(users, func(user models.User) string { return user.Name })
}

func (store store) List() []models.User {
	return utils.Slice[models.User](_store)
}

func (store store) Get(id string) (models.User, bool) {
	return _store[id], _store[id] != models.User{}
}

func (store store) Set(user models.User) {
	secret, err := k.Clientset.CoreV1().Secrets(c.Config.App.Namespace).Get(context.TODO(), c.Config.Resources.UsersName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	var users []models.User
	err = yaml.Unmarshal([]byte(secret.Data["users"]), &users)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	users = append(users, user)

	_byte, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}

	secret.Data["users"] = _byte

	_, err = k.Clientset.CoreV1().Secrets(c.Config.App.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	_store[user.Name] = user
}

var Store = store{}
