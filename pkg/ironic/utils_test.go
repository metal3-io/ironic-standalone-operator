package ironic

import (
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

func TestApplyContainerOverrides(t *testing.T) {
	testCases := []struct {
		Scenario string

		Existing  []corev1.Container
		Overrides []corev1.Container

		ExpectedNames []string
		Expected      []corev1.Container
	}{
		{
			Scenario: "No overrides",
			Existing: []corev1.Container{
				{Name: "c1", Image: "image1"},
				{Name: "c2", Image: "image2"},
			},
			Overrides: nil,

			ExpectedNames: []string{"c1", "c2"},
			Expected: []corev1.Container{
				{Name: "c1", Image: "image1"},
				{Name: "c2", Image: "image2"},
			},
		},
		{
			Scenario: "Empty overrides",
			Existing: []corev1.Container{
				{Name: "c1", Image: "image1"},
			},
			Overrides: []corev1.Container{},

			ExpectedNames: []string{"c1"},
			Expected: []corev1.Container{
				{Name: "c1", Image: "image1"},
			},
		},
		{
			Scenario: "Replace existing container",
			Existing: []corev1.Container{
				{Name: "c1", Image: "image1"},
				{Name: "c2", Image: "image2"},
			},
			Overrides: []corev1.Container{
				{Name: "c2", Image: "new-image2"},
			},

			ExpectedNames: []string{"c1", "c2"},
			Expected: []corev1.Container{
				{Name: "c1", Image: "image1"},
				{Name: "c2", Image: "new-image2"},
			},
		},
		{
			Scenario: "Append new container",
			Existing: []corev1.Container{
				{Name: "c1", Image: "image1"},
			},
			Overrides: []corev1.Container{
				{Name: "c3", Image: "image3"},
			},

			ExpectedNames: []string{"c1", "c3"},
			Expected: []corev1.Container{
				{Name: "c1", Image: "image1"},
				{Name: "c3", Image: "image3"},
			},
		},
		{
			Scenario: "Replace and append",
			Existing: []corev1.Container{
				{Name: "c1", Image: "image1"},
				{Name: "c2", Image: "image2"},
			},
			Overrides: []corev1.Container{
				{Name: "c1", Image: "new-image1"},
				{Name: "c3", Image: "image3"},
			},

			ExpectedNames: []string{"c1", "c2", "c3"},
			Expected: []corev1.Container{
				{Name: "c1", Image: "new-image1"},
				{Name: "c2", Image: "image2"},
				{Name: "c3", Image: "image3"},
			},
		},
		{
			Scenario: "Replace all containers",
			Existing: []corev1.Container{
				{Name: "c1", Image: "image1"},
				{Name: "c2", Image: "image2"},
			},
			Overrides: []corev1.Container{
				{Name: "c1", Image: "new-image1"},
				{Name: "c2", Image: "new-image2"},
			},

			ExpectedNames: []string{"c1", "c2"},
			Expected: []corev1.Container{
				{Name: "c1", Image: "new-image1"},
				{Name: "c2", Image: "new-image2"},
			},
		},
		{
			Scenario: "Multiple new containers",
			Existing: []corev1.Container{
				{Name: "c1", Image: "image1"},
			},
			Overrides: []corev1.Container{
				{Name: "c2", Image: "image2"},
				{Name: "c3", Image: "image3"},
			},

			ExpectedNames: []string{"c1", "c2", "c3"},
			Expected: []corev1.Container{
				{Name: "c1", Image: "image1"},
				{Name: "c2", Image: "image2"},
				{Name: "c3", Image: "image3"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			result := applyContainerOverrides(tc.Existing, tc.Overrides)

			var names []string
			for _, cont := range result {
				names = append(names, cont.Name)
			}

			assert.Equal(t, tc.ExpectedNames, names)
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
			InitContainers: []corev1.Container{
				{Name: "init1", Image: "init-image1"},
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

	testCases := []struct {
		Scenario string

		Overrides *metal3api.Overrides

		ExpectedAnnotations    map[string]string
		ExpectedLabels         map[string]string
		ExpectedContainers     []string
		ExpectedInitContainers []string
		ExpectedEnvs           map[string][]corev1.EnvVar
	}{
		{
			Scenario: "No overrides",

			ExpectedLabels:         map[string]string{"key1": "value1"},
			ExpectedContainers:     []string{"c1", "c2"},
			ExpectedInitContainers: []string{"init1"},
			ExpectedEnvs:           initialEnvs,
		},
		{
			Scenario:  "Empty overrides",
			Overrides: &metal3api.Overrides{},

			ExpectedLabels:         map[string]string{"key1": "value1"},
			ExpectedContainers:     []string{"c1", "c2"},
			ExpectedInitContainers: []string{"init1"},
			ExpectedEnvs:           initialEnvs,
		},
		{
			Scenario: "Keep builtin labels",
			Overrides: &metal3api.Overrides{
				Labels: map[string]string{"key1": "no value"},
			},

			ExpectedLabels:         map[string]string{"key1": "value1"},
			ExpectedContainers:     []string{"c1", "c2"},
			ExpectedInitContainers: []string{"init1"},
			ExpectedEnvs:           initialEnvs,
		},
		{
			Scenario: "New labels and annotations",
			Overrides: &metal3api.Overrides{
				Annotations: map[string]string{"key2": "value2"},
				Labels:      map[string]string{"key3": "value3"},
			},

			ExpectedAnnotations:    map[string]string{"key2": "value2"},
			ExpectedLabels:         map[string]string{"key1": "value1", "key3": "value3"},
			ExpectedContainers:     []string{"c1", "c2"},
			ExpectedInitContainers: []string{"init1"},
			ExpectedEnvs:           initialEnvs,
		},
		{
			Scenario: "Replace container",
			Overrides: &metal3api.Overrides{
				Containers: []corev1.Container{
					{
						Name: "c2",
						Env: []corev1.EnvVar{
							{Name: "new-env", Value: "new-value"},
						},
					},
				},
			},

			ExpectedLabels:         map[string]string{"key1": "value1"},
			ExpectedContainers:     []string{"c1", "c2"},
			ExpectedInitContainers: []string{"init1"},
			ExpectedEnvs: map[string][]corev1.EnvVar{
				"c1": {
					{
						Name:  "env1",
						Value: "value1",
					},
				},
				"c2": {
					{Name: "new-env", Value: "new-value"},
				},
			},
		},
		{
			Scenario: "Add new container",
			Overrides: &metal3api.Overrides{
				Containers: []corev1.Container{
					{Name: "c3", Image: "image3"},
				},
			},

			ExpectedLabels:         map[string]string{"key1": "value1"},
			ExpectedContainers:     []string{"c1", "c2", "c3"},
			ExpectedInitContainers: []string{"init1"},
			ExpectedEnvs: map[string][]corev1.EnvVar{
				"c1": {
					{
						Name:  "env1",
						Value: "value1",
					},
				},
				"c2": nil,
				"c3": nil,
			},
		},
		{
			Scenario: "Replace init container",
			Overrides: &metal3api.Overrides{
				InitContainers: []corev1.Container{
					{Name: "init1", Image: "new-init-image1"},
				},
			},

			ExpectedLabels:         map[string]string{"key1": "value1"},
			ExpectedContainers:     []string{"c1", "c2"},
			ExpectedInitContainers: []string{"init1"},
			ExpectedEnvs:           initialEnvs,
		},
		{
			Scenario: "Add new init container",
			Overrides: &metal3api.Overrides{
				InitContainers: []corev1.Container{
					{Name: "init2", Image: "init-image2"},
				},
			},

			ExpectedLabels:         map[string]string{"key1": "value1"},
			ExpectedContainers:     []string{"c1", "c2"},
			ExpectedInitContainers: []string{"init1", "init2"},
			ExpectedEnvs:           initialEnvs,
		},
		{
			Scenario: "Full override with containers and init containers",
			Overrides: &metal3api.Overrides{
				Annotations: map[string]string{"anno1": "val1"},
				Labels:      map[string]string{"label1": "val1"},
				Containers: []corev1.Container{
					{Name: "c1", Image: "new-c1"},
					{Name: "c3", Image: "c3"},
				},
				InitContainers: []corev1.Container{
					{Name: "init1", Image: "new-init1"},
					{Name: "init2", Image: "init2"},
				},
			},

			ExpectedAnnotations:    map[string]string{"anno1": "val1"},
			ExpectedLabels:         map[string]string{"key1": "value1", "label1": "val1"},
			ExpectedContainers:     []string{"c1", "c2", "c3"},
			ExpectedInitContainers: []string{"init1", "init2"},
			ExpectedEnvs: map[string][]corev1.EnvVar{
				"c1": nil,
				"c2": nil,
				"c3": nil,
			},
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

			var initContainerNames []string
			for _, cont := range result.Spec.InitContainers {
				initContainerNames = append(initContainerNames, cont.Name)
			}

			assert.Equal(t, tc.ExpectedContainers, containerNames)
			assert.Equal(t, tc.ExpectedInitContainers, initContainerNames)
			assert.Equal(t, tc.ExpectedAnnotations, result.Annotations)
			assert.Equal(t, tc.ExpectedLabels, result.Labels)
			assert.Equal(t, tc.ExpectedEnvs, envs)
		})
	}
}
