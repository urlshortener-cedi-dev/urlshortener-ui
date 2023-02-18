# \Apiv1Api

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1ShortlinkGet**](Apiv1Api.md#ApiV1ShortlinkGet) | **Get** /api/v1/shortlink/ | list shortlinks
[**ApiV1ShortlinkShortlinkDelete**](Apiv1Api.md#ApiV1ShortlinkShortlinkDelete) | **Delete** /api/v1/shortlink/{shortlink} | delete shortlink
[**ApiV1ShortlinkShortlinkGet**](Apiv1Api.md#ApiV1ShortlinkShortlinkGet) | **Get** /api/v1/shortlink/{shortlink} | get a shortlink
[**ApiV1ShortlinkShortlinkPost**](Apiv1Api.md#ApiV1ShortlinkShortlinkPost) | **Post** /api/v1/shortlink/{shortlink} | create new shortlink
[**ApiV1ShortlinkShortlinkPut**](Apiv1Api.md#ApiV1ShortlinkShortlinkPut) | **Put** /api/v1/shortlink/{shortlink} | update existing shortlink


# **ApiV1ShortlinkGet**
> []ControllerShortLink ApiV1ShortlinkGet(ctx, )
list shortlinks

list shortlinks

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]ControllerShortLink**](controller.ShortLink.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1ShortlinkShortlinkDelete**
> int32 ApiV1ShortlinkShortlinkDelete(ctx, shortlink)
delete shortlink

delete shortlink

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **shortlink** | **string**| the shortlink URL part (shortlink id) | 

### Return type

**int32**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1ShortlinkShortlinkGet**
> ControllerShortLink ApiV1ShortlinkShortlinkGet(ctx, shortlink)
get a shortlink

get a shortlink

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **shortlink** | **string**| the shortlink URL part (shortlink id) | 

### Return type

[**ControllerShortLink**](controller.ShortLink.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1ShortlinkShortlinkPost**
> int32 ApiV1ShortlinkShortlinkPost(ctx, shortlink, spec)
create new shortlink

create a new shortlink

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **shortlink** | **string**| the shortlink URL part (shortlink id) | 
  **spec** | [**V1alpha1ShortLinkSpec**](V1alpha1ShortLinkSpec.md)| shortlink spec | 

### Return type

**int32**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1ShortlinkShortlinkPut**
> int32 ApiV1ShortlinkShortlinkPut(ctx, shortlink, spec)
update existing shortlink

update a new shortlink

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **shortlink** | **string**| the shortlink URL part (shortlink id) | 
  **spec** | [**V1alpha1ShortLinkSpec**](V1alpha1ShortLinkSpec.md)| shortlink spec | 

### Return type

**int32**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

