package service_implementation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"operation-borderless/internal/domain/dto"
	"operation-borderless/pkg/config"
)

type ExternalAPIClient struct {
	Http   http.Client
	Config *config.Config
}

func NewExternalAPIClient(conf *config.Config) *ExternalAPIClient {
	return &ExternalAPIClient{Config: conf}
}

func (client *ExternalAPIClient) GetPairExchangeRate(ctx context.Context, baseCurrency, targetCurrency string) (dto.ExchangeRate, error) {

	var rate dto.ExchangeRate

	apiUrl := fmt.Sprintf("%s%s/pair/%s/%s", client.Config.ForexAPIUrl, client.Config.ForexAPIKey, baseCurrency, targetCurrency)

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Println("error calling request: ", err)
		return rate, err
	}

	// Add Basic headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := client.Http.Do(request.WithContext(ctx))
	if err != nil {
		log.Println("error getting response: ", err)
		return rate, err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated &&
		response.StatusCode != http.StatusAccepted {
		resp, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("error reading response.body to byte: ", err)
			return rate, err
		}

		return rate, errors.New(string(resp))
	}

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("error reading response.body to byte: ", err)
		return rate, err
	}

	err = json.Unmarshal(resp, &rate)
	if err != nil {
		log.Println("error unmarshalling response: ", err)
		return rate, err
	}

	return rate, nil
}

func (client *ExternalAPIClient) GetUSDExchangeRate(ctx context.Context) (dto.ExchangeRate, error) {

	var rate dto.ExchangeRate

	apiUrl := fmt.Sprintf("%s%s/latest/USD", client.Config.ForexAPIUrl, client.Config.ForexAPIKey)

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Println("error calling request: ", err)
		return rate, err
	}

	// Add Basic headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := client.Http.Do(request.WithContext(ctx))
	if err != nil {
		log.Println("error getting response: ", err)
		return rate, err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated &&
		response.StatusCode != http.StatusAccepted {
		resp, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("error reading response.body to byte: ", err)
			return rate, err
		}

		return rate, errors.New(string(resp))
	}

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("error reading response.body to byte: ", err)
		return rate, err
	}

	err = json.Unmarshal(resp, &rate)
	if err != nil {
		log.Println("error unmarshalling response: ", err)
		return rate, err
	}

	return rate, nil
}

func (client *ExternalAPIClient) GetUserCountry(ctx context.Context, ipAddress string) (string, error) {

	var res dto.IpAPIResponse

	apiUrl := fmt.Sprintf("http://ip-api.com/json/%s", ipAddress)
	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Println("error calling request: ", err)
		return "", err
	}

	// Add Basic headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := client.Http.Do(request.WithContext(ctx))
	if err != nil {
		log.Println("error getting response: ", err)
		return "", err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated &&
		response.StatusCode != http.StatusAccepted {
		resp, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("error reading response.body to byte: ", err)
			return "", err
		}

		return "", errors.New(string(resp))
	}

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("error reading response.body to byte: ", err)
		return "", err
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		log.Println("error unmarshalling response: ", err)
		return "", err
	}

	return res.Country, nil
}
