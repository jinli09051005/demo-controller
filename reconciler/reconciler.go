package reconciler

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

type reconcilerConfigMap struct {
	client.Client
	scheme *runtime.Scheme
}

func (r *reconcilerConfigMap) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	if err := r.Get(ctx, req.NamespacedName, &corev1.ConfigMap{}); err != nil {
		fmt.Printf("%v,unable to get configmap", err)
		return ctrl.Result{}, err
	}

	var deploy appsv1.Deployment
	labels := make(map[string]string)
	labels["app"] = "demo"
	deploy.Labels = labels
	if err := r.Get(ctx, req.NamespacedName, &deploy); err != nil {
		fmt.Printf("%v,unable to get deploy", err)
		return ctrl.Result{}, err
	}

	if deploy.ObjectMeta.Name != "demo" {

	}
	annotations := make(map[string]string)
	currentTime := time.Now()
	annotations["configmap-update-last-time"] = currentTime.Format("2006-01-02 15:04:05")
	deploy.Spec.Template.ObjectMeta.Annotations = annotations
	if err := r.Update(ctx, &deploy); err != nil {
		fmt.Printf("%v,unable to update deploy")
		return reconcile.Result{}, err
	}
	return ctrl.Result{}, nil
}
