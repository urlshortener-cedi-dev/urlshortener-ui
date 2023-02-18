# V1alpha1ShortLinkSpec

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**After** | **int32** | RedirectAfter specifies after how many seconds to redirect (Default&#x3D;3) +kubebuilder:default:&#x3D;0 +kubebuilder:validation:Minimum&#x3D;0 +kubebuilder:validation:Maximum&#x3D;99 | [optional] [default to null]
**Code** | **int32** | Code is the URL Code used for the redirection. leave on default (307) when using the HTML behavior. However, if you whish to use a HTTP 3xx redirect, set to the appropriate 3xx status code +kubebuilder:validation:Enum&#x3D;200;300;301;302;303;304;305;307;308 +kubebuilder:default:&#x3D;307 | [optional] [default to null]
**Owner** | **string** | Owner is the GitHub user name which created the shortlink +kubebuilder:validation:Required | [optional] [default to null]
**Owners** | **[]string** | Co-Owners are the GitHub user name which can also administrate this shortlink +kubebuilder:validation:Optional | [optional] [default to null]
**Target** | **string** | Target specifies the target to which we will redirect +kubebuilder:validation:Required +kubebuilder:validation:MinLength&#x3D;1 | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


