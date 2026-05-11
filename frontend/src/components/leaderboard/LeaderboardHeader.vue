<script setup>
import { ref } from 'vue'
import BaseButton from '../base/BaseButton.vue'
import BaseModal from '../base/BaseModal.vue'
import BaseInput from '../base/BaseInput.vue'
import { leaderboardService } from '../../services/leaderboard.service.js'
import { useAuthStore } from '../../store/auth.store.js'

const props = defineProps({
  leaderboard: {
    type: Object,
    required: true
  }
})

const authStore = useAuthStore()
const showInviteModal = ref(false)
const inviteLink = ref('')
const isGenerating = ref(false)
const copiedToClipboard = ref(false)

const isOwner = () => {
  return authStore.user?.id === props.leaderboard.ownerId
}

const openInviteModal = async () => {
  showInviteModal.value = true
  isGenerating.value = true
  copiedToClipboard.value = false

  try {
    const response = await leaderboardService.generateInvite(props.leaderboard.id)
    const baseUrl = window.location.origin
    inviteLink.value = `${baseUrl}/join?token=${response.invite_token}`
  } catch (err) {
    console.error('Failed to generate invite:', err)
  } finally {
    isGenerating.value = false
  }
}

const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(inviteLink.value)
    copiedToClipboard.value = true
    setTimeout(() => {
      copiedToClipboard.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

const closeModal = () => {
  showInviteModal.value = false
  inviteLink.value = ''
}
</script>

<template>
  <div class="leaderboard-header">
    <div class="header-info">
      <h2 class="leaderboard-name">{{ leaderboard.name }}</h2>
      <p class="member-count">{{ leaderboard.memberCount }} members</p>
    </div>
    <BaseButton
      v-if="isOwner()"
      @click="openInviteModal"
    >
      Invite
    </BaseButton>

    <BaseModal
      :open="showInviteModal"
      title="Invite to Leaderboard"
      @close="closeModal"
    >
      <div class="invite-content">
        <p class="invite-description">
          Share this link with others to invite them to join this leaderboard.
        </p>
        <div v-if="isGenerating" class="generating">
          Generating invite link...
        </div>
        <div v-else class="invite-link-container">
          <BaseInput
            :model-value="inviteLink"
            readonly
            @click="copyToClipboard"
          />
          <BaseButton @click="copyToClipboard">
            {{ copiedToClipboard ? 'Copied!' : 'Copy' }}
          </BaseButton>
        </div>
      </div>
    </BaseModal>
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

.invite-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.invite-description {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  line-height: var(--leading-normal);
}

.generating {
  padding: var(--space-4);
  text-align: center;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

.invite-link-container {
  display: flex;
  gap: var(--space-2);
}

.invite-link-container :deep(input) {
  cursor: pointer;
}
</style>
