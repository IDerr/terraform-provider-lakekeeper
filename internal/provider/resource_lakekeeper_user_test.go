//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperUser_basic(t *testing.T) {

	rID := fmt.Sprintf("oidc~%s", acctest.RandString(8))
	rName := acctest.RandString(8)
	rUpdatedName := acctest.RandString(12)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`				
				resource "lakekeeper_user" "foo" {
				  id = "%s"
				  name = "%s"
				  email = "%s@local.local"
				  user_type = "human"
				}
				`, rID, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "id", rID),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "email", rName+"@local.local"),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "user_type", "human"),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "last_updated_with", "create-endpoint"),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "updated_at", ""),
					resource.TestCheckResourceAttrSet("lakekeeper_user.foo", "created_at"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_user.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update User
			{
				Config: fmt.Sprintf(`				
				resource "lakekeeper_user" "foo" {
				  id = "%s"
				  name = "%s"
				  email = "%s@local.local"
				  user_type = "human"
				}
				`, rID, rUpdatedName, rUpdatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "id", rID),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "name", rUpdatedName),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "email", rUpdatedName+"@local.local"),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "user_type", "human"),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "last_updated_with", "create-endpoint"),
					resource.TestCheckResourceAttrSet("lakekeeper_user.foo", "updated_at"),
					resource.TestCheckResourceAttrSet("lakekeeper_user.foo", "created_at"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_user.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLakekeeperUserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_user" {
			continue
		}

		_, err := testutil.TestLakekeeperClient.GetUserByID(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("User with id %s still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}
