package activemqartemisuser

import (
	"context"
	"fmt"
	"strconv"

	mgmt "github.com/kibiluzbad/activemq-artemis-management"
	brokerv1alpha1 "github.com/kibiluzbad/activemq-artemis-operator/pkg/apis/broker/v1alpha1"
	aa "github.com/kibiluzbad/activemq-artemis-operator/pkg/controller/activemqartemis"
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

var log = logf.Log.WithName("controller_activemqartemisuser")
var namespacedNameToUserName = make(map[types.NamespacedName]brokerv1alpha1.ActiveMQArtemisUser)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ActiveMQArtemisUser Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileActiveMQArtemisUser{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("activemqartemisuser-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ActiveMQArtemisUser
	err = c.Watch(&source.Kind{Type: &brokerv1alpha1.ActiveMQArtemisUser{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner ActiveMQArtemisUser
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &brokerv1alpha1.ActiveMQArtemisUser{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileActiveMQArtemisUser{}

// ReconcileActiveMQArtemisUser reconciles a ActiveMQArtemisUser object
type ReconcileActiveMQArtemisUser struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ActiveMQArtemisUser object and makes changes based on the state read
// and what is in the ActiveMQArtemisUser.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileActiveMQArtemisUser) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ActiveMQArtemisUser")

	// Fetch the ActiveMQArtemisUser instance
	instance := &brokerv1alpha1.ActiveMQArtemisUser{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		// Delete action
		userInstance, lookupSucceeded := namespacedNameToUserName[request.NamespacedName]
		if lookupSucceeded {

			err = deleteUser(&userInstance, request, r.client)
			delete(namespacedNameToUserName, request.NamespacedName)
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

	err = createUser(instance, request, r.client)
	if nil == err {
		namespacedNameToUserName[request.NamespacedName] = *instance //.Spec.QueueName
	}

	return reconcile.Result{}, nil
}

func createUser(instance *brokerv1alpha1.ActiveMQArtemisUser, request reconcile.Request, client client.Client) error {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Creating ActiveMQArtemisUser")

	var err error
	artemisArray := getPodBrokers(instance, request, client)
	if nil != artemisArray {
		for _, a := range artemisArray {
			if nil == a {
				reqLogger.Info("Creating ActiveMQArtemisUser artemisArray had a nil!")
				continue
			}
			result, err := a.CreateUser(instance.Spec.UserName,
				instance.Spec.Password,
				instance.Spec.Roles)
			if nil != err || result.Status != 200 {
				reqLogger.Info("Error: " + result.Error + ", Status: " + strconv.Itoa(result.Status))
				reqLogger.Info("Creating ActiveMQArtemisUser error for " + instance.Spec.UserName)
				break
			} else {

				reqLogger.Info("Created ActiveMQArtemisUser for " + instance.Spec.UserName)
			}
		}
	}

	return err
}

func deleteUser(instance *brokerv1alpha1.ActiveMQArtemisUser, request reconcile.Request, client client.Client) error {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Deleting ActiveMQArtemisUser")

	var err error
	artemisArray := getPodBrokers(instance, request, client)
	if nil != artemisArray {
		for _, a := range artemisArray {
			_, err := a.RemoveUser(instance.Spec.UserName)
			if nil != err {
				reqLogger.Info("Deleting ActiveMQArtemisUser error for " + instance.Spec.UserName)
				break
			} else {
				reqLogger.Info("Deleted ActiveMQArtemisUser for " + instance.Spec.UserName)
			}
		}
	}

	return err
}

func getPodBrokers(instance *brokerv1alpha1.ActiveMQArtemisUser, request reconcile.Request, client client.Client) []*mgmt.Artemis {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Getting Pod Brokers")

	var artemisArray []*mgmt.Artemis
	var err error

	ssName, err := aa.GetStatefulSetName(request.Namespace)
	if err != nil {
		reqLogger.Error(err, "Failed to ge the statefulset name")
	}

	// Check to see if the statefulset already exists
	ssNamespacedName := types.NamespacedName{
		Name:      ssName + "-ss",
		Namespace: request.Namespace,
	}

	statefulset, err := ss.RetrieveStatefulSet(ssName, ssNamespacedName, client)
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
