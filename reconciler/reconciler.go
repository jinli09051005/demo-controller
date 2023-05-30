package reconciler

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	fmt.Println("开始调谐...")
	namespacedName := types.NamespacedName{
		Name:      "demo",
		Namespace: "default",
	}

	var configmap corev1.ConfigMap
	if err := r.Get(ctx, req.NamespacedName, &configmap); err != nil {
		fmt.Printf("%v,unable to get configmap", err)
		return ctrl.Result{}, err
	}

	var deploy appsv1.Deployment
	if err := r.Get(ctx, namespacedName, &deploy); err != nil {
		fmt.Printf("%v,unable to get deploy", err)
		return ctrl.Result{}, err
	}

	//if len(deploy.ObjectMeta.OwnerReferences) == 0 {
	//	if err := ctrl.SetControllerReference(&configmap, &deploy, r.scheme); err != nil {
	//		fmt.Printf("%v,unable to set deployment's ownerreference")
	//		return reconcile.Result{}, err
	//	}
	//}

	annotations := make(map[string]string)
	currentTime := time.Now()
	annotations["configmap-update-last-time"] = currentTime.Format("2006-01-02 15:04:05")
	deploy.Spec.Template.ObjectMeta.Annotations = annotations

	if err := r.Update(ctx, &deploy); err != nil {
		fmt.Printf("%v,unable to update deploy", err)
		return reconcile.Result{}, err
	}

	fmt.Println("..........................")

	fmt.Println("本周期调谐结束，进入下一周期...")

	return ctrl.Result{}, nil
}
