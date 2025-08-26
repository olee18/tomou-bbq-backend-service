package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"io/ioutil"
	"laotop_final/logs"
	"mime/multipart"
	"net/http"
)

type HttpClientTrail interface {
	CallApi(queryUrl string, request interface{}) ([]byte, error)
	CallApiWithFile(queryUrl, field string, ctx *fiber.Ctx, request interface{}) ([]byte, error)
	CallApiIpro(queryUrl string, request interface{}) ([]byte, error)
}

type httpClientTrail struct {
	apiClient http.Client
}

func (h httpClientTrail) CallApi(queryUrl string, request interface{}) ([]byte, error) {
	httpClientApi := h.apiClient
	marshal, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	response, err := http.NewRequest("POST", queryUrl, bytes.NewBuffer(marshal))
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	response.Header.Add("Content-Type", "application/json")
	res, err := httpClientApi.Do(response)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	defer res.Body.Close()
	readAll, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	if res.StatusCode == 200 {
		return readAll, nil
	}
	catchErr := map[string]interface{}{}
	err = json.Unmarshal([]byte(readAll), &catchErr)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	logs.Error(catchErr["error"].(string))
	return nil, errors.New(catchErr["error"].(string))

}

func (h httpClientTrail) CallApiWithFile(queryUrl, field string, ctx *fiber.Ctx, request interface{}) ([]byte, error) {
	file, err := ctx.FormFile(field)
	if err != nil {
		return nil, errors.New(field + " required")
	}

	// Open the file for reading
	fileContent, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	// Create a buffer to store the multipart request body
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	// Create a multipart file field with the same field name
	part, err := writer.CreateFormFile(field, file.Filename)
	if err != nil {
		return nil, err
	}

	// Copy the file content into the part
	_, err = io.Copy(part, fileContent)
	if err != nil {
		return nil, err
	}

	// Add other form fields
	for key, val := range request.(map[string]interface{}) {
		_ = writer.WriteField(key, fmt.Sprintf("%v", val))
	}

	// Close the multipart writer
	writer.Close()
	//fmt.Println(requestBody)

	//Create the HTTP request with the multipart body
	req, err := http.NewRequest("POST", queryUrl, requestBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		return responseBody, nil
	}
	catchErr := map[string]interface{}{}
	err = json.Unmarshal([]byte(responseBody), &catchErr)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return nil, errors.New(catchErr["error"].(string))
}

func (h httpClientTrail) CallApiIpro(queryUrl string, request interface{}) ([]byte, error) {
	httpClientApi := h.apiClient
	marshal, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	response, err := http.NewRequest("POST", queryUrl, bytes.NewBuffer(marshal))
	if err != nil {
		return nil, err
	}
	response.Header.Add("Content-Type", "application/json")
	res, err := httpClientApi.Do(response)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	defer res.Body.Close()
	readAll, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	if res.StatusCode == 200 {
		decodeError := map[string]interface{}{}
		err = json.Unmarshal([]byte(readAll), &decodeError)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		responseStatusCode := decodeError["status"].(string)
		responseDescription := decodeError["description"].(string)
		if responseStatusCode != "1" {
			logs.Error(errors.New(responseDescription))
			return nil, errors.New(responseDescription)
		}

		return readAll, nil
	}
	catchErr := map[string]interface{}{}
	err = json.Unmarshal([]byte(readAll), &catchErr)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	logs.Error(catchErr["error"].(string))
	return nil, errors.New(catchErr["error"].(string))
}

func NewHttpClientTrail(apiClient http.Client) HttpClientTrail {
	return &httpClientTrail{apiClient: apiClient}
}
