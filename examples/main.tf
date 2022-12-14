terraform {
  required_providers {
    datadog = {
      version = "0.2"
      source  = "hashicorp.com/edu/datadog"
    }
  }
}

provider "datadog" {
  
}

# module "psl" {
#   source = "./rqs"

# }

# output "psl" {
#   value = module.psl.all_restrictions
# }

resource "datadog_restriction" "restriction_test" {
  query = "scope:azerty2"
  roles {
    role_id = "b94f0ad0-b74e-11ec-802c-da7ad0900005"
  }
}
