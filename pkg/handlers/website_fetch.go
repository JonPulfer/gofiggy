package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	api_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/JonPulfer/gofiggy/pkg/events"
	"github.com/JonPulfer/gofiggy/pkg/utils"
)

const CurlAnnotation = "x-k8s.io/curl-me-that"

type WebsiteFetchHandler struct {
	logger    zerolog.Logger
	clientset kubernetes.Interface
}

// WebsiteFetchHandler watches for creation and updates to configmaps to see
// whether the annotation `x-k8s.io/curl-me-that:` has appeared. We parse the
// data from the annotation using `parseAnnotationData()` to extract the site
// and key. We then fetch the content from the site and update the config map
// setting the key and adding the site content as the data.
func NewWebsiteFetchHandler() WebsiteFetchHandler {
	return WebsiteFetchHandler{
		logger:    zerolog.New(os.Stderr).With().Timestamp().Logger(),
		clientset: utils.GetClient(),
	}
}

func (wfh WebsiteFetchHandler) ObjectCreated(obj interface{}) {
	ev := events.New(obj, "created")
	wfh.logger.Log().Fields(map[string]interface{}{"event": ev}).
		Msg("received created event")

	wfh.logger.Log().Fields(map[string]interface{}{"event": ev}).
		Msg("fetching config map")

	configMap, err := fetchConfigMap(wfh.clientset, "default",
		stripNamespaceFromName(ev.Name))
	if err != nil {
		wfh.logger.Log().Msg(err.Error())
	}
	wfh.logger.Log().Fields(map[string]interface{}{"configMaps": configMap}).
		Msg("response from fetchConfigMap")

	if err := processConfigMap(wfh.clientset, "default", configMap); err != nil {
		wfh.logger.Log().Err(err).
			Msg("failed to process the created configMap")
	}
}

func (wfh WebsiteFetchHandler) ObjectDeleted(obj interface{}) {
	ev := events.New(obj, "deleted")
	wfh.logger.Log().Fields(map[string]interface{}{"event": ev}).
		Msg("received deleted event")
}

func (wfh WebsiteFetchHandler) ObjectUpdated(oldObj interface{}, newObj interface{}) {
	ev := events.New(newObj, "updated")
	wfh.logger.Log().Fields(map[string]interface{}{"event": ev}).
		Msg("received updated event")

	wfh.logger.Log().Fields(map[string]interface{}{"event": ev}).
		Msg("fetching config map")

	configMap, err := fetchConfigMap(wfh.clientset, "default",
		stripNamespaceFromName(ev.Name))
	if err != nil {
		wfh.logger.Log().Msg(err.Error())
	}
	wfh.logger.Log().Fields(map[string]interface{}{"configMaps": configMap}).
		Msg("response from fetchConfigMap")

	if err := processConfigMap(wfh.clientset, "default", configMap); err != nil {
		wfh.logger.Log().Err(err).
			Msg("failed to process the updated configMap")
	}
}

// FetchRequest holds the URL of the site we want to fetch the content from and
// the key name to add the content to the config map with.
type FetchRequest struct {
	IntoKey  string
	FromSite *url.URL
}

// parseAnnotationData into a FetchRequest. We parse the content of the
// annotation which is expected to look like: -
//
//	joke=curl-a-joke.herokuapp.com
//
// From this we would convert `curl-a-joke.herokuapp.com` into a url.URL and
// set it as the FromSite. We would take `joke` and set that as the IntoKey.
func parseAnnotationData(annotationData string) (*FetchRequest, error) {
	parts := strings.Split(annotationData, "=")
	if len(parts) != 2 {
		return nil, errors.New(
			fmt.Sprintf(
				"unexpected value provided for annotationData: %s",
				annotationData))
	}

	withScheme := parts[1]
	if !strings.Contains(withScheme, "http") {
		withScheme = "http://" + withScheme
	}

	websiteURL, err := url.Parse(withScheme)
	if err != nil {
		return nil, err
	}

	return &FetchRequest{IntoKey: parts[0], FromSite: websiteURL}, nil
}

// FetchResponse provides the content that will be placed in the config map that
// requested it via the annotation.
type FetchResponse struct {
	Key   string
	Value string
}

func (fr FetchResponse) String() string {
	return fmt.Sprintf("%s=%s", fr.Key, fr.Value)
}

// fetchSiteData makes a simple http GET request to fetch data from the site in
// the provided FetchRequest. The result holds the Key and site data as the
// value.
func fetchSiteData(fRequest *FetchRequest) (*FetchResponse, error) {
	cl := http.Client{}
	resp, err := cl.Get(fRequest.FromSite.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("received %d status from fetch",
				resp.StatusCode))
	}

	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)

	return &FetchResponse{
		Key:   fRequest.IntoKey,
		Value: buf.String(),
	}, nil
}

// fetchConfigMap using a kubernetes client.
func fetchConfigMap(kubeClient kubernetes.Interface, namespace string,
	configMapName string) (*api_v1.ConfigMap, error) {

	configMap, err := kubeClient.CoreV1().ConfigMaps(namespace).
		Get(configMapName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return configMap, nil
}

// processConfigMap to see whether it has the appropriate annotation. Extract
// the site data request from the annotation and then add the data field with
// the request key.
func processConfigMap(
	kubeClient kubernetes.Interface,
	namespace string,
	configMap *api_v1.ConfigMap) error {
	if configMap != nil {
		if configMapHasAnnotation(configMap) {
			fReq, err := parseAnnotationData(configMap.Annotations[CurlAnnotation])
			if err != nil {
				return err
			}
			fResp, err := fetchSiteData(fReq)
			if err != nil {
				return err
			}
			configMap.Data[fResp.Key] = fResp.Value
			if err := updateConfigMap(kubeClient, namespace, configMap); err != nil {
				return err
			}
		}
	}

	return nil
}

// updateConfigMap applies the changed configMap to the namespace.
func updateConfigMap(kubeClient kubernetes.Interface, namespace string,
	configMap *api_v1.ConfigMap) error {

	_, err := kubeClient.CoreV1().ConfigMaps(namespace).Update(configMap)
	return err
}

// configMapHasAnnotation indicates whether this configMap has the annotation
// we are looking for.
func configMapHasAnnotation(configMap *api_v1.ConfigMap) bool {
	if annotation := configMap.Annotations[CurlAnnotation]; len(annotation) != 0 {
		return true
	}
	return false
}

// stripNamespaceFromName when received from an event the resource name includes
// the namespace with a forward slash separator.
func stripNamespaceFromName(raw string) string {
	parts := strings.Split(raw, "/")
	if len(parts) != 2 {
		return raw
	}
	return parts[1]
}
