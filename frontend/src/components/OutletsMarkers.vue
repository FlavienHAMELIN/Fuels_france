<template>
    <l-marker ref="marker" :key="marker.id" :lat-lng="marker.coordinates">
      <l-popup :content="marker.content" :options="{autoClose:false, closeOnClick: false}"></l-popup>
      <l-icon ref="icon">
        <img class="outlet-icon" :src="marker.imageUrl"/> 
      </l-icon>
    </l-marker>
  </template>
  
  <script>
  import { LIcon, LMarker, LPopup } from 'vue2-leaflet'
  export default {
    components: { LIcon, LMarker, LPopup },
    props: {
      marker: {
        type: Object,
        required: true
      }
    },
  }
  </script>
  
  <style>
    .outlet-icon {
      height: 15px;
      width: auto;
    }
  </style>