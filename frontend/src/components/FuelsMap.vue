  <template>
    <div id="fuels">
      <label for="fuelName-select">Choisissez un carburant :</label>
      <select id="fuelName-select" name="fuelsNames" @change="$eventName => onChangeName($eventName)" class="form-control">
        <option value="">Carburant</option>
        <option value="Gazole">Gazole</option>
        <option value="SP98">SP98</option>
        <option value="SP95">SP95</option>
        <option value="E10">E10</option>
        <option value="GPLc">GPLc</option>
        <option value="E85">E85</option>
      </select>
      <label for="region-select">Choisissez une région :</label>
      <select id="region-select" name="region" @change="$eventRegion => onChangeRegion($eventRegion)" class="form-control">
        <option value="">Localisation</option>
        <option value="Toute la France">Toute la France</option>
        <option value="Auvergne-Rhône-Alpes">Auvergne-Rhône-Alpes</option>
        <option value="Bourgogne-Franche-Comté">Bourgogne-Franche-Comté</option>
        <option value="Bretagne">Bretagne</option>
        <option value="Centre-Val de Loire">Centre-Val de Loire</option>
        <option value="Corse">Corse</option>
        <option value="Grand Est">Grand Est</option>
        <option value="Hauts-de-France">Hauts-de-France</option>
        <option value="Île-de-France">Île-de-France</option>
        <option value="Normandie">Normandie</option>
        <option value="Nouvelle-Aquitaine">Nouvelle-Aquitaine</option>
        <option value="Occitanie">Occitanie</option>
        <option value="Pays de la Loire">Pays de la Loire</option>
        <option value="Provence-Alpes-Côte d'Azur">Provence-Alpes-Côte d'Azur</option>
      </select>
      <label for="departement-select">Ou choisissez un numéro de département : </label>
      <input id="departement-select" type="text" v-model="input" placeholder="28" @change="$eventDepartement => onChangeDepartement($eventDepartement)" class="input" />
      <div id="average-price" class="average-price-display">
        {{ this.averagePrice }}
      </div>
      <l-map :center="center" :zoom="zoom" class="map" ref="map" @update:center="centerUpdated">
        <l-tile-layer :url="url"></l-tile-layer>
        <outlet v-for="marker in this.markers" :key="marker.id" :marker="marker"></outlet>
      </l-map>
    </div>
  </template>
  
  <script>
  import { LMap, LTileLayer } from 'vue2-leaflet';
  import outlet from './OutletsMarkers';
  import 'leaflet/dist/leaflet.css';
  import axios from 'axios';
  import { mean, round, quantileSeq } from 'mathjs';

  export default {
    components: {
      LMap,
      LTileLayer,
      outlet,
    },
    region: "",
    fuelName: "",
    data () {
      return {
        url: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
        zoom: 8,
        markers: [],
      };
    },
    methods: {
      async updateData(fuel, region, departement) {
        let response = null;
        if (region !== undefined && region != "Toute la France") {
          response = await axios.get(`http://localhost:8181/fuelsFrance?name=${fuel}&region=${region}`);
        } else if (region !== undefined && region == "Toute la France") {
          response = await axios.get(`http://localhost:8181/fuelsFrance?name=${fuel}`);
        }
        if (departement !== undefined) {
          response = await axios.get(`http://localhost:8181/fuelsFrance?name=${fuel}&dep_code=${departement}`);
        }
        this.localisationUpdated(response);
        this.markersUpdated(response);
        
      },
      onChangeRegion(eventRegion) {
        this.region = eventRegion.target.value;
        if (this.fuelName !== undefined && this.region !== undefined && this.fuelName != "" && this.region != "") {
          this.updateData(this.fuelName, this.region, undefined);
        }
      },
      onChangeName (eventName) {
        this.fuelName = eventName.target.value;
        if (this.fuelName !== undefined && this.fuelName != "") {
          this.updateData(this.fuelName, this.region, this.departement);
        }
      },
      onChangeDepartement (eventDepartement) {
        this.departement = eventDepartement.target.value;
        if (this.fuelName !== undefined && this.departement !== undefined && this.fuelName != "" && this.departement != "") {
          this.updateData(this.fuelName, undefined, this.departement)
        }
      },
      centerUpdated (center) {
        this.center = center;
      },
      markersUpdated (response) {
        let arrMarkers = [];
        let arrPrices = [];
        const data = response.data;

        for (let i = 0; i < data.length; i++) {
          arrPrices.push(data[i]["price"]);
        }
        this.averagePrice = "Prix moyen : " + String(round(mean(arrPrices), 3)) + " €";
        arrPrices.sort();
        const quartiles = quantileSeq(arrPrices, [0.25, 0.50, 0.75]);

        for (let i = 0; i < data.length; i++) {
          const price = data[i]["price"];
          if (price < quartiles[0]) {
            arrMarkers.push(
              {
                  id: i, 
                  imageUrl: 'station-vert.png', 
                  coordinates: [data[i]["latitude"], data[i]["longitude"]],
                  content: `${data[i]["name"]} : ${price} € | Adresse : ${data[i]["address"]}, ${data[i]["city"]} |`,
              }
            )
          }
          if (price >= quartiles[0] && price < quartiles[1]) {
            arrMarkers.push(
              {
                  id: i, 
                  imageUrl: 'station-jaune.png', 
                  coordinates: [data[i]["latitude"], data[i]["longitude"]],
                  content: `${data[i]["name"]} : ${price} € | Adresse : ${data[i]["address"]}, ${data[i]["city"]} |`,
              }
            )
          }
          if (price >= quartiles[1] && price < quartiles[2]) {
            arrMarkers.push(
              {
                  id: i, 
                  imageUrl: 'station-orange.png', 
                  coordinates: [data[i]["latitude"], data[i]["longitude"]],
                  content: `${data[i]["name"]} : ${price} € | Adresse : ${data[i]["address"]}, ${data[i]["city"]} |`,
              }
            )
          }
          if (price >= quartiles[2]) {
            arrMarkers.push(
              {
                  id: i, 
                  imageUrl: 'station-rouge.png', 
                  coordinates: [data[i]["latitude"], data[i]["longitude"]],
                  content: `${data[i]["name"]} : ${price} € | Adresse : ${data[i]["address"]}, ${data[i]["city"]} |`,
              }
            )
          }
        }
        this.markers = arrMarkers;
      },
      localisationUpdated (response) {
        let sommeLat = 0;
        let sommeLong = 0;
        const len = response.data.length;
        for (let i = 0; i < len; i++) {
          sommeLat = sommeLat + response.data[i]["latitude"];
          sommeLong = sommeLong + response.data[i]["longitude"];
        }
        let latitude = sommeLat / len;
        let longitude = sommeLong / len;
        if (this.region != "Toute la France") {
          this.centerUpdated([latitude, longitude]);
        } else if (this.region == "Toute la France") {
          this.centerUpdated([46.227638, 2.213749]);
        }
      },
    },
  }
  </script>
  
  <style>
    .map {
      position: absolute;
      width: 100%;
      height: 100%;
      overflow : hidden;
    }

    .form-control {
      padding: 0.5em;
      margin: 0.5em;
    }

    .input {
      padding: 0.5em;
      margin: 0.5em;
    }

    .average-price-display {
      padding: 0.5em;
      margin: 0.5em;
    }
  </style>