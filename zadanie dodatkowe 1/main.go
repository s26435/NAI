package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WeatherData struktura na potrzeby odczytania odpowiedzi JSON z WeatherStack
type WeatherData struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		Temperature int `json:"temperature"`
		WeatherDesc []struct {
			Value string `json:"value"`
		} `json:"weather_descriptions"`
	} `json:"current"`
}

func main() {
	apiKey := "a06e49cf288b3b7145dd817577faae59"
	url := fmt.Sprintf("http://api.weatherstack.com/current?access_key=%s&query=Gdansk", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Błąd podczas wykonywania zapytania:", err)
		return
	}
	defer resp.Body.Close()

	var weatherData WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		fmt.Println("Błąd podczas dekodowania odpowiedzi JSON:", err)
		return
	}

	fmt.Printf("Pogoda w %s, %s:\n", weatherData.Location.Name, weatherData.Location.Country)
	fmt.Printf("Temperatura: %d°C\n", weatherData.Current.Temperature)
	if len(weatherData.Current.WeatherDesc) > 0 {
		fmt.Printf("Opis pogody: %s\n", weatherData.Current.WeatherDesc[0].Value)
	}
}
