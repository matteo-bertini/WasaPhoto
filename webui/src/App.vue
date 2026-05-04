<script setup>
/**
 * @file App.vue
 * @description Root component of the application. It acts as the main container 
 * where the Vue Router injects the matched components for each route.
 */
import { RouterView } from 'vue-router'
</script>

<script>
export default {
  name: 'App'
}
</script>

<template>
  <!-- 
    RouterView is responsible for rendering the component 
    corresponding to the current URL route. 
  -->
  <RouterView />
</template>

<style>
/**
 * GLOBAL STYLES & RESET
 * These styles apply to the entire application and ensure visual consistency.
 */

/* 
 * 1. Box Sizing Reset: Ensures padding and borders are included in the element's total width/height.
 * 2. Margin/Padding Reset: Removes default browser spacing.
 */
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  overflow-x: hidden;
  
  /* Global typography and base theme setup */
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  background-color: #0a0a0a; /* Fallback for the dark theme */
  color: #ffffff;
  overflow-y: scroll; /* Forza la presenza della corsia della scrollbar verticale */
  overflow-x: hidden; /* Mantiene il fix per lo scroll orizzontale */
  
  /* Smooth scrolling for anchor links or programmatic scrolls */
  scroll-behavior: smooth;
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
}
::-webkit-scrollbar-track {
  background: #0a0a0a;
}
::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}
::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.4);
}
</style>