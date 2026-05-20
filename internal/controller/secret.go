package controller

import (
	"context"
	"fmt"

	akv1alpha1 "github.com/nayyara-airlangga/authentik-config-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ResolveSecret(ctx context.Context, k8sClient client.Client, ref akv1alpha1.SecretKeyRef, defaultNamespace string) (string, error) {
	ns := ref.Namespace
	if ns == "" {
		ns = defaultNamespace
	}

	secret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: ref.Name, Namespace: ns}, secret); err != nil {
		return "", fmt.Errorf("failed to fetch secret %s/%s: %w", ns, ref.Name, err)
	}

	value, ok := secret.Data[ref.Key]
	if !ok {
		return "", fmt.Errorf("key %q not found in secret %s/%s", ref.Key, ns, ref.Name)
	}

	return string(value), nil
}
