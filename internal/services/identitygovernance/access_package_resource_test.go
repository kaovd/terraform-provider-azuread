package identitygovernance_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/manicminer/hamilton/odata"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
)

type AccessPackageResource struct{}

func TestAccAccessPackage_simple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_access_package", "test")
	r := AccessPackageResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.AP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAccessPackage_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_access_package", "test")
	r := AccessPackageResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.AP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.APUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.AP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AccessPackageResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	_, status, err := clients.IdentityGovernance.AccessPackageClient.Get(ctx, state.ID, odata.Query{})
	if err != nil {
		if status == http.StatusNotFound {
			return nil, fmt.Errorf("Access package with object ID %q does not exist", state.ID)
		}
		return nil, fmt.Errorf("failed to retrieve access package with object ID %q: %+v", state.ID, err)
	}

	return utils.Bool(true), nil
}

func (r AccessPackageResource) AP(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_access_package" "test" {
  catalog_id = azuread_access_package_catalog.test.id
  description = <<DESC
My new
Access Package
DESC
  display_name = "acctestAP-%[2]d"
  is_hidden = true
}
`, r.catalogTemplate(data), data.RandomInteger)
}

func (r AccessPackageResource) APUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_access_package" "test" {
	catalog_id = azuread_access_package_catalog.test.id
	description = <<DESC
  My new
  and Improved!
  Access Package
DESC
	display_name = "acctestAPUpdate-%[2]d"
	is_hidden = false
}
`, r.catalogTemplate(data), data.RandomInteger)
}

func (AccessPackageResource) catalogTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_access_package_catalog" "test" {
	display_name = "acctestAPC-%[1]d"
	catalog_status = "Published"
	description = "My test Catalog"
	is_externally_visible = false
}
`, data.RandomInteger)
}