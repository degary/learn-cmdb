package k8s

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Cli *kubernetes.Clientset
}

func (c *Client) Deployments() ([]*appsv1.Deployment, error) {
	namespaces, err := c.Namespaces()
	deployments := []*appsv1.Deployment{}
	if err != nil {
		return nil, err
	}
	for _, ns := range namespaces {
		deployCli := c.Cli.AppsV1().Deployments(ns.Name)
		deployList, err := deployCli.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		//deployments = append(deployments,deployList.Items...)
		for _, deploy := range deployList.Items {
			var m appsv1.Deployment
			m = deploy
			deployments = append(deployments, &m)
		}
	}
	return deployments, nil
}

func (c *Client) Namespaces() ([]corev1.Namespace, error) {
	nsCli := c.Cli.CoreV1().Namespaces()
	namespaces := []corev1.Namespace{}
	nsList, err := nsCli.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	namespaces = append(namespaces, nsList.Items...)
	return namespaces, nil
}

func (c *Client) Services() ([]*corev1.Service, error) {
	services := []*corev1.Service{}
	namespaces, err := c.Namespaces()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, ns := range namespaces {
		serviceCli := c.Cli.CoreV1().Services(ns.Name)
		serviceList, err := serviceCli.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		for _,sv := range serviceList.Items{
			var svPtr corev1.Service
			svPtr = sv
			services = append(services,&svPtr)
		}
	}
	return services, nil
}

func NewClient(configPath string) *Client {
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil
	}

	//创建k8s client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil
	}
	return &Client{
		Cli: clientset,
	}
}
