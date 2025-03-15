package helper

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Earth radius in kilometers
const EarthRadius = 6371

type GeoLocationResponse struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Name      string  `json:"name"`
}

func GetGeoLocationByIP(ip string) (GeoLocationResponse, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return GeoLocationResponse{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
		City string  `json:"city"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GeoLocationResponse{}, err
	}

	return GeoLocationResponse{Latitude: result.Lat, Longitude: result.Lon, Name: result.City}, nil
}

func GetGeoLocationByCity(city string) (GeoLocationResponse, error) {
	externalUrl := "https://nominatim.openstreetmap.org/search"
	query := url.Values{}
	query.Set("q", city)
	query.Set("format", "json")
	query.Set("limit", "1") // Get first result only

	fullUrl := fmt.Sprintf("%s?%s", externalUrl, query.Encode())

	resp, err := http.Get(fullUrl)
	if err != nil {
		return GeoLocationResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GeoLocationResponse{}, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	var results []struct {
		Lat  string `json:"lat"`
		Lon  string `json:"lon"`
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return GeoLocationResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(results) == 0 {
		return GeoLocationResponse{}, fmt.Errorf("no location found for city: %s", city)
	}

	result := results[0]
	lat, _ := strconv.ParseFloat(result.Lat, 64)
	lon, _ := strconv.ParseFloat(result.Lon, 64)
	return GeoLocationResponse{Latitude: float64(lat), Longitude: lon, Name: result.Name}, nil
}

func CalculateDistance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadius * c
}

func GetRealIP(ctx *gin.Context) string {
	// Check X-Forwarded-For header
	forwardedFor := ctx.Request.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// Get first in multiple IP addresses
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	realIP := ctx.Request.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	return ctx.ClientIP()
}
