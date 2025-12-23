package secretutils

import (
	"maps"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	LabelEnvironmentName  = "environment.metal3.io/ironic-standalone-operator"
	LabelEnvironmentValue = "true"
)

// AddSecretSelector adds a selector to a cache.ByObject map that filters
// Secrets and ConfigMaps so that only those labelled as part of the ironic environment
// get cached. The input may be nil.
func AddSecretSelector(selectors map[client.Object]cache.ByObject) map[client.Object]cache.ByObject {
	secret := &corev1.Secret{}
	configMap := &corev1.ConfigMap{}
	newSelectors := map[client.Object]cache.ByObject{
		configMap: {
			Label: labels.SelectorFromSet(
				labels.Set{
					LabelEnvironmentName: LabelEnvironmentValue,
				}),
		},
		secret: {
			Label: labels.SelectorFromSet(
				labels.Set{
					LabelEnvironmentName: LabelEnvironmentValue,
				}),
		},
	}

	if selectors == nil {
		return newSelectors
	}

	maps.Insert(selectors, maps.All(newSelectors))
	return selectors
}
