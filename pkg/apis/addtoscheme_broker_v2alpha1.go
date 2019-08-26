package apis

import (
	"github.com/kibiluzbad/activemq-artemis-operator/pkg/apis/broker/v2alpha1"
	routev1 "github.com/openshift/api/route/v1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v2alpha1.SchemeBuilder.AddToScheme, routev1.SchemeBuilder.AddToScheme)
}
