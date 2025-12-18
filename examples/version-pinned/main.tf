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
  kubeconfig = var.kubeconfig_path
}

variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

# Install GitOps Operator with version pinning
resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
  channel   = "latest"
  source    = "redhat-operators"
  version   = "1.18.2"  # Automatically sets install_plan_approval to "Manual"

  labels = {
    "app.kubernetes.io/managed-by" = "Terraform"
    "environment"                  = "production"
  }

  namespace_labels = {
    "openshift.io/cluster-monitoring" = "true"
  }
}

output "gitops_version" {
  value       = openshift_operator.gitops.installed_csv
  description = "Installed GitOps operator version"
}
