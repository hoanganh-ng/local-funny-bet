<script setup>
import { authService } from '../services/auth.service.js'
import BaseButton from '../components/base/BaseButton.vue'

function handleLogin() {
  authService.login()
}
</script>

<template>
  <div class="login-view">
    <div class="login-container">
      <h1 class="login-title">World Cup 2026 Predictor</h1>
      <p class="login-subtitle">
        Predict match outcomes, compete with friends, and climb the leaderboard
      </p>
      <BaseButton size="lg" @click="handleLogin">
        Sign in with Google
      </BaseButton>
    </div>
  </div>
</template>

<style scoped>
.login-view {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-4);
  background: var(--color-bg);
}

.login-container {
  max-width: 600px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-8);
  text-align: center;
}

.login-title {
  font-family: var(--font-display);
  font-size: var(--text-5xl);
  font-weight: var(--font-black);
  color: var(--color-text-primary);
  line-height: var(--leading-tight);
  letter-spacing: var(--tracking-tight);
}

.login-subtitle {
  font-size: var(--text-lg);
  color: var(--color-text-secondary);
  line-height: var(--leading-normal);
  max-width: 480px;
}

@media (max-width: 768px) {
  .login-title {
    font-size: var(--text-4xl);
  }

  .login-subtitle {
    font-size: var(--text-base);
  }
}
</style>
