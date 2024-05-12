package editor

import (
	"flag"
	"path/filepath"

	c "github.com/torchiaf/code-editor/server/config"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func initKubeconfig() (*kubernetes.Clientset, client.Client, *rest.Config) {
	var config *rest.Config
	var err error

	if c.Config.IsDev {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err)
		}
	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	runtimeClient, err := client.New(config, client.Options{})
	if err != nil {
		panic(err.Error())
	}

	return clientset, runtimeClient, config
}

var Clientset, RuntimeClient, RestConfig = initKubeconfig()
