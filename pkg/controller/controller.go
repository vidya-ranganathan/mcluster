package controller

import (
	"fmt"
	"log"
	"time"

	klientset "github.com/vidya-ranganathan/mcluster/pkg/client/clientset/versioned"
	kinf "github.com/vidya-ranganathan/mcluster/pkg/client/informers/externalversions/cumulonimbus.ai/v1alpha1"
	mclister "github.com/vidya-ranganathan/mcluster/pkg/client/listers/cumulonimbus.ai/v1alpha1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
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
}

func NewController(klient klientset.Interface, mclusterInformer kinf.MclusterInformer) *Controller {
	con := &Controller{
		klient:         klient,
		mclusterSynced: mclusterInformer.Informer().HasSynced,
		mclister:       mclusterInformer.Lister(),
		wq:             workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "mcluster"),
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
	if err != nil {
		log.Printf("error %s getting the mcluster resource from lister\n", err.Error())
		return false
	}

	fmt.Printf("mcluster specs before performing the task is %+v\n", mcluster.Spec)

	return true
}
