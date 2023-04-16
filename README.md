# fuelsPricesFrance

fuelsPricesFrance allows you to observe the price of fuel on the French territory

## Installation

Install [Docker](https://www.docker.com/).

Install [npm](https://www.npmjs.com/).

Install [Vue.js](https://vuejs.org/) with npm.

Install [Go](https://go.dev/).

Start Docker.

Go in the folder fuelPricesFrance/misc/ :

```bash
cd fuelPricesFrance/misc/
```

Start the database :

```bash
docker-compose up -d
```

Go in the folder fuelPricesFrance/ :

```bash
cd fuelPricesFrance/
```

Run the server :

```bash
go run main/server.go
```

Go in the folder fuelPricesFrance/frontend/ :

```bash
cd fuelPricesFrance/frontend/
```

Run : 

```bash
npm run serve
```

You can now use the App :

```bash
App running at:
  - Local:   http://localhost:8080/
```