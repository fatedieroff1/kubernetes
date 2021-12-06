package storage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
)

// EncryptionScopesClient is the the Azure Storage Management API.
type EncryptionScopesClient struct {
	BaseClient
}

// NewEncryptionScopesClient creates an instance of the EncryptionScopesClient client.
func NewEncryptionScopesClient(subscriptionID string) EncryptionScopesClient {
	return NewEncryptionScopesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewEncryptionScopesClientWithBaseURI creates an instance of the EncryptionScopesClient client using a custom
// endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure
// stack).
func NewEncryptionScopesClientWithBaseURI(baseURI string, subscriptionID string) EncryptionScopesClient {
	return EncryptionScopesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get returns the properties for the specified encryption scope.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription. The name is case
// insensitive.
// accountName - the name of the storage account within the specified resource group. Storage account names
// must be between 3 and 24 characters in length and use numbers and lower-case letters only.
// encryptionScopeName - the name of the encryption scope within the specified storage account. Encryption
// scope names must be between 3 and 63 characters in length and use numbers, lower-case letters and dash (-)
// only. Every dash (-) character must be immediately preceded and followed by a letter or number.
func (client EncryptionScopesClient) Get(ctx context.Context, resourceGroupName string, accountName string, encryptionScopeName string) (result EncryptionScope, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/EncryptionScopesClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{
			TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{
				{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil},
			},
		},
		{
			TargetValue: accountName,
			Constraints: []validation.Constraint{
				{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil},
			},
		},
		{
			TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}},
		},
		{
			TargetValue: encryptionScopeName,
			Constraints: []validation.Constraint{
				{Target: "encryptionScopeName", Name: validation.MaxLength, Rule: 63, Chain: nil},
				{Target: "encryptionScopeName", Name: validation.MinLength, Rule: 3, Chain: nil},
			},
		},
	}); err != nil {
		return result, validation.NewError("storage.EncryptionScopesClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, accountName, encryptionScopeName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client EncryptionScopesClient) GetPreparer(ctx context.Context, resourceGroupName string, accountName string, encryptionScopeName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":         autorest.Encode("path", accountName),
		"encryptionScopeName": autorest.Encode("path", encryptionScopeName),
		"resourceGroupName":   autorest.Encode("path", resourceGroupName),
		"subscriptionId":      autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/encryptionScopes/{encryptionScopeName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client EncryptionScopesClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client EncryptionScopesClient) GetResponder(resp *http.Response) (result EncryptionScope, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List lists all the encryption scopes available under the specified storage account.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription. The name is case
// insensitive.
// accountName - the name of the storage account within the specified resource group. Storage account names
// must be between 3 and 24 characters in length and use numbers and lower-case letters only.
func (client EncryptionScopesClient) List(ctx context.Context, resourceGroupName string, accountName string) (result EncryptionScopeListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/EncryptionScopesClient.List")
		defer func() {
			sc := -1
			if result.eslr.Response.Response != nil {
				sc = result.eslr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{
			TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{
				{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil},
			},
		},
		{
			TargetValue: accountName,
			Constraints: []validation.Constraint{
				{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil},
			},
		},
		{
			TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}},
		},
	}); err != nil {
		return result, validation.NewError("storage.EncryptionScopesClient", "List", err.Error())
	}

	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName, accountName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.eslr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "List", resp, "Failure sending request")
		return
	}

	result.eslr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "List", resp, "Failure responding to request")
		return
	}
	if result.eslr.hasNextLink() && result.eslr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client EncryptionScopesClient) ListPreparer(ctx context.Context, resourceGroupName string, accountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/encryptionScopes", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client EncryptionScopesClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client EncryptionScopesClient) ListResponder(resp *http.Response) (result EncryptionScopeListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client EncryptionScopesClient) listNextResults(ctx context.Context, lastResults EncryptionScopeListResult) (result EncryptionScopeListResult, err error) {
	req, err := lastResults.encryptionScopeListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client EncryptionScopesClient) ListComplete(ctx context.Context, resourceGroupName string, accountName string) (result EncryptionScopeListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/EncryptionScopesClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, resourceGroupName, accountName)
	return
}

// Patch update encryption scope properties as specified in the request body. Update fails if the specified encryption
// scope does not already exist.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription. The name is case
// insensitive.
// accountName - the name of the storage account within the specified resource group. Storage account names
// must be between 3 and 24 characters in length and use numbers and lower-case letters only.
// encryptionScopeName - the name of the encryption scope within the specified storage account. Encryption
// scope names must be between 3 and 63 characters in length and use numbers, lower-case letters and dash (-)
// only. Every dash (-) character must be immediately preceded and followed by a letter or number.
// encryptionScope - encryption scope properties to be used for the update.
func (client EncryptionScopesClient) Patch(ctx context.Context, resourceGroupName string, accountName string, encryptionScopeName string, encryptionScope EncryptionScope) (result EncryptionScope, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/EncryptionScopesClient.Patch")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{
			TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{
				{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil},
			},
		},
		{
			TargetValue: accountName,
			Constraints: []validation.Constraint{
				{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil},
			},
		},
		{
			TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}},
		},
		{
			TargetValue: encryptionScopeName,
			Constraints: []validation.Constraint{
				{Target: "encryptionScopeName", Name: validation.MaxLength, Rule: 63, Chain: nil},
				{Target: "encryptionScopeName", Name: validation.MinLength, Rule: 3, Chain: nil},
			},
		},
	}); err != nil {
		return result, validation.NewError("storage.EncryptionScopesClient", "Patch", err.Error())
	}

	req, err := client.PatchPreparer(ctx, resourceGroupName, accountName, encryptionScopeName, encryptionScope)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "Patch", nil, "Failure preparing request")
		return
	}

	resp, err := client.PatchSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "Patch", resp, "Failure sending request")
		return
	}

	result, err = client.PatchResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "Patch", resp, "Failure responding to request")
		return
	}

	return
}

// PatchPreparer prepares the Patch request.
func (client EncryptionScopesClient) PatchPreparer(ctx context.Context, resourceGroupName string, accountName string, encryptionScopeName string, encryptionScope EncryptionScope) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":         autorest.Encode("path", accountName),
		"encryptionScopeName": autorest.Encode("path", encryptionScopeName),
		"resourceGroupName":   autorest.Encode("path", resourceGroupName),
		"subscriptionId":      autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/encryptionScopes/{encryptionScopeName}", pathParameters),
		autorest.WithJSON(encryptionScope),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PatchSender sends the Patch request. The method will close the
// http.Response Body if it receives an error.
func (client EncryptionScopesClient) PatchSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// PatchResponder handles the response to the Patch request. The method always
// closes the http.Response Body.
func (client EncryptionScopesClient) PatchResponder(resp *http.Response) (result EncryptionScope, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Put synchronously creates or updates an encryption scope under the specified storage account. If an encryption scope
// is already created and a subsequent request is issued with different properties, the encryption scope properties
// will be updated per the specified request.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription. The name is case
// insensitive.
// accountName - the name of the storage account within the specified resource group. Storage account names
// must be between 3 and 24 characters in length and use numbers and lower-case letters only.
// encryptionScopeName - the name of the encryption scope within the specified storage account. Encryption
// scope names must be between 3 and 63 characters in length and use numbers, lower-case letters and dash (-)
// only. Every dash (-) character must be immediately preceded and followed by a letter or number.
// encryptionScope - encryption scope properties to be used for the create or update.
func (client EncryptionScopesClient) Put(ctx context.Context, resourceGroupName string, accountName string, encryptionScopeName string, encryptionScope EncryptionScope) (result EncryptionScope, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/EncryptionScopesClient.Put")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{
			TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{
				{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil},
			},
		},
		{
			TargetValue: accountName,
			Constraints: []validation.Constraint{
				{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil},
			},
		},
		{
			TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}},
		},
		{
			TargetValue: encryptionScopeName,
			Constraints: []validation.Constraint{
				{Target: "encryptionScopeName", Name: validation.MaxLength, Rule: 63, Chain: nil},
				{Target: "encryptionScopeName", Name: validation.MinLength, Rule: 3, Chain: nil},
			},
		},
	}); err != nil {
		return result, validation.NewError("storage.EncryptionScopesClient", "Put", err.Error())
	}

	req, err := client.PutPreparer(ctx, resourceGroupName, accountName, encryptionScopeName, encryptionScope)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "Put", nil, "Failure preparing request")
		return
	}

	resp, err := client.PutSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "Put", resp, "Failure sending request")
		return
	}

	result, err = client.PutResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.EncryptionScopesClient", "Put", resp, "Failure responding to request")
		return
	}

	return
}

// PutPreparer prepares the Put request.
func (client EncryptionScopesClient) PutPreparer(ctx context.Context, resourceGroupName string, accountName string, encryptionScopeName string, encryptionScope EncryptionScope) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":         autorest.Encode("path", accountName),
		"encryptionScopeName": autorest.Encode("path", encryptionScopeName),
		"resourceGroupName":   autorest.Encode("path", resourceGroupName),
		"subscriptionId":      autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/encryptionScopes/{encryptionScopeName}", pathParameters),
		autorest.WithJSON(encryptionScope),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PutSender sends the Put request. The method will close the
// http.Response Body if it receives an error.
func (client EncryptionScopesClient) PutSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// PutResponder handles the response to the Put request. The method always
// closes the http.Response Body.
func (client EncryptionScopesClient) PutResponder(resp *http.Response) (result EncryptionScope, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
