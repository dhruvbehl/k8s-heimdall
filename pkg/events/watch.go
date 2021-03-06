package events

import (
	"time"

	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var EventList []v1.Event

func WatchEvents(clientset *kubernetes.Clientset) {
	eventListWatch := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(),
	string("Events"), v1.NamespaceAll, fields.Everything())

	_, controller := cache.NewInformer(eventListWatch, &v1.Event{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) {
			event := obj.(*v1.Event)
			EventList = append(EventList, *event)
			log.Info().Msgf("new event received; operation: [%v] kind: [%v] name: [%v] reason: [%v] message: [%v]",
			 event.ManagedFields[0].Operation, event.InvolvedObject.Kind, event.InvolvedObject.Name, event.Reason, event.Message)
		},
	})

	stop := make(chan struct{})
	defer close(stop)

	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}