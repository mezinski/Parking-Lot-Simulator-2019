# Go-with-Vue-2
# A Parking Lot Application

This is a basic Parking Lot SPA utilizing Golang &amp; Vue

## Packages used in this software:
GORM  :            Golang ORM          : github.com/jinzhu/gorm

PQ    :     Golang Postgres Driver     : github.com/lib/pq

VIPER :   Golang Configuration Package : github.com/spf13/viper

ECHO  :       Golang HTTP routing      : gopkg.in/echo.v3

Glide is used to manage package dependencies, and package depencies can be found in the glide.yml file.

## Docker Commands:

<em>Please ensure you have an up-to-date and working version of Docker</em>
```
1. docker run -d -p 5439:5432 --name postgres postgres
This will create the postgres container that we will use for this application

2. docker run -d -p 5439:5432 --name postgres postgres
This will let you enter the postgres DB and poke around while the application is running
```
