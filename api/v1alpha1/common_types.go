/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

// AuthentikEndpoint defines the configuration to connect to an Authentik endpoint.
type AuthentikEndpoint struct {
	// The base URL of the Authentik instance (e.g. https://authentik.example.com).
	// +required
	// +kubebuilder:validation:Pattern=`^https?://.+`
	Host string `json:"host"`

	// TokenSecret is a reference to the Kubernetes secret that contains the
	// Authentik API token for authentication. The key used in the secret must
	// contain the API token value (typically "apiToken").
	// +required
	TokenSecret SecretKeyRef `json:"tokenSecret"`

	// TLSConfig defines custom TLS settings for connecting to the Authentik
	// instance. If omitted, the system's default CA bundle is used.
	// +optional
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty"`
}

// SecretKeyRef is a reference to a key in a Kubernetes secret.
type SecretKeyRef struct {
	// Name of the secret.
	// +required
	Name string `json:"name"`

	// Namespace of the secret. Defaults to the namespace of the referencing resource.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Key inside the secret that will be refered.
	// +required
	Key string `json:"key"`
}

// TLSConfig defines the TLS configuration that will be used when connecting
// to an endpoint.
type TLSConfig struct {
	// CASecret is a reference to the Kubernetes secret that contains a
	// CA certificate bundle. The key used in the secret must be the one
	// containing the PEM-encoded CA certificate (typically "ca.crt").
	// +optional
	CASecret *SecretKeyRef `json:"caSecret,omitempty"`

	// InsecureSkipVerify disables TLS certificate verification.
	// +optional
	// +kubebuilder:default=false
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
}
