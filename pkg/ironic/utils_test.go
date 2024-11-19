package ironic

import (
	"maps"
	"testing"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBuildEndpoints(t *testing.T) {
	testCases := []struct {
		Scenario string

		IPs          []string
		Port         int
		IncludeProto string

		Expected []string
	}{
		{
			Scenario: "non-standard-port-no-protocol",

			IPs:          []string{"2001:db8::42", "192.0.2.42"},
			Port:         6385,
			IncludeProto: "",

			Expected: []string{"192.0.2.42:6385", "[2001:db8::42]:6385"},
		},
		{
			Scenario: "non-standard-port-with-protocol",

			IPs:          []string{"2001:db8::42", "192.0.2.42"},
			Port:         6385,
			IncludeProto: "http",

			Expected: []string{"http://192.0.2.42:6385", "http://[2001:db8::42]:6385"},
		},
		{
			Scenario: "http-with-protocol",

			IPs:          []string{"2001:db8::42", "192.0.2.42"},
			Port:         80,
			IncludeProto: "http",

			Expected: []string{"http://192.0.2.42", "http://[2001:db8::42]"},
		},
		{
			Scenario: "https-with-protocol",

			IPs:          []string{"2001:db8::42", "192.0.2.42"},
			Port:         443,
			IncludeProto: "https",

			Expected: []string{"https://192.0.2.42", "https://[2001:db8::42]"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			result := buildEndpoints(tc.IPs, tc.Port, tc.IncludeProto)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func TestApplyOverridesToPod(t *testing.T) {
	initial := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{"key1": "value1"},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "c1",
					Env: []corev1.EnvVar{
						{
							Name:  "env1",
							Value: "value1",
						},
					},
				},
				{Name: "c2"},
			},
		},
	}

	initialEnvs := map[string][]corev1.EnvVar{
		"c1": {
			{
				Name:  "env1",
				Value: "value1",
			},
		},
		"c2": nil,
	}
	addedC3Envs := maps.Clone(initialEnvs)
	addedC3Envs["c3"] = nil
	updatedEnvs := map[string][]corev1.EnvVar{
		"c1": {
			{
				Name:  "env1",
				Value: "value1",
			},
			{
				Name:  "env2",
				Value: "value2",
			},
		},
		"c2": {
			{
				Name:  "env2",
				Value: "value2",
			},
		},
		"c3": {
			{
				Name:  "env2",
				Value: "value2",
			},
		},
	}

	testCases := []struct {
		Scenario string

		Overrides *metal3api.Overrides

		ExpectedContainerNames []string
		ExpectedAnnotations    map[string]string
		ExpectedLabels         map[string]string
		ExpectedEnvs           map[string][]corev1.EnvVar
	}{
		{
			Scenario: "No overrides",

			ExpectedContainerNames: []string{"c1", "c2"},
			ExpectedLabels:         map[string]string{"key1": "value1"},
			ExpectedEnvs:           initialEnvs,
		},
		{
			Scenario:  "Empty overrides",
			Overrides: &metal3api.Overrides{},

			ExpectedContainerNames: []string{"c1", "c2"},
			ExpectedLabels:         map[string]string{"key1": "value1"},
			ExpectedEnvs:           initialEnvs,
		},
		{
			Scenario: "Keep builtin labels",
			Overrides: &metal3api.Overrides{
				Labels: map[string]string{"key1": "no value"},
			},

			ExpectedContainerNames: []string{"c1", "c2"},
			ExpectedLabels:         map[string]string{"key1": "value1"},
			ExpectedEnvs:           initialEnvs,
		},
		{
			Scenario: "New containers and meta",
			Overrides: &metal3api.Overrides{
				Annotations: map[string]string{"key2": "value2"},
				Containers: []corev1.Container{
					{Name: "c3"},
				},
				Labels: map[string]string{"key3": "value3"},
			},

			ExpectedContainerNames: []string{"c1", "c2", "c3"},
			ExpectedAnnotations:    map[string]string{"key2": "value2"},
			ExpectedLabels:         map[string]string{"key1": "value1", "key3": "value3"},
			ExpectedEnvs:           addedC3Envs,
		},
		{
			Scenario: "New containers and env",
			Overrides: &metal3api.Overrides{
				Containers: []corev1.Container{
					{Name: "c3"},
				},
				Env: []corev1.EnvVar{
					{
						Name:  "env2",
						Value: "value2",
					},
				},
			},

			ExpectedContainerNames: []string{"c1", "c2", "c3"},
			ExpectedLabels:         map[string]string{"key1": "value1"},
			ExpectedEnvs:           updatedEnvs,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			result := applyOverridesToPod(tc.Overrides, *initial.DeepCopy())

			var containerNames []string
			envs := make(map[string][]corev1.EnvVar)

			for _, cont := range result.Spec.Containers {
				containerNames = append(containerNames, cont.Name)
				envs[cont.Name] = cont.Env
			}

			assert.Equal(t, tc.ExpectedContainerNames, containerNames)
			assert.Equal(t, tc.ExpectedAnnotations, result.Annotations)
			assert.Equal(t, tc.ExpectedLabels, result.Labels)
			assert.Equal(t, tc.ExpectedEnvs, envs)
		})
	}
}
