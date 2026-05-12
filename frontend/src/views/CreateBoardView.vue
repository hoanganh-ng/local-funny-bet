<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { leaderboardService } from '../services/leaderboard.service.js'
import BaseButton from '../components/base/BaseButton.vue'
import BaseInput from '../components/base/BaseInput.vue'

const router = useRouter()
const boardName = ref('')
const isCreating = ref(false)
const error = ref('')
const inviteLink = ref(null)
const isGeneratingLink = ref(false)

async function createBoard() {
  if (!boardName.value.trim()) {
    error.value = 'Board name is required'
    return
  }

  isCreating.value = true
  error.value = ''

  try {
    const board = await leaderboardService.create(boardName.value.trim())
    // After creating, generate invite link
    await generateInviteLink(board.id)
  } catch (err) {
    error.value = err.message || 'Failed to create board'
    isCreating.value = false
  }
}

async function generateInviteLink(boardId) {
  isGeneratingLink.value = true

  try {
    const link = await leaderboardService.generateInviteLink(boardId)
    inviteLink.value = link
  } catch (err) {
    console.error('Failed to generate invite link:', err)
  } finally {
    isGeneratingLink.value = false
    isCreating.value = false
  }
}

function copyLink() {
  if (!inviteLink.value) return
  navigator.clipboard.writeText(inviteLink.value.url)
}

function goToBoard() {
  if (!inviteLink.value) return
  router.push(`/leaderboards/${inviteLink.value.leaderboardId}`)
}
</script>

<template>
  <div class="create-board-view">
    <div class="container">
      <header class="page-header">
        <button class="back-btn" @click="$router.back()">
          ← Back to boards
        </button>
        <div class="header-content">
          <p class="header-label">New Leaderboard</p>
          <h1 class="page-title">Start a board.</h1>
          <p class="page-subtitle">
            You'll own it. Send the link, watch it fill up.
          </p>
        </div>
      </header>

      <div v-if="!inviteLink" class="create-form">
        <div class="form-section">
          <h3 class="form-section-title">Board Name</h3>
          <BaseInput
            v-model="boardName"
            placeholder="e.g. Engineering Cup"
            :error="error"
            @keyup.enter="createBoard"
          />
        </div>

        <div class="form-section">
          <h3 class="form-section-title">Settings</h3>
          <div class="settings-list">
            <div class="setting-row">
              <span class="setting-label">Visibility</span>
              <span class="setting-value">Invite link only</span>
            </div>
            <div class="setting-row">
              <span class="setting-label">Scoring</span>
              <span class="setting-value">1 pt per correct outcome</span>
            </div>
            <div class="setting-row">
              <span class="setting-label">Joiners need</span>
              <span class="setting-value">Company email</span>
            </div>
          </div>
        </div>

        <div class="info-box">
          <span class="info-icon">ℹ️</span>
          <div>
            <strong>Invite links expire in 10 minutes</strong>
            <p>Signed-out links, regenerate a fresh one whenever you need.</p>
          </div>
        </div>

        <BaseButton
          size="lg"
          :loading="isCreating"
          :disabled="!boardName.trim()"
          @click="createBoard"
        >
          Create & generate invite link
        </BaseButton>
      </div>

      <div v-else class="invite-result">
        <div class="success-message">
          <span class="success-icon">✓</span>
          <h2 class="success-title">Board created!</h2>
          <p class="success-subtitle">
            Share this link with your team to get them in.
          </p>
        </div>

        <div class="invite-link-box">
          <input
            readonly
            :value="inviteLink.url"
            class="invite-link-input"
            @click="(e) => e.target.select()"
          />
          <BaseButton @click="copyLink">
            Copy link
          </BaseButton>
        </div>

        <div class="invite-actions">
          <BaseButton variant="secondary" @click="goToBoard">
            Go to board
          </BaseButton>
          <BaseButton variant="ghost" @click="$router.push('/leaderboards')">
            Create another
          </BaseButton>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.create-board-view {
  min-height: 100vh;
  padding: var(--space-8) var(--space-6);
}

.container {
  max-width: 600px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: var(--space-10);
}

.back-btn {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  padding: var(--space-2) 0;
  cursor: pointer;
  margin-bottom: var(--space-6);
  transition: color var(--duration-fast);
}

.back-btn:hover {
  color: var(--color-text-primary);
}

.header-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.header-label {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-accent);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.page-title {
  font-family: var(--font-display);
  font-size: var(--text-5xl);
  font-weight: var(--font-black);
  color: var(--color-text-primary);
  line-height: var(--leading-tight);
}

.page-subtitle {
  font-size: var(--text-lg);
  color: var(--color-text-secondary);
  line-height: var(--leading-normal);
}

/* Create Form */
.create-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.form-section-title {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.settings-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-5);
  background: var(--color-surface);
  border: var(--border-hairline);
}

[data-theme="dark"] .settings-list {
  border-radius: var(--radius-md);
}

.setting-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-4);
}

.setting-label {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.setting-value {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-primary);
  text-align: right;
}

.info-box {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-5);
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  line-height: var(--leading-normal);
}

[data-theme="dark"] .info-box {
  border-radius: var(--radius-md);
  background: rgba(94,106,210,0.05);
  border-color: rgba(94,106,210,0.2);
}

.info-icon {
  font-size: var(--text-xl);
}

.info-box strong {
  display: block;
  color: var(--color-text-primary);
  margin-bottom: var(--space-1);
}

/* Invite Result */
.invite-result {
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
}

.success-message {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
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
  background: rgba(63,185,80,0.2);
  box-shadow: 0 0 40px rgba(63,185,80,0.3);
}

.success-title {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.success-subtitle {
  font-size: var(--text-base);
  color: var(--color-text-secondary);
}

.invite-link-box {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-5);
  background: var(--color-surface);
  border: var(--border-hairline);
}

[data-theme="dark"] .invite-link-box {
  border-radius: var(--radius-lg);
}

.invite-link-input {
  flex: 1;
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: transparent;
  border: none;
  outline: none;
  padding: var(--space-2);
}

.invite-actions {
  display: flex;
  gap: var(--space-3);
  justify-content: center;
}

@media (max-width: 768px) {
  .create-board-view {
    padding: var(--space-6) var(--space-4);
  }

  .page-title {
    font-size: var(--text-4xl);
  }

  .invite-link-box {
    flex-direction: column;
  }

  .invite-actions {
    flex-direction: column;
  }
}
</style>
