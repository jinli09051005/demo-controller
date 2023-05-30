package reconciler

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func Start() {
	//scheme := runtime.NewScheme()
	//_ = corev1.AddToScheme(scheme)
	//_ = appsv1.AddToScheme(scheme)
	selectorLabel := make(map[string]string)
	selectorLabel["app"] = "demo"
	config := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(config,
		ctrl.Options{
			Scheme: scheme,
			NewCache: func(config *rest.Config, opts cache.Options) (cache.Cache, error) {
				return cache.New(config, cache.Options{
					ByObject: map[client.Object]cache.ByObject{
						&corev1.ConfigMap{}: {
							Label: labels.Set(selectorLabel).AsSelector(),
						},
						&appsv1.Deployment{}: {
							Label: labels.Set(selectorLabel).AsSelector(),
						},
					},
				})
			},
		},
	)
	if err != nil {
		fmt.Printf("%v,unable to init manager", err)
		os.Exit(1)
	}

	//c, _ := controller.New("app", mgr, controller.Options{})
	//_ = c.Watch(source.Kind(), &handler.EnqueueRequestForObject{}, predicate.Funcs{
	//	UpdateFunc: func(updateEvent event.UpdateEvent) bool {
	//		return true
	//	},
	//})

	err = ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).WithEventFilter(
		predicate.Funcs{
			CreateFunc: func(createEvent event.CreateEvent) bool {
				return false
			},
			DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
				return false
			},
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				return true
			},
			GenericFunc: func(genericEvent event.GenericEvent) bool {
				return false
			}}).
		Owns(&appsv1.Deployment{}).
		Complete(&reconcilerConfigMap{
			Client: mgr.GetClient(),
			scheme: mgr.GetScheme(),
		})
	if err != nil {
		fmt.Printf("%v,unable to create controller", err)
		os.Exit(1)
	}
	if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		fmt.Printf("%v,stop manager", err)
		os.Exit(1)
	}
}
