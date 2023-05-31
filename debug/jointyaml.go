package debug

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	goyml "gopkg.in/yaml.v3"
	"io"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ymlcodec "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	discache "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"os"
)

func GetJointYaml() {
	file, err := os.Open("./debug/joint.yaml")
	if err != nil {
		fmt.Println(err)
	}
	dec := goyml.NewDecoder(file)
	y := make(map[string]interface{})

	err = dec.Decode(&y)
	for err == nil {
		dataType, _ := json.Marshal(y)
		dataString := string(dataType)
		fmt.Println(dataString)
		err = dec.Decode(&y)
	}
	if !errors.Is(err, io.EOF) {
		fmt.Println(err)
	}
	fmt.Println("Parse yaml complete!")
}

func dynamicCreate(ctx context.Context, cfg *rest.Config, codec runtime.Serializer, data []byte) error {
	// 创建discoveryClient
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return err
	}
	// 创建dynamicClient
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return err
	}
	// 创建RESTMapper接口实例
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discache.NewMemCacheClient(discoveryClient))

	obj := &unstructured.Unstructured{}
	// 解码yml，获取gvk和Object
	_, gvk, err := codec.Decode(data, nil, obj)
	if err != nil {
		return err
	}
	// 获取gvr
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}
	gvr := mapping.Resource

	var dynamicResource dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		namesapce := obj.GetNamespace()
		if namesapce == "" {
			namesapce = "default"
		}
		dynamicResource = dynamicClient.Resource(gvr).Namespace(namesapce)
	}
	// 创建k8s资源
	if _, err = dynamicResource.Create(ctx, obj, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func CreateJointYml() {
	ctx := context.TODO()
	// 初始化编解码器
	codec := ymlcodec.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	// 从容器中获取集群配置
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	// 读取联合yml
	file, err := os.Open("./debug/joint.yaml")
	if err != nil {
		panic(err)
	}
	dec := goyml.NewDecoder(file)
	ymls := make(map[string]interface{})
	err = dec.Decode(&ymls)
	for err == nil {
		subyml, _ := json.Marshal(ymls)
		// 创建k8s资源对象
		if err = dynamicCreate(ctx, restConfig, codec, subyml); err != nil {
			panic(err)
		}
		err = dec.Decode(&ymls)
	}
	// 读取完毕
	if !errors.Is(err, io.EOF) {
		fmt.Println(err)
	}
	fmt.Println("All yml apply complete!")
}
