/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// KrknWebhookSpec defines the desired state of KrknWebhook.
// A webhook posts scenario lifecycle events to an external endpoint, so a chaos
// run can notify a dashboard or a chat channel without polling the API.
type KrknWebhookSpec struct {
	// URL is the HTTPS endpoint that receives the event payload
	URL string `json:"url"`

	// Event selects which scenario lifecycle event triggers the post
	// +optional
	// +kubebuilder:default=completed
	// +kubebuilder:validation:Enum=started;completed;failed
	Event string `json:"event,omitempty"`

	// SigningToken is the shared secret used to sign the payload body,
	// sent as the X-Krkn-Signature header
	// +optional
	SigningToken string `json:"signingToken,omitempty"`

	// TokenSecretRef names a Secret holding the signing token, preferred over
	// SigningToken so the value never appears in the resource
	// +optional
	TokenSecretRef string `json:"tokenSecretRef,omitempty"`

	// RetryLimit is how many times a failed post is retried before the event
	// is dropped. Set it to 0 to disable retries, so the first failure drops
	// the event immediately.
	// +optional
	// +kubebuilder:default=3
	RetryLimit int32 `json:"retryLimit,omitempty"`
}

// KrknWebhookStatus defines the observed state of KrknWebhook.
type KrknWebhookStatus struct {
	// LastDelivery is when the most recent post was accepted by the endpoint
	// +optional
	LastDelivery *metav1.Time `json:"lastDelivery,omitempty"`

	// Failures counts consecutive failed deliveries, reset on success
	// +optional
	Failures int32 `json:"failures,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Event",type=string,JSONPath=`.spec.event`
// +kubebuilder:printcolumn:name="Failures",type=integer,JSONPath=`.status.failures`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:shortName=kwh

// KrknWebhook is the Schema for the krknwebhooks API.
// It registers an external endpoint to notify when a scenario run changes state.
type KrknWebhook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KrknWebhookSpec   `json:"spec,omitempty"`
	Status KrknWebhookStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KrknWebhookList contains a list of KrknWebhook.
type KrknWebhookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KrknWebhook `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(GroupVersion, &KrknWebhook{}, &KrknWebhookList{})
		return nil
	})
}
