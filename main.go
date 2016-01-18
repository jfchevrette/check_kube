package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"

	client "k8s.io/kubernetes/pkg/client/unversioned"
)

const (
	appVersion = "0.1.1"

	// Nagios status codes
	nagiosStatusOK       = 0
	nagiosStatusWarning  = 1
	nagiosStatusCritical = 2
	nagiosStatusUnknown  = 3
)

var (
	kubeConfig = &client.Config{}
)

func checkKubeNodes(c *cli.Context) {
	var err error
	var statusCode = nagiosStatusOK
	var statusLine []string

	kubeClient, err := client.New(kubeConfig)
	if err != nil {
		fmt.Printf("CRITICAL: %s\n", err)
		os.Exit(nagiosStatusUnknown)
	}

	// TODO: Allow for selecting nodes based on labels or fields
	nodes, err := kubeClient.Nodes().List(labels.Everything(), fields.Everything())
	if err != nil {
		fmt.Printf("CRITICAL: %s\n", err)
		os.Exit(nagiosStatusUnknown)
	}

	// Loop over all the nodes
	for _, node := range nodes.Items {

		// Loop over all the node conditions
		for _, condition := range node.Status.Conditions {

			// Check the NodeReady condition
			if condition.Type == api.NodeReady && condition.Status != api.ConditionTrue {
				msg := fmt.Sprintf("%s, %s, %s", node.Name, condition.Reason, condition.Message)
				statusLine = append(statusLine, msg)
				statusCode = nagiosStatusCritical
			}

		}

	}

	if statusCode != nagiosStatusOK {
		fmt.Println(strings.Join(statusLine, "\n"))
		os.Exit(statusCode)
	}

	fmt.Println("OK")
	os.Exit(nagiosStatusOK)
}

func checkKubePods(c *cli.Context) {
	var err error
	var statusCode = nagiosStatusOK
	var statusLine []string
	var notReadyCount int

	kubeClient, err := client.New(kubeConfig)
	if err != nil {
		fmt.Printf("CRITICAL: %s\n", err)
		os.Exit(nagiosStatusUnknown)
	}

	// TODO: Allow for selecting pods based on labels or fields
	pods, err := kubeClient.Pods("").List(labels.Everything(), fields.Everything())
	if err != nil {
		fmt.Printf("CRITICAL: %s\n", err)
		os.Exit(nagiosStatusUnknown)
	}

	// Loop over all the pods
	for _, pod := range pods.Items {
		for _, cond := range pod.Status.Conditions {
			if cond.Type == "Ready" && cond.Status != "True" {
				notReadyCount++
			}
		}
	}

	if notReadyCount != 0 {
		msg := fmt.Sprintf("%d pods not in READY status.", notReadyCount)
		statusLine = append(statusLine, msg)
		statusCode = nagiosStatusCritical
	}

	if statusCode != nagiosStatusOK {
		fmt.Println(strings.Join(statusLine, "\n"))
		os.Exit(statusCode)
	}

	fmt.Println("OK")
	os.Exit(nagiosStatusOK)
}

func main() {
	app := cli.NewApp()
	app.Name = "check_kube_nodes"
	app.HelpName = app.Name
	app.Usage = "Nagios check to verify Kubernetes resources status"
	app.Version = appVersion

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "api-endpoint",
			Value:       "",
			Usage:       "Kubernetes API Endpoint",
			Destination: &kubeConfig.Host,
		},
		cli.StringFlag{
			Name:        "username",
			Value:       "",
			Usage:       "Kubernetes API Username",
			Destination: &kubeConfig.Username,
		},
		cli.StringFlag{
			Name:        "password",
			Value:       "",
			Usage:       "Kubernetes API Password",
			Destination: &kubeConfig.Password,
		},
		cli.BoolFlag{
			Name:        "skip-tls-verify",
			Usage:       "Skip TLS certificate verification",
			Destination: &kubeConfig.Insecure,
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:    "node",
			Aliases: []string{"n"},
			Usage:   "check node status",
			Action: func(c *cli.Context) {
				checkKubeNodes(c)
			},
		},
		cli.Command{
			Name:    "pod",
			Aliases: []string{"p"},
			Usage:   "check pod status",
			Action: func(c *cli.Context) {
				checkKubePods(c)
			},
		},
	}

	app.Run(os.Args)
}
