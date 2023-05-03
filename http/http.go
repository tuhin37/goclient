package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ------------------------------- models -------------------------------

// this is a request object on which .Get() and .Post() methods are called
type HttpRequest struct {
	Url     string
	Method  string
	headers map[string]string
	body    []byte
}

// this is a response object which is returned by .Get() and .Post() methods
type HttpResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

// ------------------------------- constructor -------------------------------

// create a new http request object, use it when you are forming this from a gin context
func NewHttpRequestFromContext(method string, url string, body []byte, header http.Header) (*HttpRequest, error) {
	req := &HttpRequest{}
	req.Url = url
	req.Method = method

	// convert headers to map
	req.headers = make(map[string]string)
	for key, values := range header {
		headerValue := ""
		for _, v := range values {
			headerValue += headerValue + v
		}
		req.headers[key] = headerValue
	}

	req.body = body

	fmt.Println()
	return req, nil
}

// create a new http request object, use it when you are forming this manually
func NewHttpRequest(method string, url string, body any, header map[string]string) (*HttpRequest, error) {
	req := &HttpRequest{}
	req.Url = url
	req.Method = method
	req.headers = header
	// conver map to json []byte
	var err error
	req.body, err = json.Marshal(body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return req, nil
}

// ------------------------------- request methods -------------------------------

// set header by key in the request
func (req *HttpRequest) SetHeader(key string, value string) {
	req.headers[key] = value
}

// get header by key from the request
func (req *HttpRequest) GetHeader(key string) string {
	return req.headers[key]
}

// delete header by key from the request
func (req *HttpRequest) DeleteHeader(key string) {
	delete(req.headers, key)
}

// get the entire header map from the request
func (req *HttpRequest) GetHeaders() map[string]string {
	return req.headers
}

// delete all headers from the request
func (req *HttpRequest) DeleteHeaders() {
	req.headers = make(map[string]string)
}

// load a structure into the request body
func (req *HttpRequest) Encode(v any) error {
	var err error
	req.body, err = json.Marshal(v)
	if err != nil {
		return err
	}
	return nil
}

// ------------------------------- response methods -------------------------------

// get a map from response body
func (res *HttpResponse) GetBodyMap() (interface{}, error) {
	// Create a new HTTP request with the provided JSON payload and headers
	var bodyMap interface{}
	err := json.Unmarshal(res.Body, &bodyMap)
	if err != nil {
		return nil, err
	}
	return bodyMap, nil
}

// decode the response body into a struct
func (res *HttpResponse) Decode(v any) error {
	reader := bytes.NewReader(res.Body)
	err := json.NewDecoder(reader).Decode(&v)
	if err != nil {
		return err
	}
	return nil
}

func (res *HttpResponse) GetHeaders() map[string]string {
	return res.Headers
}

// --------------------------------- http methods ---------------------------------

// make get request
func (r *HttpRequest) Get(urls ...string) (int, HttpResponse, error) {
	url := r.Url
	if len(urls) == 1 {
		url = urls[0]
	}

	// Create a new HTTP request with the provided JSON payload and headers
	payloadJson, _ := json.Marshal(r.body)
	payload := bytes.NewReader(payloadJson)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}

	// Set the request headers
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}
	defer resp.Body.Close()

	// Convert the response headers into a map

	ResponseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		headerValue := ""
		for _, v := range values {
			headerValue += headerValue + v
		}
		ResponseHeaders[key] = headerValue
	}

	// Read the response body into a []byte variable
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}

	// return this object
	httpResponse := HttpResponse{}
	httpResponse.StatusCode = resp.StatusCode
	httpResponse.Body = body
	httpResponse.Headers = ResponseHeaders

	return resp.StatusCode, httpResponse, nil
}

// make post request
func (r *HttpRequest) Post(urls ...string) (int, HttpResponse, error) {
	url := r.Url
	if len(urls) == 1 {
		url = urls[0]
	}

	// Create a new HTTP request with the provided JSON payload and headers
	payloadJson, _ := json.Marshal(r.body)
	payload := bytes.NewReader(payloadJson)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}

	// Set the request headers
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}
	defer resp.Body.Close()

	// Convert the response headers into a map

	ResponseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		headerValue := ""
		for _, v := range values {
			headerValue += headerValue + v
		}
		ResponseHeaders[key] = headerValue
	}

	// Read the response body into a []byte variable
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}

	// return this object
	httpResponse := HttpResponse{}
	httpResponse.StatusCode = resp.StatusCode
	httpResponse.Body = body
	httpResponse.Headers = ResponseHeaders

	return resp.StatusCode, httpResponse, nil
}

// make put request
func (r *HttpRequest) Put(urls ...string) (int, HttpResponse, error) {
	url := r.Url
	if len(urls) == 1 {
		url = urls[0]
	}

	// Create a new HTTP request with the provided JSON payload and headers
	payloadJson, _ := json.Marshal(r.body)
	payload := bytes.NewReader(payloadJson)
	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}

	// Set the request headers
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}
	defer resp.Body.Close()

	// Convert the response headers into a map

	ResponseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		headerValue := ""
		for _, v := range values {
			headerValue += headerValue + v
		}
		ResponseHeaders[key] = headerValue
	}

	// Read the response body into a []byte variable
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}

	// return this object
	httpResponse := HttpResponse{}
	httpResponse.StatusCode = resp.StatusCode
	httpResponse.Body = body
	httpResponse.Headers = ResponseHeaders

	return resp.StatusCode, httpResponse, nil
}

// make patch request
func (r *HttpRequest) Patch(urls ...string) (int, HttpResponse, error) {
	url := r.Url
	if len(urls) == 1 {
		url = urls[0]
	}

	// Create a new HTTP request with the provided JSON payload and headers
	payloadJson, _ := json.Marshal(r.body)
	payload := bytes.NewReader(payloadJson)
	req, err := http.NewRequest("PATCH", url, payload)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}

	// Set the request headers
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}
	defer resp.Body.Close()

	// Convert the response headers into a map

	ResponseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		headerValue := ""
		for _, v := range values {
			headerValue += headerValue + v
		}
		ResponseHeaders[key] = headerValue
	}

	// Read the response body into a []byte variable
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}

	// return this object
	httpResponse := HttpResponse{}
	httpResponse.StatusCode = resp.StatusCode
	httpResponse.Body = body
	httpResponse.Headers = ResponseHeaders

	return resp.StatusCode, httpResponse, nil
}

// make delete request
func (r *HttpRequest) Delete(urls ...string) (int, HttpResponse, error) {
	url := r.Url
	if len(urls) == 1 {
		url = urls[0]
	}

	// Create a new HTTP request with the provided JSON payload and headers
	payloadJson, _ := json.Marshal(r.body)
	payload := bytes.NewReader(payloadJson)
	req, err := http.NewRequest("DELETE", url, payload)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}

	// Set the request headers
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}
	defer resp.Body.Close()

	// Convert the response headers into a map

	ResponseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		headerValue := ""
		for _, v := range values {
			headerValue += headerValue + v
		}
		ResponseHeaders[key] = headerValue
	}

	// Read the response body into a []byte variable
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, HttpResponse{}, err // the msg body becomes the error msg when response code is -1
	}

	// return this object
	httpResponse := HttpResponse{}
	httpResponse.StatusCode = resp.StatusCode
	httpResponse.Body = body
	httpResponse.Headers = ResponseHeaders

	return resp.StatusCode, httpResponse, nil
}

// ------------------------------ quick http functions -------------------------------
// this is a function, not a method. This can be called directly to perform a quick post request
func QuickHttpPOST(url string, postBody map[string]interface{}, bearerToken string) ([]byte, int, error) {
	// Set the data to be sent to the API
	payload, err := json.Marshal(postBody)
	if err != nil {
		return nil, 0, err
	}

	// Make the API call using the POST method
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

func QuickHttpGET(url string, postBody map[string]interface{}, bearerToken string) ([]byte, int, error) {
	// Set the data to be sent to the API
	payload, err := json.Marshal(postBody)
	if err != nil {
		return nil, 0, err
	}

	// Make the API call using the POST method
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

/* ------------------------------------ Examples -------------------------------------


POST
------------------------------------------
	payload := make(map[string]interface{})
	payload["phone_number"] = "254712345678"
	payloadJSON, _ := json.Marshal(payload)
	code, body := client.Http.Post("http://localhost:6000/get-otp", payloadJSON)
	fmt.Println("code: ", code)
	fmt.Println("body: ", body)


GET
-----------------------------------------------------------------
	// create a http request object
	httpRequest, _ := client.NewHttpRequestFromContext("GET", "https://x.com/get-list", payloadBytes, nil)
	httpRequest.SetHeader("Authorization", "Bearer AWUSRW53X7NMCODWI9YPAT0L4HJNR29A")

	// make the request
	code, resp, err := httpRequest.Get()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// if server replied json
	if resp.Headers["Content-Type"] == "application/json" {
		// load response body in a map
		var responseBody interface{}
		err = resp.Decode(&responseBody)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// all went well. return the response
		c.AsciiJSON(code, responseBody)
	}


DELETE
-----------------------------------------------------------------
	client.Http.BearerToken = "90f53ba5c74499477cd7880bbdc0159c1cc20e1ba51f4a96336cce52bb29b2121dcfe23d4b3b6f68805ab9c273b6c6b19292e5eb3c65dd36af98acdb5b2f9193"
	code, body := client.Http.Delete("http://localhost:5000/api/v2/ratelimit/delete/644d69bfd9d8cb1a99683aed")
	fmt.Println("code: ", code)
	fmt.Println("body: ", body)

*/
