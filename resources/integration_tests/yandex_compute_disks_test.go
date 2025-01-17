package integration_tests

import (
	"fmt"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/yandex-cloud/cq-provider-yandex/resources"

	"github.com/cloudquery/cq-provider-sdk/provider/providertest"
)

func TestIntegrationComputeDisks(t *testing.T) {
	var tfTmpl = fmt.Sprintf(`
resource "yandex_compute_disk" "foo" {
  name = "cq-disk-test-disk-%[1]s"
}
`, suffix)
	testIntegrationHelper(t, resources.ComputeDisks(), func(res *providertest.ResourceIntegrationTestData) providertest.ResourceIntegrationVerification {
		return providertest.ResourceIntegrationVerification{
			Name: "yandex_compute_disks",
			Filter: func(sq squirrel.SelectBuilder, _ *providertest.ResourceIntegrationTestData) squirrel.SelectBuilder {
				return sq
			},
			ExpectedValues: []providertest.ExpectedValue{{
				Count: 1,
				Data: map[string]interface{}{
					"name": fmt.Sprintf("cq-disk-test-disk-%s", suffix),
				},
			}},
		}
	}, tfTmpl)
}
