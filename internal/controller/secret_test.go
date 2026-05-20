package controller

import (
	"fmt"
	"os"
	"testing"

	akv1alpha1 "github.com/nayyara-airlangga/authentik-config-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func TestMain(m *testing.M) {
	testEnv := &envtest.Environment{}

	if getFirstFoundEnvTestBinaryDir() != "" {
		testEnv.BinaryAssetsDirectory = getFirstFoundEnvTestBinaryDir()
	}

	cfg, err := testEnv.Start()
	if err != nil {
		panic(fmt.Errorf("failed to start test environment: %w", err))
	}

	k8sClient, err = client.New(cfg, client.Options{})
	if err != nil {
		panic(fmt.Errorf("failed to initialize k8s test client: %w", err))
	}

	code := m.Run()

	_ = testEnv.Stop()

	os.Exit(code)
}

func TestResolveSecret(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "authentik-token",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"token": []byte("my-api-token"),
		},
	}

	require.NoError(t, k8sClient.Create(t.Context(), secret))

	t.Cleanup(func() {
		_ = k8sClient.Delete(t.Context(), secret)
	})

	t.Run("resolves secret successfully", func(t *testing.T) {
		val, err := ResolveSecret(t.Context(), k8sClient, akv1alpha1.SecretKeyRef{
			Name:      "authentik-token",
			Namespace: "default",
			Key:       "token",
		}, "default")

		require.NoError(t, err)
		assert.Equal(t, "my-api-token", val)
	})
}
