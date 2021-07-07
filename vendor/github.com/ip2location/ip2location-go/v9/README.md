[![Go Report Card](https://goreportcard.com/badge/github.com/ip2location/ip2location-go)](https://goreportcard.com/report/github.com/ip2location/ip2location-go)


IP2Location Go Package
======================

This Go package provides a fast lookup of country, region, city, latitude, longitude, ZIP code, time zone, ISP, domain name, connection type, IDD code, area code, weather station code, station name, mcc, mnc, mobile brand, elevation, and usage type from IP address by using IP2Location database. This package uses a file based database available at IP2Location.com. This database simply contains IP blocks as keys, and other information such as country, region, city, latitude, longitude, ZIP code, time zone, ISP, domain name, connection type, IDD code, area code, weather station code, station name, mcc, mnc, mobile brand, elevation, usage type, address type and IAB category as values. It supports both IP address in IPv4 and IPv6.

This package can be used in many types of projects such as:

 - select the geographically closest mirror
 - analyze your web server logs to determine the countries of your visitors
 - credit card fraud detection
 - software export controls
 - display native language and currency 
 - prevent password sharing and abuse of service 
 - geotargeting in advertisement

The database will be updated in monthly basis for the greater accuracy. Free LITE databases are available at https://lite.ip2location.com/ upon registration.

The paid databases are available at https://www.ip2location.com under Premium subscription package.


Installation
=======

```
go get github.com/ip2location/ip2location-go/v9
```

Example
=======

```go
package main

import (
	"fmt"
	"github.com/ip2location/ip2location-go"
)

func main() {
	db, err := ip2location.OpenDB("./IP-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ZIPCODE-TIMEZONE-ISP-DOMAIN-NETSPEED-AREACODE-WEATHER-MOBILE-ELEVATION-USAGETYPE-ADDRESSTYPE-CATEGORY.BIN")
	
	if err != nil {
		fmt.Print(err)
		return
	}
	ip := "8.8.8.8"
	results, err := db.Get_all(ip)
	
	if err != nil {
		fmt.Print(err)
		return
	}
	
	fmt.Printf("country_short: %s\n", results.Country_short)
	fmt.Printf("country_long: %s\n", results.Country_long)
	fmt.Printf("region: %s\n", results.Region)
	fmt.Printf("city: %s\n", results.City)
	fmt.Printf("isp: %s\n", results.Isp)
	fmt.Printf("latitude: %f\n", results.Latitude)
	fmt.Printf("longitude: %f\n", results.Longitude)
	fmt.Printf("domain: %s\n", results.Domain)
	fmt.Printf("zipcode: %s\n", results.Zipcode)
	fmt.Printf("timezone: %s\n", results.Timezone)
	fmt.Printf("netspeed: %s\n", results.Netspeed)
	fmt.Printf("iddcode: %s\n", results.Iddcode)
	fmt.Printf("areacode: %s\n", results.Areacode)
	fmt.Printf("weatherstationcode: %s\n", results.Weatherstationcode)
	fmt.Printf("weatherstationname: %s\n", results.Weatherstationname)
	fmt.Printf("mcc: %s\n", results.Mcc)
	fmt.Printf("mnc: %s\n", results.Mnc)
	fmt.Printf("mobilebrand: %s\n", results.Mobilebrand)
	fmt.Printf("elevation: %f\n", results.Elevation)
	fmt.Printf("usagetype: %s\n", results.Usagetype)
	fmt.Printf("addresstype: %s\n", results.Addresstype)
	fmt.Printf("category: %s\n", results.Category)
	fmt.Printf("api version: %s\n", ip2location.Api_version())
	
	db.Close()
}
```

Dependencies
============

The complete database is available at https://www.ip2location.com under subscription package.


IPv4 BIN vs IPv6 BIN
====================

Use the IPv4 BIN file if you just need to query IPv4 addresses.
Use the IPv6 BIN file if you need to query BOTH IPv4 and IPv6 addresses.


Copyright
=========

Copyright (C) 2021 by IP2Location.com, support@ip2location.com
