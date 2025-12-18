terraform {
  required_version = ">= 1.5"

  required_providers {
    openshift = {
      source  = "registry.terraform.io/rh-mobb/openshift"
      version = "~> 0.1"
    }
  }
}

provider "openshift" {
  kubeconfig = "~/.kube/config"
}

# Read existing operator status
data "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
}

# Only create resource if operator is ready
resource "null_resource" "gitops_ready" {
  count = data.openshift_operator.gitops.csv_phase == "Succeeded" ? 1 : 0

  triggers = {
    csv_phase = data.openshift_operator.gitops.csv_phase
    csv_version = data.openshift_operator.gitops.csv_version
  }

  provisioner "local-exec" {
    command = "echo 'GitOps operator is ready: ${data.openshift_operator.gitops.installed_csv}'"
  }
}

output "operator_status" {
  value = {
    installed_csv      = data.openshift_operator.gitops.installed_csv
    csv_phase          = data.openshift_operator.gitops.csv_phase
    csv_version        = data.openshift_operator.gitops.csv_version
    subscription_state = data.openshift_operator.gitops.subscription_state
    channel            = data.openshift_operator.gitops.channel
    source             = data.openshift_operator.gitops.source
  }
}

output "is_ready" {
  value       = data.openshift_operator.gitops.csv_phase == "Succeeded"
  description = "Whether the operator is in Succeeded phase"
}
