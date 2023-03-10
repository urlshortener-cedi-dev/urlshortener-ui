/*
 * URL Shortener
 *
 * A url shortener, written in Go running on Kubernetes
 *
 * API version: 1.0
 * Contact: urlshortener-api@cedi.dev
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type V1alpha1ShortLinkStatus struct {
	// ChangedBy indicates who (GitHub User) changed the Shortlink last +kubebuilder:validation:Optional
	Changedby string `json:"changedby,omitempty"`
	// Count represents how often this ShortLink has been called +kubebuilder:default:=0 +kubebuilder:validation:Minimum=0
	Count int32 `json:"count,omitempty"`
	// LastModified is a date-time when the ShortLink was last modified +kubebuilder:validation:Format:date-time +kubebuilder:validation:Optional
	Lastmodified string `json:"lastmodified,omitempty"`
}
