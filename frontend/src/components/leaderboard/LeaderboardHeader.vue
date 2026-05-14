<script setup>
import { ref } from 'vue'
import BaseButton from '../base/BaseButton.vue'

const props = defineProps({
  leaderboard: {
    type: Object,
    required: true
  },
  isOwner: {
    type: Boolean,
    default: false
  },
  onGetInviteLink: {
    type: Function,
    default: null
  }
})

// 'idle' | 'loading' | 'success' | 'error'
const buttonState = ref('idle')

async function handleInviteClick() {
  if (buttonState.value === 'loading' || !props.onGetInviteLink) return
  buttonState.value = 'loading'
  try {
    await props.onGetInviteLink()
    buttonState.value = 'success'
    setTimeout(() => { buttonState.value = 'idle' }, 2000)
  } catch {
    buttonState.value = 'error'
  }
}
</script>

<template>
  <div class="leaderboard-header">
    <div class="header-info">
      <h2 class="leaderboard-name">{{ leaderboard.name }}</h2>
      <p class="member-count">{{ leaderboard.memberCount }} members</p>
    </div>
    <BaseButton
      v-if="isOwner"
      variant="secondary"
      size="sm"
      :loading="buttonState === 'loading'"
      @click="handleInviteClick"
    >
      <span v-if="buttonState === 'success'">Copied!</span>
      <span v-else-if="buttonState === 'error'">Failed, try again</span>
      <span v-else>Get Invite Link</span>
    </BaseButton>
  </div>
</template>

<style scoped>
.leaderboard-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--space-4);
  padding: var(--space-6);
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
  border-radius: var(--radius-lg);
}

.header-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.leaderboard-name {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.member-count {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}
</style>
