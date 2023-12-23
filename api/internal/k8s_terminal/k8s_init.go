package k8s_terminal

import (
	"errors"
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sConfig struct {
	restConfig   *rest.Config
	clientConfig clientcmd.ClientConfig
	clientset    *kubernetes.Clientset
}

func K8sInit() (*kubernetes.Clientset, error) {
	var k8sConfig K8sConfig
	initKubeClientConfig(&k8sConfig)
	err := initKubeClient(&k8sConfig)
	if err != nil {
		return nil, err
	}
	return k8sConfig.clientset, nil
}

func initKubeClientConfig(k8sConfig *K8sConfig) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	overrides := clientcmd.ConfigOverrides{}

	k8sConfig.clientConfig = clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, &overrides, os.Stdin)
}

func initKubeClient(k8sConfig *K8sConfig) (err error) {
	if k8sConfig.clientset != nil {
		return errors.New("fail to k8sConfig.clientset init")
	}
	k8sConfig.restConfig, err = k8sConfig.clientConfig.ClientConfig()
	if err != nil {
		log.Fatal(err)
		return err
	}

	k8sConfig.clientset, err = kubernetes.NewForConfig(k8sConfig.restConfig)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func ReturnK8sConfig() (*K8sConfig, error) {
	var k8sConfig K8sConfig
	initKubeClientConfig(&k8sConfig)
	err := initKubeClient(&k8sConfig)
	if err != nil {
		return nil, err
	}
	return &k8sConfig, nil
}
