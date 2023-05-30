package reconciler

import (
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Scheme负责各个资源版本的统一注册和管理，
// 为其他组件提供根据GVK来获取对应的资源对象(GVK与go type转换)，
// 也提供通过资源对象获取版本等操作(go type与GVK转换)，
// 内部还持有convert对象是包装了源对象到目标对象的转换操作(不同资源版本之间的转换)

// GVR GVK

var scheme = runtime.NewScheme()

func init() {
	log.SetLogger(zap.New())
	//clientgoscheme包含了除crd资源外的所有种类资源
	_ = clientgoscheme.AddToScheme(scheme)
}
