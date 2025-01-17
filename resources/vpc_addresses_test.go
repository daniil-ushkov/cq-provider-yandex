package resources_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"google.golang.org/grpc"

	"github.com/cloudquery/cq-provider-sdk/logging"
	"github.com/cloudquery/cq-provider-sdk/provider/providertest"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/cloudquery/faker/v3"
	"github.com/hashicorp/go-hclog"
	"github.com/yandex-cloud/cq-provider-yandex/client"
	"github.com/yandex-cloud/cq-provider-yandex/resources"
	vpc1 "github.com/yandex-cloud/go-genproto/yandex/cloud/vpc/v1"
	"github.com/yandex-cloud/go-sdk/gen/vpc"
)

func TestVPCAddresses(t *testing.T) {
	vpcSvc, serv, err := createAddressServer()
	if err != nil {
		t.Fatal(err)
	}
	resource := providertest.ResourceTestData{
		Table: resources.VPCAddresses(),
		Config: client.Config{
			FolderIDs: []string{"test"},
		},
		Configure: func(logger hclog.Logger, _ interface{}) (schema.ClientMeta, error) {
			c := client.NewYandexClient(logging.New(&hclog.LoggerOptions{
				Level: hclog.Warn,
			}), []string{"test"}, nil, nil, &client.Services{
				VPC: vpcSvc,
			}, nil)
			return c, nil
		},
		Verifiers: []providertest.Verifier{
			providertest.VerifyAtLeastOneRow("yandex_vpc_addresses"),
		},
	}
	providertest.TestResource(t, resources.Provider, resource)
	serv.Stop()
}

type FakeAddressServiceServer struct {
	vpc1.UnimplementedAddressServiceServer
	Address *vpc1.Address
}

func NewFakeAddressServiceServer() (*FakeAddressServiceServer, error) {
	var address vpc1.Address
	faker.SetIgnoreInterface(true)
	err := faker.FakeData(&address)
	if err != nil {
		return nil, err
	}
	return &FakeAddressServiceServer{Address: &address}, nil
}

func (s *FakeAddressServiceServer) List(context.Context, *vpc1.ListAddressesRequest) (*vpc1.ListAddressesResponse, error) {
	return &vpc1.ListAddressesResponse{Addresses: []*vpc1.Address{s.Address}}, nil
}

func createAddressServer() (*vpc.VPC, *grpc.Server, error) {
	lis, err := net.Listen("tcp", ":0")

	if err != nil {
		return nil, nil, err
	}

	serv := grpc.NewServer()
	fakeAddressServiceServer, err := NewFakeAddressServiceServer()

	if err != nil {
		return nil, nil, err
	}

	vpc1.RegisterAddressServiceServer(serv, fakeAddressServiceServer)

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

	return vpc.NewVPC(
		func(ctx context.Context) (*grpc.ClientConn, error) {
			return conn, nil
		},
	), serv, nil
}
