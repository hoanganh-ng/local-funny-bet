<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth.store.js'
import { authService } from '../services/auth.service.js'
import { useTournament } from '../composables/useTournament.js'
import BaseButton from '../components/base/BaseButton.vue'

const router = useRouter()
const authStore = useAuthStore()
const { tournament } = useTournament()

// Redirect if already authenticated
onMounted(() => {
  if (authStore.isAuthenticated) {
    router.push('/')
  }
})

function handleLogin() {
  authService.login()
}
</script>

<template>
  <div class="login-view">
    <!-- Brand panel - desktop only -->
    <div class="brand-panel">
      <div class="brand-content">
        <p class="tournament-label">{{ tournament?.name || 'Continental Cup 2026' }}</p>
        <h1 class="headline">
          Pick the<br>
          <span class="accent-word">winners.</span><br>
          Beat your<br>
          team.
        </h1>
        <p class="tagline">
          Predict every match of the tournament. One point per
          correct call. Join your team's leaderboards. The office
          bragging rights are everything.
        </p>
        <div class="features">
          <span class="feature-item">1 pt / correct outcome</span>
          <span class="feature-dot">·</span>
          <span class="feature-item">Locks at kickoff</span>
          <span class="feature-dot">·</span>
          <span class="feature-item">Real-time leaderboards</span>
        </div>
      </div>
    </div>

    <!-- Auth panel -->
    <div class="auth-panel">
      <div class="auth-card">
        <div class="auth-content">
          <h2 class="auth-title">Sign in to predict</h2>
          <p class="auth-subtitle">
            We only let people with a company email in.
          </p>

          <BaseButton size="lg" @click="handleLogin">
            <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M17.64 9.205c0-.639-.057-1.252-.164-1.841H9v3.481h4.844a4.14 4.14 0 0 1-1.796 2.716v2.259h2.908c1.702-1.567 2.684-3.875 2.684-6.615Z" fill="#4285F4"/>
              <path d="M9 18c2.43 0 4.467-.806 5.956-2.18l-2.908-2.259c-.806.54-1.837.86-3.048.86-2.344 0-4.328-1.584-5.036-3.711H.957v2.332A8.997 8.997 0 0 0 9 18Z" fill="#34A853"/>
              <path d="M3.964 10.71A5.41 5.41 0 0 1 3.682 9c0-.593.102-1.17.282-1.71V4.958H.957A8.996 8.996 0 0 0 0 9c0 1.452.348 2.827.957 4.042l3.007-2.332Z" fill="#FBBC05"/>
              <path d="M9 3.58c1.321 0 2.508.454 3.44 1.345l2.582-2.58C13.463.891 11.426 0 9 0A8.997 8.997 0 0 0 .957 4.958L3.964 7.29C4.672 5.163 6.656 3.58 9 3.58Z" fill="#EA4335"/>
            </svg>
            Continue with Google
          </BaseButton>

          <p class="auth-note">
            We never see your password. Sign-in is handled by Google
            with our company SSO.
          </p>

          <p class="invite-note">
            Get an invite link? Open it in a new tab once signed in.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-view {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 1fr 1fr;
  background: var(--color-bg);
}

/* ──── Brand Panel ────  */
.brand-panel {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-12) var(--space-8);
  position: relative;
}

.brand-content {
  max-width: 540px;
}

.tournament-label {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
  margin-bottom: var(--space-6);
}

.headline {
  font-family: var(--font-display);
  font-size: clamp(48px, 6vw, 72px);
  font-weight: var(--font-black);
  color: var(--color-text-primary);
  line-height: var(--leading-tight);
  letter-spacing: var(--tracking-tighter);
  margin-bottom: var(--space-8);
}

.accent-word {
  color: var(--color-accent);
  position: relative;
}

[data-theme="dark"] .accent-word {
  text-shadow: 0 0 40px var(--color-accent-glow);
}

.tagline {
  font-size: var(--text-lg);
  color: var(--color-text-secondary);
  line-height: var(--leading-relaxed);
  margin-bottom: var(--space-8);
}

.features {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: var(--space-3);
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.feature-dot {
  color: var(--color-border);
}

/* ──── Auth Panel ──── */
.auth-panel {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-12) var(--space-8);
}

.auth-card {
  width: 100%;
  max-width: 480px;
  padding: var(--space-12) var(--space-10);
  background: var(--color-surface);
  border: var(--border-hairline);
  box-shadow: var(--shadow-lg);
}

[data-theme="dark"] .auth-card {
  border-radius: var(--radius-lg);
  backdrop-filter: blur(20px);
}

[data-theme="light"] .auth-card {
  border-radius: 0;
  border-top: var(--border-accent);
  box-shadow: none;
}

.auth-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.auth-title {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
  letter-spacing: var(--tracking-tight);
}

.auth-subtitle {
  font-size: var(--text-base);
  color: var(--color-text-secondary);
  line-height: var(--leading-normal);
}

.auth-note,
.invite-note {
  font-size: var(--text-sm);
  color: var(--color-text-disabled);
  line-height: var(--leading-normal);
}

/* Mobile: stack vertically, hide brand panel decorations */
@media (max-width: 768px) {
  .login-view {
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr;
    padding: var(--space-6) var(--space-4);
  }

  .brand-panel {
    padding: var(--space-8) var(--space-4) var(--space-6);
  }

  .headline {
    font-size: clamp(32px, 10vw, 48px);
    margin-bottom: var(--space-6);
  }

  .tagline {
    font-size: var(--text-base);
    margin-bottom: var(--space-6);
  }

  .auth-panel {
    padding: var(--space-6) var(--space-4);
  }

  .auth-card {
    padding: var(--space-8) var(--space-6);
  }

  .auth-title {
    font-size: var(--text-2xl);
  }
}
</style>
