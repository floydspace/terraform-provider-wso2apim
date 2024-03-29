/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package apim handles the interactions with APIM.
package apim

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/pkg/errors"
	"github.com/wso2/openservicebroker-apim/pkg/client"
	"github.com/wso2/openservicebroker-apim/pkg/log"
	"github.com/wso2/openservicebroker-apim/pkg/token"
	"github.com/wso2/openservicebroker-apim/pkg/utils"
)

const (
	CreateAPIContext                  = "create API"
	CreateApplicationContext          = "create application"
	CreateSubscriptionContext         = "create subscription"
	CreateMultipleSubscriptionContext = "create multiple subscriptions"
	ChangeAPILifeCycleContext         = "change API lifecycle"
	UpdateAPIContext                  = "update API"
	UpdateApplicationContext          = "update application"
	UpdateSubscriptionContext         = "update subscription"
	GenerateKeyContext                = "Generate application keys"
	RegenerateKeyContext              = "Regenerate application keys"
	CleanupKeysContext                = "cleanup application keys"
	UnSubscribeContext                = "unsubscribe api"
	ApplicationDeleteContext          = "delete application"
	APIDeleteContext                  = "delete API"
	APISearchContext                  = "search API"
	ApplicationSearchContext          = "search Application"
	KeyManagerSearchContext           = "search key manager"
	SubscriptionSearchContext         = "search Subscription"
	ApplicationKeySearchContext       = "search application keys"
	ErrMsgAPPIDEmpty                  = "application id is empty"
)

var (
	publisherAPIEndpoint              string
	storeApplicationEndpoint          string
	storeKeyManagerEndpoint           string
	storeSubscriptionEndpoint         string
	storeMultipleSubscriptionEndpoint string
	applicationDashBoardURLBase       string
	tokenManager                      token.Manager
	once                              sync.Once
)

// Init function initialize the API-M client. If there is an error, process exists with a panic.
func Init(manager token.Manager, conf APIM) {
	once.Do(func() {
		tokenManager = manager
		publisherAPIEndpoint = createEndpoint(conf.PublisherEndpoint, conf.PublisherAPIContext)
		storeApplicationEndpoint = createEndpoint(conf.StoreEndpoint, conf.StoreApplicationContext)
		storeKeyManagerEndpoint = createEndpoint(conf.StoreEndpoint, conf.StoreKeyManagerContext)
		storeSubscriptionEndpoint = createEndpoint(conf.StoreEndpoint, conf.StoreSubscriptionContext)
		storeMultipleSubscriptionEndpoint = createEndpoint(conf.StoreEndpoint, conf.StoreMultipleSubscriptionContext)
		applicationDashBoardURLBase = createEndpoint(conf.StoreEndpoint, "/devportal/applications/")
	})
}

// createEndpoint returns a endpoint from the given paths.
func createEndpoint(paths ...string) string {
	endpoint, err := utils.ConstructURL(paths...)
	if err != nil {
		log.HandleErrorAndExit("cannot construct endpoint", err)
	}
	return endpoint
}

// CreateAPI function creates an API with the provided API spec.
// Returns the API ID and any error encountered.
func CreateAPI(reqBody *APIReqBody) (*APICreateResp, error) {
	req, err := creatHTTPPOSTAPIRequest(publisherAPIEndpoint, reqBody)
	if err != nil {
		return nil, err
	}
	var resBody APICreateResp
	err = send(CreateAPIContext, req, &resBody, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	return &resBody, nil
}

// UpdateAPI updates an existing API under the given ID with the provided API spec.
// Returns the updated API and any error encountered.
func UpdateAPI(id string, reqBody *APIReqBody) (*APICreateResp, error) {
	endpoint, err := utils.ConstructURL(publisherAPIEndpoint, id)
	if err != nil {
		return nil, err
	}
	req, err := creatHTTPPUTAPIRequest(endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	var resBody APICreateResp
	err = send(UpdateAPIContext, req, &resBody, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resBody, nil
}

func ChangeLifeCycleStatus(apiID, action string) (*APIChangeLifeCycleResp, error) {
	endpoint, err := utils.ConstructURL(publisherAPIEndpoint, "change-lifecycle")
	if err != nil {
		return nil, err
	}
	req, err := creatHTTPPOSTAPIRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("apiId", apiID)
	q.Add("action", action)
	req.HTTPRequest().URL.RawQuery = q.Encode()
	var resBody APIChangeLifeCycleResp
	err = send(ChangeAPILifeCycleContext, req, &resBody, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resBody, nil
}

// GetAppDashboardURL returns DashBoard URL for the given Application.
func GetAppDashboardURL(appID string) string {
	return applicationDashBoardURLBase + "/" + appID + "/overview"
}

// CreateApplication creates an application with provided Application spec.
// Returns the Application ID and any error encountered.
func CreateApplication(reqBody *ApplicationCreateReq) (*ApplicationSearchInfo, error) {
	req, err := creatHTTPPOSTAPIRequest(storeApplicationEndpoint, reqBody)
	if err != nil {
		return nil, err
	}
	var resBody ApplicationSearchInfo
	err = send(CreateApplicationContext, req, &resBody, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	return &resBody, nil
}

// UpdateApplication updates an existing Application under the given ID with the provided Application spec.
// Returns any error encountered.
func UpdateApplication(id string, reqBody *ApplicationCreateReq) (*ApplicationSearchInfo, error) {
	endpoint, err := utils.ConstructURL(storeApplicationEndpoint, id)
	if err != nil {
		return nil, err
	}
	req, err := creatHTTPPUTAPIRequest(endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	var resBody ApplicationSearchInfo
	err = send(UpdateApplicationContext, req, &resBody, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resBody, nil
}

// GenerateKeys generates keys for the given application.
// Returns generated keys and any error encountered.
func GenerateKeys(appID string, reqBody *ApplicationKeyGenerateRequest) (*ApplicationKeyResp, error) {
	if appID == "" {
		return nil, errors.New(ErrMsgAPPIDEmpty)
	}
	generateApplicationKeyEndpoint, err := utils.ConstructURL(storeApplicationEndpoint, appID, "/generate-keys")
	if err != nil {
		return nil, errors.Wrap(err, "cannot construct endpoint")
	}
	req, err := creatHTTPPOSTAPIRequest(generateApplicationKeyEndpoint, reqBody)
	if err != nil {
		return nil, err
	}
	var resBody ApplicationKeyResp
	err = send(GenerateKeyContext, req, &resBody, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resBody, nil
}

func UpdateApplicationKeys(appID string, keyMappingID string, reqBody *ApplicationKeyGenerateRequest) (*ApplicationKeyResp, error) {
	if appID == "" {
		return nil, errors.New(ErrMsgAPPIDEmpty)
	}
	endpoint, err := utils.ConstructURL(storeApplicationEndpoint, appID, "/oauth-keys", keyMappingID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot construct endpoint")
	}
	req, err := creatHTTPPUTAPIRequest(endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	var resBody ApplicationKeyResp
	err = send(GenerateKeyContext, req, &resBody, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resBody, nil
}

func RegenerateKeys(appID string, keyMappingID string) (*ApplicationKeyResp, error) {
	if appID == "" {
		return nil, errors.New(ErrMsgAPPIDEmpty)
	}
	regenerateApplicationKeyEndpoint, err := utils.ConstructURL(storeApplicationEndpoint, appID, "/oauth-keys", keyMappingID, "/regenerate-secret")
	if err != nil {
		return nil, errors.Wrap(err, "cannot construct endpoint")
	}
	req, err := creatHTTPPOSTAPIRequest(regenerateApplicationKeyEndpoint, nil)
	if err != nil {
		return nil, err
	}
	var resBody ApplicationKeyResp
	err = send(RegenerateKeyContext, req, &resBody, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resBody, nil
}

func CleanupKeys(appID string, keyMappingID string) error {
	if appID == "" {
		return errors.New(ErrMsgAPPIDEmpty)
	}
	endpoint, err := utils.ConstructURL(storeApplicationEndpoint, appID, "/oauth-keys", keyMappingID, "/clean-up")
	if err != nil {
		return err
	}
	req, err := creatHTTPPOSTAPIRequest(endpoint, nil)
	if err != nil {
		return err
	}
	err = send(CleanupKeysContext, req, nil, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

func CreateSubscription(sub *SubscriptionReq) (*SubscriptionSearchInfo, error) {
	req, err := creatHTTPPOSTAPIRequest(storeSubscriptionEndpoint, sub)
	if err != nil {
		return nil, err
	}
	var resBody SubscriptionSearchInfo
	err = send(CreateSubscriptionContext, req, &resBody, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	return &resBody, nil
}

func UpdateSubscription(id string, reqBody *SubscriptionReq) (*SubscriptionSearchInfo, error) {
	endpoint, err := utils.ConstructURL(storeSubscriptionEndpoint, id)
	if err != nil {
		return nil, err
	}
	req, err := creatHTTPPUTAPIRequest(endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	var resBody SubscriptionSearchInfo
	err = send(UpdateSubscriptionContext, req, &resBody, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resBody, nil
}

// CreateMultipleSubscriptions creates the given subscriptions.
// Returns list of SubscriptionResp and any error encountered.
func CreateMultipleSubscriptions(subs []SubscriptionReq) ([]SubscriptionResp, error) {
	req, err := creatHTTPPOSTAPIRequest(storeMultipleSubscriptionEndpoint, subs)
	if err != nil {
		return nil, err
	}
	resBody := make([]SubscriptionResp, 0)
	err = send(CreateMultipleSubscriptionContext, req, &resBody, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}

// UnSubscribe method removes the given subscription.
// Returns any error encountered.
func UnSubscribe(subscriptionID string) error {
	endpoint, err := utils.ConstructURL(storeSubscriptionEndpoint, subscriptionID)
	if err != nil {
		return err
	}
	req, err := creatHTTPDELETEAPIRequest(endpoint)
	if err != nil {
		return err
	}
	err = send(UnSubscribeContext, req, nil, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

// DeleteApplication method deletes the given application.
// Returns any error encountered.
func DeleteApplication(applicationID string) error {
	endpoint, err := utils.ConstructURL(storeApplicationEndpoint, applicationID)
	if err != nil {
		return err
	}
	req, err := creatHTTPDELETEAPIRequest(endpoint)
	if err != nil {
		return err
	}
	err = send(ApplicationDeleteContext, req, nil, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAPI method deletes the given API.
// Returns any error encountered.
func DeleteAPI(apiID string) error {
	endpoint, err := utils.ConstructURL(publisherAPIEndpoint, apiID)
	if err != nil {
		return err
	}
	req, err := creatHTTPDELETEAPIRequest(endpoint)
	if err != nil {
		return err
	}
	err = client.Invoke(APIDeleteContext, req, nil, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

// send sends the given HTTP request, initialize the given response body if it is expected response code.
// Returns any error encountered.
func send(context string, req *client.HTTPRequest, resBody interface{}, expectedRespCode int) error {
	err := client.Invoke(context, req, resBody, expectedRespCode)
	if err != nil {
		return err
	}
	return nil
}

// getBodyReaderAndToken returns a token, a Reader for the given HTTP request body and any error encountered.
func getBodyReaderAndToken(reqBody interface{}) (string, io.ReadSeeker, error) {
	aT, err := tokenManager.Token()
	if err != nil {
		return "", nil, err
	}
	var bodyReader io.ReadSeeker
	if reqBody != nil {
		bodyReader, err = client.BodyReader(reqBody)
		if err != nil {
			return "", nil, err
		}
	}
	return aT, bodyReader, nil
}

func creatHTTPGETAPIRequest(endpoint string) (*client.HTTPRequest, error) {
	aT, err := tokenManager.Token()
	if err != nil {
		return nil, err
	}
	req, err := client.CreateHTTPGETRequest(aT, endpoint)
	if err != nil {
		return nil, err
	}
	return req, err
}

func creatHTTPPOSTAPIRequest(endpoint string, reqBody interface{}) (*client.HTTPRequest, error) {
	aT, bodyReader, err := getBodyReaderAndToken(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := client.CreateHTTPPOSTRequest(aT, endpoint, bodyReader)
	if err != nil {
		return nil, err
	}
	return req, err
}

func creatHTTPDELETEAPIRequest(endpoint string) (*client.HTTPRequest, error) {
	aT, err := tokenManager.Token()
	if err != nil {
		return nil, err
	}
	req, err := client.CreateHTTPDELETERequest(aT, endpoint)
	if err != nil {
		return nil, err
	}
	return req, err
}

// creatAPIMSearchHTTPRequest returns a API-M resource search request and any error encountered.
func creatAPIMSearchHTTPRequest(endpoint, query string) (*client.HTTPRequest, error) {
	aT, err := tokenManager.Token()
	if err != nil {
		return nil, err
	}
	req, err := client.CreateHTTPGETRequest(aT, endpoint)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("query", query)
	req.HTTPRequest().URL.RawQuery = q.Encode()
	return req, err
}

func creatHTTPPUTAPIRequest(endpoint string, reqBody interface{}) (*client.HTTPRequest, error) {
	aT, bodyReader, err := getBodyReaderAndToken(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := client.CreateHTTPPUTRequest(aT, endpoint, bodyReader)
	if err != nil {
		return nil, err
	}
	return req, err
}

// SearchAPIByNameVersion method returns API ID of the Given API.
// An error is returned if the number of result for the search is not equal to 1.
// Returns API ID and any error encountered.
func SearchAPIByNameVersion(apiName, version string) (string, error) {
	query := "name:" + apiName + " version:" + version
	req, err := creatAPIMSearchHTTPRequest(publisherAPIEndpoint, query)
	if err != nil {
		return "", err
	}
	var resp APISearchResp
	err = send(APISearchContext, req, &resp, http.StatusOK)
	if err != nil {
		return "", err
	}
	if resp.Count == 0 {
		return "", errors.New(fmt.Sprintf("couldn't find the API %s", apiName))
	}
	if resp.Count > 1 {
		return "", errors.New(fmt.Sprintf("returned more than one API for API %s", apiName))
	}
	return resp.List[0].ID, nil
}

func GetAPI(apiID string) (*APISearchInfo, error) {
	endpoint, err := utils.ConstructURL(publisherAPIEndpoint, apiID)
	if err != nil {
		return nil, err
	}
	req, err := creatHTTPGETAPIRequest(endpoint)
	if err != nil {
		return nil, err
	}
	var resp APISearchInfo
	err = send(APISearchContext, req, &resp, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchApplication method returns Application ID of the Given Application.
// An error is returned if the number of result for the search is not equal to 1.
// Returns Application ID and any error encountered.
func SearchApplication(appName string) (string, error) {
	req, err := creatAPIMSearchHTTPRequest(storeApplicationEndpoint, appName)
	if err != nil {
		return "", err
	}
	var resp ApplicationSearchResp
	err = send(ApplicationSearchContext, req, &resp, http.StatusOK)
	if err != nil {
		return "", err
	}
	if resp.Count == 0 {
		return "", errors.New(fmt.Sprintf("couldn't find the Application %s", appName))
	}
	if resp.Count > 1 {
		return "", errors.New(fmt.Sprintf("returned more than one Application for %s", appName))
	}
	return resp.List[0].ApplicationID, nil
}

func GetApplication(applicationID string) (*ApplicationSearchInfo, error) {
	endpoint, err := utils.ConstructURL(storeApplicationEndpoint, applicationID)
	if err != nil {
		return nil, err
	}
	req, err := creatHTTPGETAPIRequest(endpoint)
	if err != nil {
		return nil, err
	}
	var resp ApplicationSearchInfo
	err = send(ApplicationSearchContext, req, &resp, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func GetSubscription(subID string) (*SubscriptionSearchInfo, error) {
	endpoint, err := utils.ConstructURL(storeSubscriptionEndpoint, subID)
	if err != nil {
		return nil, err
	}
	req, err := creatHTTPGETAPIRequest(endpoint)
	if err != nil {
		return nil, err
	}
	var resp SubscriptionSearchInfo
	err = send(SubscriptionSearchContext, req, &resp, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func SearchKeyManager(keyManagerName string) (*KeyManagerSearchInfo, error) {
	req, err := creatHTTPGETAPIRequest(storeKeyManagerEndpoint)
	if err != nil {
		return nil, err
	}
	var resp KeyManagerSearchResp
	err = send(KeyManagerSearchContext, req, &resp, http.StatusOK)
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, errors.New(fmt.Sprintf("couldn't find the KeyManager %s", keyManagerName))
	}
	var keyManager KeyManagerSearchInfo
	for _, km := range resp.List {
		if km.Name == keyManagerName {
			keyManager = km
			break
		}
	}
	if keyManager.ID == "" {
		return nil, errors.New(fmt.Sprintf("couldn't find the KeyManager %s", keyManagerName))
	}
	return &keyManager, nil
}

func GetKeyManager(keyManagerID string) (*KeyManagerSearchInfo, error) {
	req, err := creatHTTPGETAPIRequest(storeKeyManagerEndpoint)
	if err != nil {
		return nil, err
	}
	var resp KeyManagerSearchResp
	err = send(KeyManagerSearchContext, req, &resp, http.StatusOK)
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, errors.New(fmt.Sprintf("couldn't find the KeyManager %s", keyManagerID))
	}
	var keyManager KeyManagerSearchInfo
	for _, km := range resp.List {
		if km.ID == keyManagerID {
			keyManager = km
			break
		}
	}
	if keyManager.ID == "" {
		return nil, errors.New(fmt.Sprintf("couldn't find the KeyManager %s", keyManagerID))
	}
	return &keyManager, nil
}

func GetApplicationKeys(applicationID string, keyMappingID string) (*ApplicationKeyResp, error) {
	endpoint, err := utils.ConstructURL(storeApplicationEndpoint, applicationID, "/oauth-keys", keyMappingID)
	if err != nil {
		return nil, err
	}
	req, err := creatHTTPGETAPIRequest(endpoint)
	if err != nil {
		return nil, err
	}
	var resp ApplicationKeyResp
	err = send(ApplicationKeySearchContext, req, &resp, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
