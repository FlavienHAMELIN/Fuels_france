  <template>
    <div id="fuels">
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
  import { quantileSeq } from 'mathjs';

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
        center: [48.866667, 2.333333],
        zoom: 8,
        markers: [],
      };
    },
    methods: {
      async updateData(fuel, region) {
        let response = null;
        if (region != "Toute la France") {
          response = await axios.get(`http://localhost:8181/fuelsFrance?name=${fuel}&region=${region}`);
          this.localisationUpdated(response);
          this.markersUpdated(response);
        } else if (region == "Toute la France"){
          response = await axios.get(`http://localhost:8181/fuelsFrance?name=${fuel}`);
          this.localisationUpdated(response);
          this.markersUpdated(response);
        }
        
      },
      onChangeRegion(eventRegion) {
        this.region = eventRegion.target.value;
        if (this.fuelName !== undefined && this.region !== undefined && this.fuelName != "" && this.region != "") {
          this.updateData(this.fuelName, this.region);
        }
      },
      onChangeName (eventName) {
        this.fuelName = eventName.target.value;
        if (this.fuelName !== undefined && this.region !== undefined && this.fuelName != "" && this.region != "") {
          this.updateData(this.fuelName, this.region);
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
                  content: `${data[i]["name"]} : ${price} €`,
              }
            )
          }
          if (price >= quartiles[0] && price < quartiles[1]) {
            arrMarkers.push(
              {
                  id: i, 
                  imageUrl: 'station-jaune.png', 
                  coordinates: [data[i]["latitude"], data[i]["longitude"]],
                  content: `${data[i]["name"]} : ${price} €`,
              }
            )
          }
          if (price >= quartiles[1] && price < quartiles[2]) {
            arrMarkers.push(
              {
                  id: i, 
                  imageUrl: 'station-orange.png', 
                  coordinates: [data[i]["latitude"], data[i]["longitude"]],
                  content: `${data[i]["name"]} : ${price} €`,
              }
            )
          }
          if (price >= quartiles[2]) {
            arrMarkers.push(
              {
                  id: i, 
                  imageUrl: 'station-rouge.png', 
                  coordinates: [data[i]["latitude"], data[i]["longitude"]],
                  content: `${data[i]["name"]} : ${price} €`,
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
  </style>