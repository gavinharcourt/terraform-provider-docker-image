package dockerImage

import (
	"github.com/zongoose/terraform-provider-docker-image/src"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = dockerImage.Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"dockerImage": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := dockerImage.Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
