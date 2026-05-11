<template>
  <div id="app">
    <ThemeToggle />
    <router-view />
  </div>
</template>

<script setup>
import ThemeToggle from './components/base/ThemeToggle.vue'
</script>

<style>
@import './assets/css/global.css';
</style>
