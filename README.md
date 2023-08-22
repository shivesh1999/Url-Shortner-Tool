# URL Shortener

URL Shortener is a simple web application that shortens long URLs into more manageable and shareable short URLs. It's built using GoLang and uses Redis as the database for storing URL mappings.

## Table of Contents
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Features

- Shorten long URLs into short and manageable links.
- Custom short URL preference available.
- Rate limiting to prevent abuse.
- Expiry time for short URLs.
- Simple API endpoints for URL shortening and resolving.

## Getting Started

### Prerequisites

- Docker (for running the application and Redis)

### Installation

1. Clone this repository to your local machine:

   ```bash
   git clone https://github.com/your-username/URL_Shortner.git
   cd URL_Shortner
   ```

2. Create a `.env` file in the api repo of the project with the following content:

   ```env
   DB_ADDR="redis-db:6379"
   DB_PASS=""
   APP_PORT=":3000"
   DOMAIN="localhost:3000"
   API_QUOTA=10
   ```

   Adjust the values as needed.

3. Build and run the application using Docker Compose:

   ```bash
   docker-compose up
   ```

The application should now be up and running. You can access the API at http://localhost:3000.

## Usage

- To shorten a URL, make a `POST` request to `/shorten` with a JSON body containing the URL to be shortened.
- To access a shortened URL, visit `http://localhost:3000/shortURL` where `shortURL` is the short code generated.

## Configuration

- `DB_ADDR`: Address of the Redis database.
- `DB_PASS`: Password for the Redis database.
- `APP_PORT`: Port on which the GoLang application will run.
- `DOMAIN`: Domain used for generating short URLs.
- `API_QUOTA`: Rate limit for API requests.

## Contributing

Contributions are welcome! Please feel free to fork this repository and submit pull requests to contribute to the project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```

Remember to replace placeholders like `your-username` with your actual GitHub username. You can customize the README further to include more detailed information about your project, such as additional features, use cases, and any other relevant details.