package controller

import (
	"github.com/kibiluzbad/activemq-artemis-operator/pkg/controller/activemqartemis"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, activemqartemis.Add)
}
