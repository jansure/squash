package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/squash/pkg/api/v1"
)

// Purpose of script:
// This is a helper for interacting with squash CRDs

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func run() error {
	var name string
	var namespace string
	flag.StringVar(&name, "name", "some-name", "debug attachment name")
	flag.StringVar(&namespace, "namespace", "default", "debug attachment namespace")
	flag.Parse()

	attClient, err := getDebugAttachmentClient()
	if err != nil {
		return err
	}
	initialAtt := &v1.DebugAttachment{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Debugger: "dlv",
	}

	wResponse, err := attClient.Write(initialAtt, clients.WriteOpts{})
	if err != nil {
		return err
	}
	fmt.Println(wResponse)
	return nil
}

func getDebugAttachmentClient() (v1.DebugAttachmentClient, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, err
	}

	cache := kube.NewKubeCache()
	kubeRC := &factory.KubeResourceClientFactory{
		Crd:         v1.DebugAttachmentCrd,
		Cfg:         cfg,
		SharedCache: cache,
	}

	attClient, err := v1.NewDebugAttachmentClient(kubeRC)
	if err != nil {
		return nil, err
	}

	if err := attClient.Register(); err != nil {
		return nil, err
	}
	return attClient, nil

}
