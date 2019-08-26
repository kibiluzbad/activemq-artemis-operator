package activemqartemisdivert

import (
	"context"
	"fmt"
	"strconv"

	mgmt "github.com/kibiluzbad/activemq-artemis-management"
	brokerv2alpha1 "github.com/kibiluzbad/activemq-artemis-operator/pkg/apis/broker/v2alpha1"
	ss "github.com/kibiluzbad/activemq-artemis-operator/pkg/resources/statefulsets"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_activemqartemisdivert")
var namespacedNameToDivertName = make(map[types.NamespacedName]brokerv2alpha1.ActiveMQArtemisDivert)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ActiveMQArtemisDivert Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileActiveMQArtemisDivert{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("activemqartemisdivert-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ActiveMQArtemisDivert
	err = c.Watch(&source.Kind{Type: &brokerv2alpha1.ActiveMQArtemisDivert{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner ActiveMQArtemisDivert
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &brokerv2alpha1.ActiveMQArtemisDivert{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileActiveMQArtemisDivert{}

// ReconcileActiveMQArtemisDivert reconciles a ActiveMQArtemisDivert object
type ReconcileActiveMQArtemisDivert struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ActiveMQArtemisDivert object and makes changes based on the state read
// and what is in the ActiveMQArtemisDivert.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileActiveMQArtemisDivert) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ActiveMQArtemisDivert")

	// Fetch the ActiveMQArtemisDivert instance
	instance := &brokerv2alpha1.ActiveMQArtemisDivert{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		// Delete action
		divertInstance, lookupSucceeded := namespacedNameToDivertName[request.NamespacedName]
		if lookupSucceeded {

			err = deleteDivert(&divertInstance, request, r.client)
			delete(namespacedNameToDivertName, request.NamespacedName)
		}
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	err = createDivert(instance, request, r.client)
	if nil == err {
		namespacedNameToDivertName[request.NamespacedName] = *instance //.Spec.QueueName
	}

	return reconcile.Result{}, nil
}

func createDivert(instance *brokerv2alpha1.ActiveMQArtemisDivert, request reconcile.Request, client client.Client) error {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Creating ActiveMQArtemisDivert")

	var err error
	artemisArray := getPodBrokers(instance, request, client)
	if nil != artemisArray {
		for _, a := range artemisArray {
			if nil == a {
				reqLogger.Info("Creating ActiveMQArtemisDivert artemisArray had a nil!")
				continue
			}
			result, err := a.CreateDivert(instance.Spec.Name,
				instance.Spec.RoutingName,
				instance.Spec.Address,
				instance.Spec.ForwardingAddress,
				instance.Spec.Exclusive,
				instance.Spec.Filter,
				instance.Spec.Transformer)
			if nil != err || result.Status != 200 {
				reqLogger.Info("Error: " + result.Error + ", Status: " + strconv.Itoa(result.Status))
				reqLogger.Info("Creating ActiveMQArtemisDivert error for " + instance.Spec.Name)
				break
			} else {

				reqLogger.Info("Created ActiveMQArtemisDivert for " + instance.Spec.Name)
			}
		}
	}

	return err
}

func deleteDivert(instance *brokerv2alpha1.ActiveMQArtemisDivert, request reconcile.Request, client client.Client) error {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Deleting ActiveMQArtemisDivert")

	var err error
	artemisArray := getPodBrokers(instance, request, client)
	if nil != artemisArray {
		for _, a := range artemisArray {
			_, err := a.DeleteDivert(instance.Spec.Name)
			if nil != err {
				reqLogger.Info("Deleting ActiveMQArtemisDivert error for " + instance.Spec.Name)
				break
			} else {
				reqLogger.Info("Deleted ActiveMQArtemisDivert for " + instance.Spec.Name)
			}
		}
	}

	return err
}

func getPodBrokers(instance *brokerv2alpha1.ActiveMQArtemisDivert, request reconcile.Request, client client.Client) []*mgmt.Artemis {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Getting Pod Brokers")

	var artemisArray []*mgmt.Artemis
	var err error

	ss.NameBuilder.Name()
	if err != nil {
		reqLogger.Error(err, "Failed to ge the statefulset name")
	}

	// Check to see if the statefulset already exists
	ssNamespacedName := types.NamespacedName{
		Name:      ss.NameBuilder.Name(),
		Namespace: request.Namespace,
	}

	statefulset, err := ss.RetrieveStatefulSet(ss.NameBuilder.Name(), ssNamespacedName, client)
	if nil != err {
		reqLogger.Info("Statefulset: " + ssNamespacedName.Name + " not found")
	} else {
		reqLogger.Info("Statefulset: " + ssNamespacedName.Name + " found")
		fmt.Printf("%v+", statefulset)

		pod := &corev1.Pod{}
		podNamespacedName := types.NamespacedName{
			Name:      statefulset.Name + "-0",
			Namespace: request.Namespace,
		}

		// For each of the replicas
		i := 0
		replicas := int(*statefulset.Spec.Replicas)
		artemisArray = make([]*mgmt.Artemis, 0, replicas)
		for i = 0; i < replicas; i++ {
			s := statefulset.Name + "-" + strconv.Itoa(i)
			podNamespacedName.Name = s
			if err = client.Get(context.TODO(), podNamespacedName, pod); err != nil {
				if errors.IsNotFound(err) {
					reqLogger.Error(err, "Pod IsNotFound", "Namespace", request.Namespace, "Name", request.Name)
				} else {
					reqLogger.Error(err, "Pod lookup error", "Namespace", request.Namespace, "Name", request.Name)
				}
			} else {
				reqLogger.Info("Pod found", "Namespace", request.Namespace, "Name", request.Name)
				artemis := mgmt.NewArtemis(pod.Status.PodIP, "8161", "amq-broker")
				artemisArray = append(artemisArray, artemis)
			}
		}
	}

	return artemisArray
}
