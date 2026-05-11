<template>
  <div id="app">
    <div class="ambient-blob ambient-blob--1"></div>
    <div class="ambient-blob ambient-blob--2"></div>
    <div class="ambient-blob ambient-blob--3"></div>
    <div class="noise-overlay"></div>
    <ThemeToggle />
    <router-view />
  </div>
</template>

<script setup>
import ThemeToggle from './components/base/ThemeToggle.vue'
</script>

<style>
@import './assets/css/global.css';

/* Ambient blobs - dark theme only */
[data-theme="dark"] .ambient-blob {
  position: fixed;
  border-radius: 50%;
  filter: blur(120px);
  opacity: 0.15;
  pointer-events: none;
  z-index: -1;
  animation: blob-float var(--blob-duration) ease-in-out infinite;
}

[data-theme="dark"] .ambient-blob--1 {
  top: 10%;
  left: 15%;
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(94,106,210,0.4), transparent 70%);
  animation-delay: 0s;
}

[data-theme="dark"] .ambient-blob--2 {
  bottom: 15%;
  right: 20%;
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, rgba(63,185,80,0.3), transparent 70%);
  animation-delay: 3s;
}

[data-theme="dark"] .ambient-blob--3 {
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 700px;
  height: 700px;
  background: radial-gradient(circle, rgba(94,106,210,0.2), transparent 70%);
  animation-delay: 6s;
}

/* Hide blobs in light theme */
[data-theme="light"] .ambient-blob {
  display: none;
}

/* Noise texture - dark theme only */
[data-theme="dark"] .noise-overlay {
  position: fixed;
  inset: 0;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='2' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
  opacity: 0.03;
  pointer-events: none;
  z-index: 10;
  mix-blend-mode: overlay;
}

[data-theme="light"] .noise-overlay {
  display: none;
}
</style>
