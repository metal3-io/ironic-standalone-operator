package v1alpha1

import (
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

const (
	FeatureHighAvailability featuregate.Feature = "HighAvailability"
)

var (
	availableFeatures = map[featuregate.Feature]featuregate.FeatureSpec{
		FeatureHighAvailability: {Default: false, PreRelease: featuregate.Alpha},
	}

	CurrentFeatureGate = featuregate.NewFeatureGate()
)

func init() {
	utilruntime.Must(CurrentFeatureGate.Add(availableFeatures))
}
