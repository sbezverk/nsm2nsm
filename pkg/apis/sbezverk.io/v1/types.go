package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServerEndpoint CRD
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ServerEndpoint struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            ServerEndpointSpec `json:"spec"`
}

// ServerEndpointSpec describes Server's Endpoint related info
type ServerEndpointSpec struct {
	ServerAddress string `json:"serverAddress"`
}

// ServerEndpointList is the list schema for this CRD
// -genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ServerEndpointList struct {
	meta.TypeMeta `json:",inline"`
	// +optional
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []ServerEndpoint `json:"items"`
}
