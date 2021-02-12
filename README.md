# Scaffolding Go

[![YourActionName Actions Status](https://github.com/chiahsoon/go_scaffold/workflows/go_scaffold/badge.svg)](https://github.com/chiahsoon/go_scaffold/actions)

Robust backend API scaffold written in Go with the following features:
1. ORM with [GORM](https://github.com/go-gorm/gorm).
2. JWT Authentication with [jwt-go](https://github.com/dgrijalva/jwt-go).
3. Logging using [zap](https://github.com/uber-go/zap).
4. Configuration Management [viper](https://github.com/spf13/viper).

## Points to Note
> Please read this section first before proceeding to the next sections.

### Configuration
1. The config file - environment name mapping is as follows: `./configs/<ENV>.yaml`.
2. You can define environment variables, and they'll automatically override the ones in your config files. 
   However, nested environment variables are not yet supported.

## Setup Local `dev` Environment

### Manual
1. Export the appropriate environment variables.
   ``` 
   export ENV=dev 
   ```
2. Start the task-runner for live-reload during development.
   ``` 
   realize start 
   ```
### Docker
1. Export your appropriate environment variables.
   ``` 
   export ENV=docker-dev 
   ```
2. Edit the `docker-compose.yml` if needed.
2. Build the images using the `docker-compose.yml` file.
    ``` 
   docker-compose build 
   ```
3. Run the containers.
    ``` 
    docker-compose up   
    ```