package main

import (
	//"context"
	//"flag"
	//"fmt"
	//"log"
	//"path/filepath"

	//klient "github.com/vidya-ranganathan/mcluster/pkg/client/clientset/versioned"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/client-go/rest"
	//"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/util/homedir"
	"fmt"

	"github.com/vidya-ranganathan/mcluster/pkg/apis/cumulonimbus.ai/v1alpha1"
)

func main() {
	/*
		var kubeconfig *string

		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to kubeconfig")
		}
		flag.Parse()

		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Printf("Building config from flags, %s", err.Error())
			config, err = rest.InClusterConfig()
			if err != nil {
				log.Printf("error %s building inclusterconfig", err.Error())
			}
		}
		klientset, err := klient.NewForConfig(config)
		if err != nil {
			log.Printf("getting klient set %s\n", err.Error())
		}

		klusters, err := klientset.CumulonimbusV1alpha1().Mclusters("").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			log.Printf("listing klusters - %s\n", err.Error())
		}

		fmt.Printf("listing klusters %d\n", len(klusters.Items))
	*/

	/* step #1 */
	mc := v1alpha1.Mcluster{}
	fmt.Println(mc)

}
