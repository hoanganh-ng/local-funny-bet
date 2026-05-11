<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth.store.js'
import { useLeaderboardStore } from '../store/leaderboard.store.js'
import { leaderboardService } from '../services/leaderboard.service.js'
import { matchService } from '../services/match.service.js'
import BaseButton from '../components/base/BaseButton.vue'
import BaseModal from '../components/base/BaseModal.vue'
import BaseInput from '../components/base/BaseInput.vue'
import BaseCard from '../components/base/BaseCard.vue'
import MatchList from '../components/match/MatchList.vue'
import AdSlot from '../components/base/AdSlot.vue'

const router = useRouter()
const authStore = useAuthStore()
const leaderboardStore = useLeaderboardStore()

const showCreateModal = ref(false)
const newLeaderboardName = ref('')
const isCreating = ref(false)
const createError = ref('')

const upcomingMatches = ref([])
const isLoadingMatches = ref(false)

onMounted(async () => {
  await leaderboardStore.fetchAll()

  isLoadingMatches.value = true
  try {
    upcomingMatches.value = await matchService.list('scheduled')
  } catch (err) {
    console.error('Failed to load matches:', err)
  } finally {
    isLoadingMatches.value = false
  }
})

function openCreateModal() {
  showCreateModal.value = true
  newLeaderboardName.value = ''
  createError.value = ''
}

function closeCreateModal() {
  showCreateModal.value = false
  newLeaderboardName.value = ''
  createError.value = ''
}

async function createLeaderboard() {
  if (!newLeaderboardName.value.trim()) {
    createError.value = 'Name is required'
    return
  }

  isCreating.value = true
  createError.value = ''

  try {
    await leaderboardService.create(newLeaderboardName.value.trim())
    await leaderboardStore.fetchAll()
    closeCreateModal()
  } catch (err) {
    createError.value = err.message || 'Failed to create leaderboard'
  } finally {
    isCreating.value = false
  }
}

function goToLeaderboard(id) {
  router.push(`/leaderboard/${id}`)
}
</script>

<template>
  <div class="home-view">
    <div class="container">
      <AdSlot position="top" />

      <section class="welcome-section">
        <h1 class="welcome-title">
          Welcome, {{ authStore.user?.name || 'Player' }}
        </h1>
        <BaseButton @click="openCreateModal">
          Create Leaderboard
        </BaseButton>
      </section>

      <section class="leaderboards-section">
        <h2 class="section-title">Your Leaderboards</h2>
        <div v-if="leaderboardStore.leaderboards.length === 0" class="empty-state">
          <p>You haven't joined any leaderboards yet.</p>
          <p class="empty-hint">Create one or ask a friend for an invite link.</p>
        </div>
        <div v-else class="leaderboards-grid">
          <BaseCard
            v-for="leaderboard in leaderboardStore.leaderboards"
            :key="leaderboard.id"
            hoverable
            class="leaderboard-card"
            @click="goToLeaderboard(leaderboard.id)"
          >
            <h3 class="leaderboard-name">{{ leaderboard.name }}</h3>
            <p class="leaderboard-meta">{{ leaderboard.memberCount }} members</p>
          </BaseCard>
        </div>
      </section>

      <section class="matches-section">
        <h2 class="section-title">Upcoming Matches</h2>
        <div v-if="isLoadingMatches" class="loading-state">
          Loading matches...
        </div>
        <MatchList v-else :matches="upcomingMatches" />
      </section>

      <BaseModal
        :open="showCreateModal"
        title="Create Leaderboard"
        @close="closeCreateModal"
      >
        <div class="modal-content">
          <BaseInput
            v-model="newLeaderboardName"
            label="Leaderboard Name"
            placeholder="e.g., Office Cup 2026"
            :error="createError"
            @keyup.enter="createLeaderboard"
          />
        </div>
        <template #footer>
          <BaseButton variant="ghost" @click="closeCreateModal">
            Cancel
          </BaseButton>
          <BaseButton
            :loading="isCreating"
            :disabled="!newLeaderboardName.trim()"
            @click="createLeaderboard"
          >
            Create
          </BaseButton>
        </template>
      </BaseModal>
    </div>
  </div>
</template>

<style scoped>
.home-view {
  min-height: 100vh;
  background: var(--color-bg);
  padding: var(--space-6) var(--space-4);
}

.container {
  max-width: var(--max-width);
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: var(--space-12);
}

.welcome-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-4);
  flex-wrap: wrap;
}

.welcome-title {
  font-family: var(--font-display);
  font-size: var(--text-4xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.section-title {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  margin-bottom: var(--space-6);
}

.leaderboards-section,
.matches-section {
  display: flex;
  flex-direction: column;
}

.empty-state {
  padding: var(--space-12) var(--space-6);
  text-align: center;
  color: var(--color-text-secondary);
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
  border-radius: var(--radius-lg);
}

.empty-state p {
  font-size: var(--text-base);
  line-height: var(--leading-normal);
}

.empty-hint {
  margin-top: var(--space-2);
  font-size: var(--text-sm);
  opacity: 0.7;
}

.leaderboards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--space-4);
}

.leaderboard-card {
  cursor: pointer;
}

.leaderboard-name {
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  margin-bottom: var(--space-2);
}

.leaderboard-meta {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.loading-state {
  padding: var(--space-8);
  text-align: center;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

.modal-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

@media (max-width: 768px) {
  .welcome-title {
    font-size: var(--text-3xl);
  }

  .section-title {
    font-size: var(--text-xl);
  }

  .leaderboards-grid {
    grid-template-columns: 1fr;
  }
}
</style>
