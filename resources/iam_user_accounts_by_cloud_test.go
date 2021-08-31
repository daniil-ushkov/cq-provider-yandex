package resources_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"testing"

	iam1 "github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
	"github.com/yandex-cloud/go-sdk/gen/iam"
	"google.golang.org/grpc"

	"github.com/cloudquery/cq-provider-sdk/logging"
	"github.com/cloudquery/cq-provider-sdk/provider/providertest"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/hashicorp/go-hclog"
	"github.com/yandex-cloud/cq-provider-yandex/client"
	"github.com/yandex-cloud/cq-provider-yandex/resources"
)

func TestUserAccountsByCloud(t *testing.T) {
	accessBindingsByCloudServer, serv1, err := createAccessBindingsByCloudServer()
	defer serv1.Stop()
	if err != nil {
		t.Fatal(err)
	}

	userAccountsByCloudServer, serv2, err := createUserAccountsByCloudServer()
	defer serv2.Stop()
	if err != nil {
		t.Fatal(err)
	}

	resource := providertest.ResourceTestData{
		Table:  resources.IAMUserAccountsByCloud(),
		Config: client.Config{CloudIDs: []string{"test"}},
		Configure: func(logger hclog.Logger, _ interface{}) (schema.ClientMeta, error) {
			c := client.NewYandexClient(
				logging.New(&hclog.LoggerOptions{
					Level: hclog.Warn,
				}),
				nil,
				[]string{"test"},
				nil,
				&client.Services{
					IAM:             userAccountsByCloudServer,
					ResourceManager: accessBindingsByCloudServer,
				}, nil)
			return c, nil
		},
		Verifiers: []providertest.Verifier{
			providertest.VerifyAtLeastOneRow("yandex_iam_user_accounts_by_cloud"),
			providertest.VerifyNoEmptyColumnsExcept(
				"yandex_iam_user_accounts_by_cloud",
				"user_account_yandex_passport_user_account_login",
				"user_account_yandex_passport_user_account_default_email",
				"user_account_saml_user_account_federation_id",
				"user_account_saml_user_account_name_id",
				"user_account_saml_user_account_attributes",
			),
			providertest.VerifyOneOf(
				"yandex_iam_user_accounts_by_cloud",
				"user_account_yandex_passport_user_account_login",
				"user_account_saml_user_account_name_id",
			),
		},
	}
	providertest.TestResource(t, resources.Provider, resource)
}

type FakeUserAccountByCloudServer struct {
	iam1.UnimplementedUserAccountServiceServer
}

func (f *FakeUserAccountByCloudServer) Get(_ context.Context, req *iam1.GetUserAccountRequest) (*iam1.UserAccount, error) {
	switch req.UserAccountId {
	case "1":
		return &iam1.UserAccount{Id: "1", UserAccount: &iam1.UserAccount_SamlUserAccount{
			SamlUserAccount: &iam1.SamlUserAccount{
				FederationId: "1",
				NameId:       "1",
			},
		}}, nil
	case "2":
		return &iam1.UserAccount{Id: "2", UserAccount: &iam1.UserAccount_YandexPassportUserAccount{
			YandexPassportUserAccount: &iam1.YandexPassportUserAccount{
				Login:        "qwerty",
				DefaultEmail: "qwerty@qwerty.com",
			},
		}}, nil
	default:
		return nil, errors.New("no such user account")
	}
}

func createUserAccountsByCloudServer() (*iam.IAM, *grpc.Server, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, nil, err
	}

	serv := grpc.NewServer()
	fakeCloudServiceServer := &FakeUserAccountByCloudServer{}

	iam1.RegisterUserAccountServiceServer(serv, fakeCloudServiceServer)

	go func() {
		err := serv.Serve(lis)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	return iam.NewIAM(
		func(ctx context.Context) (*grpc.ClientConn, error) {
			return conn, nil
		},
	), serv, nil
}