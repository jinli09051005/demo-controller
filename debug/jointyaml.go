package debug

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	yamlCodec "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
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
	dec := yaml.NewDecoder(file)
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
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discache.NewMemCacheClient(discoveryClient))
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return err
	}

	obj := &unstructured.Unstructured{}
	_, gvk, err := codec.Decode(data, nil, obj)
	if err != nil {
		return err
	}

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	var dynamicResource dynamic.ResourceInterface = dynamicClient.Resource(mapping.Resource)
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		namesapce := obj.GetNamespace()
		if namesapce == "" {
			namesapce = "default"
		}
		dynamicResource = dynamicClient.Resource(mapping.Resource).Namespace(namesapce)
	}

	if _, err = dynamicResource.Create(ctx, obj, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func CreateJointYml() {
	ctx := context.TODO()
	codec := yamlCodec.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	file, err := os.Open("./debug/joint.yaml")
	if err != nil {
		panic(err)
	}
	dec := yaml.NewDecoder(file)
	ymls := make(map[string]interface{})
	err = dec.Decode(&ymls)
	for err == nil {
		subyml, _ := json.Marshal(ymls)
		if err = dynamicCreate(ctx, restConfig, codec, subyml); err != nil {
			panic(err)
		}
		err = dec.Decode(&ymls)
	}
	if !errors.Is(err, io.EOF) {
		fmt.Println(err)
	}
	fmt.Println("All yml apply complete!")
}
