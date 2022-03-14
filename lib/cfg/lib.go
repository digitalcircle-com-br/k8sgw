package cfg

import (
	"context"
	"log"
	"os"

	"github.com/digitalcircle-com-br/k8sgw/lib/muxmgr"
	nanoapigorm "github.com/digitalcircle-com-br/nanoapi-gorm"
	nanoapiredis "github.com/digitalcircle-com-br/nanoapi-redis"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func Setup() error {

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	watcher, err := clientset.CoreV1().ConfigMaps("default").Watch(context.Background(), metav1.SingleObject(metav1.ObjectMeta{Name: "gw", Namespace: "default"}))
	if err != nil {
		return err
	}
	ch := watcher.ResultChan()

	for {
		evt, ok := <-ch

		if ok {
			updatedMap, ok := evt.Object.(*corev1.ConfigMap)
			if ok {
				switch evt.Type {
				case watch.Added:
					fallthrough
				case watch.Modified:
					log.Printf("Got new config: %s\n", updatedMap.Data["gw"])
					muxmgr.Update(updatedMap.Data["gw"])
					os.Setenv("REDIS", updatedMap.Data["redis"])
					os.Setenv("DSN", updatedMap.Data["dsn"])
					err := nanoapiredis.Setup()
					if err != nil {
						log.Printf("Error setting up redis: %s", err.Error())
					}
					err = nanoapigorm.Setup()
					if err != nil {
						log.Printf("Error setting up postgres: %s", err.Error())
					}

				}
			}
		} else {
			log.Printf("Config chan closed")
			return nil
		}
	}
	log.Printf("Left Config Loop")
	return nil
}
