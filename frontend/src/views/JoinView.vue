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
const joinedBoard = ref(null) // { id, name, alreadyMember }

onMounted(async () => {
  if (!token.value) {
    error.value = 'No invite token provided'
    isJoining.value = false
    return
  }

  try {
    const response = await leaderboardService.join(token.value)
    const board = await leaderboardService.getOne(response.leaderboard_id)
    joinedBoard.value = { id: response.leaderboard_id, name: board.name, alreadyMember: false }
    isJoining.value = false
    setTimeout(() => router.push(`/leaderboards/${response.leaderboard_id}`), 2500)
  } catch (err) {
    if (err.status === 409 && err.data?.leaderboard_id) {
      const board = await leaderboardService.getOne(err.data.leaderboard_id).catch(() => null)
      joinedBoard.value = {
        id: err.data.leaderboard_id,
        name: board?.name ?? 'this leaderboard',
        alreadyMember: true,
      }
      isJoining.value = false
      return
    }
    const msg = err.message || ''
    if (msg.includes('expired')) {
      error.value = 'This invite link has expired. Ask the board owner to generate a new one.'
    } else if (msg.includes('invalid') || msg.includes('required')) {
      error.value = 'This invite link is invalid or malformed.'
    } else if (msg.includes('already')) {
      error.value = "You're already a member of this leaderboard."
    } else {
      error.value = 'Failed to join leaderboard. Please try again.'
    }
    isJoining.value = false
  }
})

function goToBoard() {
  router.push(`/leaderboards/${joinedBoard.value.id}`)
}

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

      <div v-else-if="joinedBoard" class="success-state">
        <div class="success-icon">✓</div>
        <h1 class="success-title">
          {{ joinedBoard.alreadyMember ? 'Already a member' : "You're in!" }}
        </h1>
        <p class="success-board-name">{{ joinedBoard.name }}</p>
        <p v-if="!joinedBoard.alreadyMember" class="redirect-hint">
          Redirecting to board...
        </p>
        <BaseButton @click="goToBoard">Go to board</BaseButton>
      </div>

      <div v-else-if="error" class="error-state">
        <div class="error-icon">⚠️</div>
        <h1 class="error-title">Unable to Join</h1>
        <p class="error-message">{{ error }}</p>
        <p v-if="error.includes('expired') || error.includes('invalid')" class="error-hint">
          Ask the board owner to generate a fresh invite link.
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

/* Success state */
.success-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-8);
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
  border-radius: var(--radius-lg);
}

.success-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-4xl);
  background: var(--color-correct);
  color: var(--color-accent-text);
  border-radius: var(--radius-pill);
}

[data-theme="dark"] .success-icon {
  background: rgba(63, 185, 80, 0.2);
  box-shadow: 0 0 40px rgba(63, 185, 80, 0.3);
}

.success-title {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.success-board-name {
  font-size: var(--text-lg);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
}

.redirect-hint {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
}

/* Error state */
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

[data-theme="light"] .error-state,
[data-theme="light"] .success-state {
  border: var(--border-medium);
}
</style>
