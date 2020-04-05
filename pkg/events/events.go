package events

import (
	api_v1 "k8s.io/api/core/v1"
)

type EventHandler interface {
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(oldObj, newObj interface{})
}

// Event received from Kubernetes from the watcher.
type Event struct {
	Namespace string
	Kind      string
	Component string
	Host      string
	Reason    string
	Status    string
	Name      string
}

func New(obj interface{}, action string) Event {

	var name, kind, namespace string

	switch object := obj.(type) {
	case *api_v1.ConfigMap:
		kind = "configmap"
	case Event:
		name = object.Name
		kind = object.Kind
		namespace = object.Namespace
	}

	return Event {
		Name: name,
		Namespace: namespace,
		Kind: kind,
		Reason: action,
		Status: convertActionToStatus(action),
	}
}

func convertActionToStatus(action string) string {
	switch action {
	case "created":
		return "Normal"
	case "deleted":
		return "Danger"
	case "updated":
		return "Warning"
	}
	return ""
}


