// @author : vidya.ranganathan@cumulonimbus.ai

package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"

	klient "github.com/vidya-ranganathan/mcluster/pkg/client/clientset/versioned"
	kInfFac "github.com/vidya-ranganathan/mcluster/pkg/client/informers/externalversions"
	"github.com/vidya-ranganathan/mcluster/pkg/controller"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to kubeconfig")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Printf("Building config from flags, %s - next trying to build with InClusterConfig()", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("error %s building inclusterconfig", err.Error())
		}
	}

	klientset, err := klient.NewForConfig(config)
	if err != nil {
		log.Printf("getting klient set %s\n", err.Error())
	}

	/* step#1
	mclusters, err := klientset.CumulonimbusV1alpha1().Mclusters("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("listing klusters - %s\n", err.Error())
	}

	fmt.Printf("listing klusters %d\n", len(mclusters.Items))
	*/

	infoFactory := kInfFac.NewSharedInformerFactory(klientset, 20*time.Minute)
	ch := make(chan struct{})
	con := controller.NewController(klientset, infoFactory.Cumulonimbus().V1alpha1().Mclusters())
	infoFactory.Start(ch)
	if err := con.Run(ch); err != nil {
		log.Printf("error running controller %s\n", err.Error())
	}

}
