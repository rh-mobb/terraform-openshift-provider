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
  # Configure provider authentication
  # Option 1: Use kubeconfig
  # kubeconfig = "~/.kube/config"

  # Option 2: Use host and token
  # host     = "https://api.example.com:6443"
  # token    = var.k8s_token
  # insecure = false
}

# Install GitOps Operator
resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
  channel   = "latest"
  source    = "redhat-operators"
  version   = "1.18.2"
}

# Install Service Mesh Operator
resource "openshift_operator" "service_mesh" {
  name      = "servicemeshoperator"
  namespace = "openshift-operators"
  channel   = "stable"
  source    = "redhat-operators"

  wait_for_csv = true
  wait_timeout = "15m"
}

# Install Serverless Operator
resource "openshift_operator" "serverless" {
  name      = "serverless-operator"
  namespace = "openshift-serverless"
  channel   = "stable"
  source    = "redhat-operators"

  create_namespace = true
  namespace_labels = {
    "openshift.io/cluster-monitoring" = "true"
  }
}

output "operators" {
  value = {
    gitops = {
      csv   = openshift_operator.gitops.installed_csv
      phase = openshift_operator.gitops.csv_phase
    }
    service_mesh = {
      csv   = openshift_operator.service_mesh.installed_csv
      phase = openshift_operator.service_mesh.csv_phase
    }
    serverless = {
      csv   = openshift_operator.serverless.installed_csv
      phase = openshift_operator.serverless.csv_phase
    }
  }
}
