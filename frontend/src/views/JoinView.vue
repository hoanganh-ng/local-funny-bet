<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { leaderboardService } from '../services/leaderboard.service.js'
import BaseButton from '../components/base/BaseButton.vue'
import BaseSpinner from '../components/base/BaseSpinner.vue'

const route = useRoute()
const router = useRouter()

const isJoining = ref(true)
const error = ref(null)
const token = ref(route.query.token)

onMounted(async () => {
  if (!token.value) {
    error.value = 'No invite token provided'
    isJoining.value = false
    return
  }

  try {
    const response = await leaderboardService.join(token.value)
    router.push(`/leaderboard/${response.leaderboard_id}`)
  } catch (err) {
    error.value = err.message || 'Failed to join leaderboard'
    isJoining.value = false
  }
})

function goHome() {
  router.push('/')
}
</script>

<template>
  <div class="join-view">
    <div class="join-container">
      <div v-if="isJoining" class="joining-state">
        <BaseSpinner />
        <p class="status-text">Joining leaderboard...</p>
      </div>

      <div v-else-if="error" class="error-state">
        <div class="error-icon">⚠️</div>
        <h1 class="error-title">Unable to Join</h1>
        <p class="error-message">{{ error }}</p>
        <p class="error-hint">
          The invite link may have expired or is invalid.
          Ask for a new link from the leaderboard owner.
        </p>
        <BaseButton @click="goHome">
          Go to Home
        </BaseButton>
      </div>
    </div>
  </div>
</template>

<style scoped>
.join-view {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-4);
  background: var(--color-bg);
}

.join-container {
  max-width: 500px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.joining-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-6);
}

.status-text {
  font-size: var(--text-lg);
  color: var(--color-text-secondary);
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-8);
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
  border-radius: var(--radius-lg);
}

.error-icon {
  font-size: var(--text-6xl);
}

.error-title {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
  margin-top: var(--space-2);
}

.error-message {
  font-size: var(--text-base);
  color: var(--color-danger);
  font-weight: var(--font-medium);
}

.error-hint {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  line-height: var(--leading-normal);
  max-width: 400px;
}

[data-theme="light"] .error-state {
  border: var(--border-medium);
}
</style>
