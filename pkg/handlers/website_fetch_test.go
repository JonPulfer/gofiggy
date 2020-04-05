package handlers

import (
	"testing"

	api_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestParseAnnotationData(t *testing.T) {
	exampleAnnotationContent := "joke=curl-a-joke.herokuapp.com"

	fetchReq, err := parseAnnotationData(exampleAnnotationContent)
	if err != nil {
		t.Logf("failed to parse the annotation data: %s\n", err.Error())
		t.FailNow()
	}

	if fetchReq.IntoKey != "joke" {
		t.FailNow()
	}

	if fetchReq.FromSite.Host != "curl-a-joke.herokuapp.com" {
		t.FailNow()
	}
}

func TestStripNamespaceFromName(t *testing.T) {
	rawName := "default/the-name"
	gotName := stripNamespaceFromName(rawName)
	if gotName != "the-name" {
		t.FailNow()
	}

	if stripNamespaceFromName("short-name") != "short-name" {
		t.FailNow()
	}
}

func TestFetchSiteData(t *testing.T) {
	fRequest, _ := parseAnnotationData("joke=curl-a-joke.herokuapp.com")
	fResp, err := fetchSiteData(fRequest)
	if err != nil {
		t.Logf("received error from fetchSiteData: %s\n", err.Error())
		t.FailNow()
	}

	if fResp.Key != fRequest.IntoKey {
		t.FailNow()
	}

	if len(fResp.Value) == 0 {
		t.FailNow()
	}
	t.Logf("fResp: %s\n", fResp)
}

func TestFetchConfigMap(t *testing.T) {
	configMapToCreate := &api_v1.ConfigMap{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name: "simple-config",
			Annotations: map[string]string{
				CurlAnnotation: "joke=curl-a-joke.herokuapp.com",
			},
		},
		Data: map[string]string{},
	}

	kubeClient := fake.NewSimpleClientset()
	kubeClient.CoreV1().ConfigMaps("default").Create(configMapToCreate)
	configMap, err := fetchConfigMap(kubeClient, "default", "simple-config")
	if err != nil {
		t.Logf("error getting configMap: %s", err.Error())
		t.FailNow()
	}

	if configMap.Name != configMapToCreate.Name {
		t.Fail()
	}
}

func TestProcessConfigMap(t *testing.T) {
	configMapToCreate := &api_v1.ConfigMap{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name: "simple-config",
			Annotations: map[string]string{
				CurlAnnotation: "joke=curl-a-joke.herokuapp.com",
			},
		},
		Data: map[string]string{},
	}

	kubeClient := fake.NewSimpleClientset()
	kubeClient.CoreV1().ConfigMaps("default").Create(configMapToCreate)

	err := processConfigMap(kubeClient, "default", configMapToCreate)
	if err != nil {
		t.Logf("error processConfigMap: %s", err.Error())
		t.FailNow()
	}

	configMap, err := fetchConfigMap(kubeClient, "default", "simple-config")
	if err != nil {
		t.Logf("error getting configMap: %s", err.Error())
		t.FailNow()
	}

	joke := configMap.Data["joke"]
	if len(joke) == 0 {
		t.FailNow()
	}
	t.Logf("joke = %s", joke)
}
