go build -o terraform-provider-datadog
export OS_ARCH="$(go env GOHOSTOS)_$(go env GOHOSTARCH)"
export TF_LOG=TRACE
export TF_LOG_PATH=terraform.log
mkdir -p ~/.terraform.d/plugins/hashicorp.com/edu/datadog/0.2/$OS_ARCH
mv terraform-provider-datadog ~/.terraform.d/plugins/hashicorp.com/edu/datadog/0.2/$OS_ARCH
rm -r /terraform-provider-datadog-custom/examples/.terraform
rm -r /terraform-provider-datadog-custom/examples/.terraform.lock.hcl
rm -r /terraform-provider-datadog-custom/examples/terraform.tfstate
rm -r /terraform-provider-datadog-custom/examples/terraform.tfstate.backup