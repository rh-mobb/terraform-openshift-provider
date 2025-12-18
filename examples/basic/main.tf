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
  # Use default kubeconfig or set explicitly
  # kubeconfig = "~/.kube/config"
}

# Install GitOps Operator
resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
  channel   = "latest"
  source    = "redhat-operators"
}

output "gitops_csv" {
  value       = openshift_operator.gitops.installed_csv
  description = "Installed GitOps CSV name"
}

output "gitops_phase" {
  value       = openshift_operator.gitops.csv_phase
  description = "GitOps CSV phase"
}
