# Parking Lot Simulator 2019
<hr>
This is a basic 'Parking Lot' SPA utilizing Golang &amp; Vue

## Packages used in this software:
<strong>GORM</strong>  : github.com/jinzhu/gorm (ORM for Golang)

<strong>PQ</strong>    : github.com/lib/pq (Postgresql Library for Golang

<strong>VIPER</strong> : github.com/spf13/viper (Configuration package for Golang)

<strong>ECHO</strong>  : gopkg.in/echo.v3 (HTTP Routing package for Golang)

Glide is used to manage package dependencies, and package depencies can be found in the glide.yml file.

## Docker Commands:

<em>Please ensure you have an up-to-date and working version of Docker</em>
```
1. docker run -d -p 5439:5432 --name postgres postgres
This will create the postgres container that we will use for this application

2. docker run -d -p 5439:5432 --name postgres postgres
This will let you enter the postgres DB and poke around while the application is running
```

