# URL Shortener Documentation

## Overview

The URL Shortener service is a Golang-based application that allows you to shorten long URLs.
It supports two data storage engines: local machine memory and Redis, which can be configured in the configuration file.
You can run this service using Docker Compose, and it is also deployed on Render.

## Table of Contents

1. [Installation](#installation)
    - [Docker Compose](#docker-compose)

2. [Endpoints](#endpoints)
    - [Shorten URL Endpoint](#shorten-url-endpoint)
    - [Expand URL Endpoint](#expand-url-endpoint)
    - [Metrics URL Endpoint](#metrics-url-endpoint)

3. [Swagger Documentation](#swagger-documentation)

4. [Configuration](#configuration)
    - [Config File](#config-file)
    - [Environment Variables](#environment-variables)

## Installation

### Docker Compose

To run the URL Shortener service using Docker Compose, follow these steps:

1. Clone the repository containing the service code.
    ```bash
    git clone git@github.com:harshabangi/url-shortener.git
    ```
2. Configure the Service with the data storage engine you want to use: either "memory" or "redis."

   - To use local machine memory for storage, set the `data_storage_engine` field to "memory" in the configuration file.

     Example configuration for "memory":
       ```yaml
       data_storage_engine: "memory"
       ```

   - To use Redis as the data storage engine, set the `data_storage_engine` field to "redis" and provide the `redis_url`
     field with the appropriate Redis server URL.

     Example configuration for "redis" (replace with your Redis server URL):
     ```yaml
     data_storage_engine: "redis"
     redis_url: "redis://redis:6379"
     ```

3. If you have made changes to the configuration file and need to update the Docker images, use the following command to
   build the images from the modified source code and configuration:
    ```
    docker-compose build --no-cache
    ```
4. Once the images are built or if you haven't made changes to the configuration file, start the URL Shortener service
   with the following command:
   ```
   docker-compose up
   ```

### Endpoints

#### Shorten URL Endpoint

To shorten a URL using the URL Shortener service, you can make a POST request to the `/v1/shorten` endpoint.

**Render Web Service Example:**

```bash
curl --request POST --url "https://url-shortener-xk2d.onrender.com/v1/shorten" --data '{"url": "https://harsha.com"}' --header 'Content-Type: application/json'
```

*Example Response:*

```
{"url":"https://url-shortener-xk2d.onrender.com/g2GJ99W"}
```

**Localhost Example:**

Assuming your service is running locally on port 8080:

```bash
curl --request POST --url "http://localhost:8080/v1/shorten" --data '{"url": "https://www.google.com"}' --header 'Content-Type: application/json'
```

*Example Response:*

```
{"url":"http://localhost:8080/g2GJ99W"}
```

#### Expand URL Endpoint

Expand URL request, redirects you (302) to the original URL.

**Render Web Service Example:**

```bash
curl --request GET --url "https://url-shortener-xk2d.onrender.com/g2GJ99W" --header 'Content-Type: application/json'
```

**Localhost Example:**

Assuming your service is running locally on port 8080:

```bash
curl --request GET --url "http://localhost:8080/g2GJ99W" --header 'Content-Type: application/json'
```

#### Metrics URL Endpoint

To retrieve metrics from the URL Shortener service, make a GET request to the /v1/metrics endpoint.
You can also include a 'limit' query parameter to specify the maximum number of metrics to retrieve.

**Render Web Service Example:**

```bash
curl --request GET --url "https://url-shortener-xk2d.onrender.com/v1/metrics" --header 'Content-Type: application/json'
```

**Localhost Example:**

Assuming your service is running locally on port 8080:

```bash
curl --request GET --url "http://localhost:8080/v1/metrics" --header 'Content-Type: application/json'
```

*Example Response:*

```
[
  {
    "domain_name": "www.google.com",
    "frequency": 2
  },
  {
    "domain_name": "www.youtube.com",
    "frequency": 1
  },
  {
    "domain_name": "www.facebook.com",
    "frequency": 1
  }
]
```

### Swagger-documentation

You can access the Swagger documentation for the URL Shortener service both locally and via the Render web service.

#### Localhost:

If you are running the service locally, you can access the Swagger documentation
at http://localhost:8080/swagger/index.html#/.

#### Render Web Service:

The Swagger documentation for the Render web service is available
at https://url-shortener-xk2d.onrender.com/swagger/index.html#/.

### Configuration

#### Config File

The URL Shortener service should be configured using a `config.yaml` file. Below are the configuration parameters you
can set in this file:

- **`short_url_domain`** (string): The domain where the shortened URLs will be accessible. For local development, it's
  set to `"http://localhost:8080"`. For the Render deployment, it's set to `"https://url-shortener-xk2d.onrender.com"`.
  This is the base URL for the generated shortened URLs.

- **`listen_addr`** (string): The address and port on which the service listens for incoming requests. In this example,
  it's set to `":8080"`, which means it listens on all available network interfaces on port 8080.

- **`data_storage_engine`** (string): Specifies the data storage engine used by the service. In this configuration, it's
  set to `"memory"` for local machine memory storage.

- **`short_url_length`** (integer): Defines the desired length of the generated short URLs. In this example, it's set
  to `7`.

- **`redis_url`** (string): The URL for the Redis server used for data storage when the `data_storage_engine` is set
  to `"redis"`. In this case, it's set to `"redis://redis:6379"`.

Example `config.yaml`:

```yaml
# config.yaml
short_url_domain: "http://localhost:8080" # For local development
# For Render deployment
# short_url_domain: "https://url-shortener-xk2d.onrender.com"
listen_addr: ":8080"
data_storage_engine: "memory"
short_url_length: 7
redis_url: "redis://redis:6379"
```

#### Environment Variables

The server configuration file should be supplied via environment variables.

```
CONFIG_FILE=/app/config/config.yaml
```
