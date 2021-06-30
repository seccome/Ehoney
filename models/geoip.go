package models

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/geoip"
	_ "github.com/aliyun/alibaba-cloud-sdk-go/services/geoip"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const apiKey string = ""
const ipStackURL string = "http://api.ipstack.com/"
const torExitURL string = "https://check.torproject.org/torbulkexitlist"

var isExitNodeIP map[string]bool

// GeoData for the json
type GeoData struct {
	ContinentCode string `json:"continent_code"`
	CountryCode   string `json:"country_code"`
	City          string `json:"city"`
	ContinentName string `json:"continent_name"`
	CountryName	  string `json:"country_name"`
	RegionName	  string `json:"region_name"`
	TorExitNode   bool
}

// SetupExitNodeMap setups of the exit node map
func SetupExitNodeMap() error {
	// Make map
	isExitNodeIP = make(map[string]bool)

	f, err := ioutil.ReadFile("conf/torbulkexitlist")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	temp := strings.Split(string(f), "\n")

	// Finally, turn into map
	for _, val := range temp {
		if val != "" {
			trimmedVal := strings.TrimSpace(val)
			isExitNodeIP[trimmedVal] = true
			// fmt.Println(trimmedVal)
		}
	}
	return nil
}

// IsTorExitNode returns true if it's an exit node
func isTorExitNode(ip string) bool {
	_, exists := isExitNodeIP[ip]
	if exists {
		return true
	}
	return false
}

// GetGeoData returns geographical data for an IP address
func GetGeoData(ip string) GeoData {
	// Create default GeoData
	geoData := GeoData{
		ContinentCode: "unk",
		CountryCode:   "unk",
		City:          "unk",
		ContinentName: "unk",
		CountryName:   "unk",
		RegionName:    "unk",
		TorExitNode:   isTorExitNode(ip),
	}

	// Create HTTP request with the IP address
	req, err := http.NewRequest("GET", ipStackURL+ip, nil)
	if err != nil {
		log.Println(fmt.Sprintf("Error creating API Request: %v", err))
		return geoData
	}

	// Add the API token
	q := req.URL.Query()
	q.Add("access_key", apiKey)
	req.URL.RawQuery = q.Encode()

	// Accept json as response
	req.Header.Add("Accept", "application/json")

	// Create HTTP client and submit request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(fmt.Sprintf("Error submitting API Request: %v", err))
		return geoData
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(fmt.Sprintf("Error reading response data: %v", err))
		return geoData
	}

	err = json.Unmarshal(respBody, &geoData)
	if err != nil {
		log.Println(fmt.Sprintf("Unable to get IP data for %s: %v", ip, err))
	}
	return geoData
}

func GetGeoDataForAliYun(ip string) *geoip.DescribeIpv4LocationResponse {
	urlKeyId := ""
	urlKeySec := ""
	resp, err := http.Get(urlKeyId)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	accessKeyId := gjson.Get(string(body),"decryptContent")
	resp1, err := http.Get(urlKeySec)
	defer resp1.Body.Close()
	body1, err := ioutil.ReadAll(resp1.Body)
	accessKeySec := gjson.Get(string(body1),"decryptContent")

	client, err := geoip.NewClientWithAccessKey("cn-hangzhou", accessKeyId.Str, accessKeySec.Str)
	request := geoip.CreateDescribeIpv4LocationRequest()
	request.Scheme = "https"
	request.Ip = ip
	response, err := client.DescribeIpv4Location(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
	return response
}