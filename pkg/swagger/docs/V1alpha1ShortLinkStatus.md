# V1alpha1ShortLinkStatus

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Changedby** | **string** | ChangedBy indicates who (GitHub User) changed the Shortlink last +kubebuilder:validation:Optional | [optional] [default to null]
**Count** | **int32** | Count represents how often this ShortLink has been called +kubebuilder:default:&#x3D;0 +kubebuilder:validation:Minimum&#x3D;0 | [optional] [default to null]
**Lastmodified** | **string** | LastModified is a date-time when the ShortLink was last modified +kubebuilder:validation:Format:date-time +kubebuilder:validation:Optional | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


