terraform {
  required_providers {
    datadog = {
      version = "0.2"
      source  = "hashicorp.com/edu/datadog"
    }
  }
}

data "datadog_restrictions" "all" {}

# Returns all coffees
output "all_restrictions" {
  value = data.datadog_restrictions.all.restrictions
}
