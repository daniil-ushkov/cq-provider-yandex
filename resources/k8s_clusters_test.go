package resources_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	k8s "github.com/yandex-cloud/go-sdk/gen/kubernetes"
	"google.golang.org/grpc"

	"github.com/cloudquery/cq-provider-sdk/logging"
	"github.com/cloudquery/cq-provider-sdk/provider/providertest"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/cloudquery/faker/v3"
	"github.com/hashicorp/go-hclog"
	"github.com/yandex-cloud/cq-provider-yandex/client"
	"github.com/yandex-cloud/cq-provider-yandex/resources"
	k8s1 "github.com/yandex-cloud/go-genproto/yandex/cloud/k8s/v1"
)

func TestK8SClusters(t *testing.T) {
	k8sSvc, serv, err := createClusterServer()
	if err != nil {
		t.Fatal(err)
	}
	resource := providertest.ResourceTestData{
		Table: resources.K8SClusters(),
		Config: client.Config{
			FolderIDs: []string{"test"},
		},
		Configure: func(logger hclog.Logger, _ interface{}) (schema.ClientMeta, error) {
			c := client.NewYandexClient(logging.New(&hclog.LoggerOptions{
				Level: hclog.Warn,
			}), []string{"test"}, nil, nil, &client.Services{
				K8S: k8sSvc,
			}, nil)
			return c, nil
		},
		Verifiers: []providertest.Verifier{
			providertest.VerifyAtLeastOneRow("yandex_k8s_clusters"),
		},
	}
	providertest.TestResource(t, resources.Provider, resource)
	serv.Stop()
}

type FakeClusterServiceServer struct {
	k8s1.UnimplementedClusterServiceServer
	Cluster *k8s1.Cluster
}

func NewFakeClusterServiceServer() (*FakeClusterServiceServer, error) {
	var cluster k8s1.Cluster
	faker.SetIgnoreInterface(true)
	err := faker.FakeData(&cluster)
	if err != nil {
		return nil, err
	}
	return &FakeClusterServiceServer{Cluster: &cluster}, nil
}

func (s *FakeClusterServiceServer) List(context.Context, *k8s1.ListClustersRequest) (*k8s1.ListClustersResponse, error) {
	return &k8s1.ListClustersResponse{Clusters: []*k8s1.Cluster{s.Cluster}}, nil
}

func createClusterServer() (*k8s.Kubernetes, *grpc.Server, error) {
	lis, err := net.Listen("tcp", ":0")

	if err != nil {
		return nil, nil, err
	}

	serv := grpc.NewServer()
	fakeClusterServiceServer, err := NewFakeClusterServiceServer()

	if err != nil {
		return nil, nil, err
	}

	k8s1.RegisterClusterServiceServer(serv, fakeClusterServiceServer)

	go func() {
		err := serv.Serve(lis)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())

	if err != nil {
		return nil, nil, err
	}

	return k8s.NewKubernetes(
		func(ctx context.Context) (*grpc.ClientConn, error) {
			return conn, nil
		},
	), serv, nil
}
