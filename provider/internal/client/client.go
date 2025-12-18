package client

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Kubernetes kubernetes.Interface
	Dynamic    dynamic.Interface
	Config     *rest.Config
}

func NewClient(kubeconfig, host, token string, insecure bool) (*Client, error) {
	var config *rest.Config
	var err error

	// Priority: kubeconfig file > host+token > in-cluster config
	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
		}
	} else if host != "" && token != "" {
		config = &rest.Config{
			Host:        host,
			BearerToken: token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: insecure,
			},
		}
	} else {
		// Try to use kubeconfig from default locations
		home, _ := os.UserHomeDir()
		kubeconfigPath := filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(kubeconfigPath); err == nil {
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load default kubeconfig: %w", err)
			}
		} else {
			// Try in-cluster config
			config, err = rest.InClusterConfig()
			if err != nil {
				return nil, fmt.Errorf("no valid Kubernetes configuration found: %w", err)
			}
		}
	}

	// Override host and token if provided
	if host != "" {
		config.Host = host
	}
	if token != "" {
		config.BearerToken = token
	}
	if insecure {
		config.TLSClientConfig = rest.TLSClientConfig{
			Insecure: true,
		}
	} else if config.TLSClientConfig.CAFile == "" && config.TLSClientConfig.CAData == nil {
		// Set insecure if no CA is configured
		config.TLSClientConfig = rest.TLSClientConfig{
			Insecure: true,
		}
	}

	// Create clients
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &Client{
		Kubernetes: k8sClient,
		Dynamic:    dynamicClient,
		Config:     config,
	}, nil
}

// VerifyTLSConfig ensures TLS is properly configured
func (c *Client) VerifyTLSConfig() error {
	if c.Config.TLSClientConfig.Insecure {
		return nil
	}

	// Test connection
	transport, err := rest.TransportFor(c.Config)
	if err != nil {
		return fmt.Errorf("failed to create transport: %w", err)
	}

	// Try to get TLS config
	if tlsConfig, err := rest.TLSConfigFor(c.Config); err != nil {
		return fmt.Errorf("failed to get TLS config: %w", err)
	} else if tlsConfig == nil {
		return fmt.Errorf("TLS config is nil")
	}

	_ = transport // Use transport to avoid unused variable
	return nil
}
