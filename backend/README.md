# Fuels_france

Fuels_france allows you to observe the price of fuel on the French territory

## Installation

Install Docker.

Install docker-compose.

Install npm.

Install Vue.js with npm.

Install Go.

Start Docker.

Go in the folder Fuels_france/misc/ :

```bash
cd Fuels_france/misc/
```

Start the database :

```bash
docker-compose up -d
```

Go in the folder Fuels_france/ :

```bash
cd Fuels_france/
```

Run the server :

```bash
go run main/server.go
```

Go in the folder fuelPricesFrance/frontend/ :

```bash
cd Fuels_france/frontend/
```

Install :

```bash
npm install
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
