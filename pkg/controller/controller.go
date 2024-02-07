package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vidya-ranganathan/mcluster/pkg/apis/cumulonimbus.ai/v1alpha1"
	klientset "github.com/vidya-ranganathan/mcluster/pkg/client/clientset/versioned"
	"github.com/vidya-ranganathan/mcluster/pkg/client/clientset/versioned/scheme"
	kinf "github.com/vidya-ranganathan/mcluster/pkg/client/informers/externalversions/cumulonimbus.ai/v1alpha1"
	mclister "github.com/vidya-ranganathan/mcluster/pkg/client/listers/cumulonimbus.ai/v1alpha1"
	"github.com/vidya-ranganathan/mcluster/pkg/todo"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	// EVENT_IMPORTS_START
	customscheme "github.com/vidya-ranganathan/mcluster/pkg/client/clientset/versioned/scheme"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"

	// "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"

	// EVENT_IMPORTS_END

	// DELETE_IMPORT_START
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	// DELETE_IMPORT_END

	// STATUS_START
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// STATUS_END
)

type Controller struct {
	// clientset
	klient klientset.Interface
	// sync mcluster
	mclusterSynced cache.InformerSynced
	// queue to hold mcluster objects
	mclister mclister.MclusterLister
	// lister for mcluster
	wq workqueue.RateLimitingInterface

	// EVENT_START
	client   kubernetes.Interface
	recorder record.EventRecorder
	// EVENT_END
}

func NewController(client kubernetes.Interface,
	klient klientset.Interface,
	mclusterInformer kinf.MclusterInformer) *Controller {
	// Initialize the EVENT module ; add controller types to the scheme
	// And let events be recorded for the types

	runtime.Must(customscheme.AddToScheme(scheme.Scheme))

	// creating a new event broadcaster
	eveBroadCaster := record.NewBroadcaster()
	// setting up event logging mechanisms
	eveBroadCaster.StartStructuredLogging(0)

	eveBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{
		Interface: client.CoreV1().Events(""),
	})

	recorder := eveBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "Mcluster"})

	con := &Controller{
		client:         client,
		klient:         klient,
		mclusterSynced: mclusterInformer.Informer().HasSynced,
		mclister:       mclusterInformer.Lister(),
		wq:             workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "mcluster"),
		recorder:       recorder,
	}

	mclusterInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    con.handleAdd,
			DeleteFunc: con.handleDel,
		},
	)

	return con
}

func (con *Controller) handleAdd(obj interface{}) {
	fmt.Println("handleAdd was called...")
	con.wq.Add(obj)
}

func (con *Controller) handleDel(obj interface{}) {
	fmt.Println("handleDel was called")
	con.wq.Add(obj)
}

func (con *Controller) Run(ch chan struct{}) error {
	if ok := cache.WaitForCacheSync(ch, con.mclusterSynced); !ok {
		log.Println("cache was not sycned")
	}

	fmt.Println("worker#1..................")
	go wait.Until(con.worker, time.Second, ch)
	fmt.Println("worker#2..................")
	go wait.Until(con.worker, time.Second, ch)

	<-ch
	return nil
}

func (con *Controller) worker() {
	for con.processNextItem() {

	}
}

func (con *Controller) processNextItem() bool {
	// get the object from the workqueue
	item, shutDown := con.wq.Get()
	if shutDown {
		// logs as well
		return false
	}

	// after object is processed , remove it from the workqueue
	// so that it is not re-processed once more.
	defer con.wq.Forget(item)

	// based on the item , retrieve the key
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		log.Printf("error %s calling Namespace key func on cache for item\n", err.Error())
		return false
	}

	// fetch the name and namespace
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		log.Printf("spilitting key into the namespace and name, error - %s\n", err.Error())
		return false
	}

	// get the CRD and its specs.
	mcluster, err := con.mclister.Mclusters(ns).Get(name)
	// DELETE_START
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Delete the KIND cluster with given name in the specification..
			return todo.Delete(name)
		}
		log.Printf("error %s getting the mcluster resource from lister\n", err.Error())
		// EVENT_START - logging
		// con.recorder.Event(mcluster, corev1.EventTypeNormal, "KindClusterDelete", "KIND cluster deleted")
		// EVENT_END
		return false
	}
	// DELETE_END

	fmt.Printf("mcluster specs before performing the task is %+v\n", mcluster.Spec)

	// EVENT_START - logging
	con.recorder.Event(mcluster, corev1.EventTypeNormal, "KindClusterCreation", "Calling the KIND cluster creation custom REST API endpoint")
	// EVENT_END

	// perform the controller job here..
	clusterID := todo.Add(mcluster.Spec)

	// STATUS_START
	err = con.updateStatus(clusterID, mcluster)
	if err != nil {
		log.Printf("error %s updating clustus status after creation of cluster", err.Error())
	}
	// STATUS_END

	// EVENT_START - logging
	con.recorder.Event(mcluster, corev1.EventTypeNormal, "KindClusterCreationComplete", "KIND cluster created")
	// EVENT_END

	return true
}

/* v0.4 adding status */
func (con *Controller) updateStatus(clusterID string, mcluster *v1alpha1.Mcluster) error {
	// get the latest version of mcluster
	mc, err := con.klient.CumulonimbusV1alpha1().Mclusters(mcluster.Namespace).Get(context.Background(), mcluster.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// now update the status with clusterID
	mc.Status.MclusterID = clusterID
	_, err = con.klient.CumulonimbusV1alpha1().Mclusters(mcluster.Namespace).UpdateStatus(context.Background(), mc, metav1.UpdateOptions{})
	return err
}
