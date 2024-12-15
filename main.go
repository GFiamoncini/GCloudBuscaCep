package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

type ViaCep struct {
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type Temperature struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func main() {
	http.HandleFunc("/", WeatherHandler)
	port := ":8080"
	log.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s\n", r.Method, r.URL.String())

	cep := r.URL.Query().Get("cep")
	if cep == "" || !isValidCep(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	location, err := BuscaCep(cep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	tempC, err := GetTemperature(location.Localidade, location.Uf)
	if err != nil {
		http.Error(w, "failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	response := Temperature{
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273.15,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func isValidCep(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

func BuscaCep(cep string) (*ViaCep, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var c ViaCep
	err = json.Unmarshal(body, &c)
	if err != nil || c.Localidade == "" {
		return nil, errors.New("invalid CEP")
	}

	return &c, nil
}

func GetTemperature(city, uf string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		apiKey = "5ef06d7fbd5743b69ed150449241512"
	} else {
		return 0, fmt.Errorf("WEATHER_API_KEY is not set")
	}

	location := url.QueryEscape(fmt.Sprintf("%s,%s", city, uf))
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, location)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("WeatherAPI response: %s\n", string(body))
		return 0, fmt.Errorf("failed to fetch weather data: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return 0, err
	}

	return weather.Current.TempC, nil
}
