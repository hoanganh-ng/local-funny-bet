<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useLeaderboardStore } from '../store/leaderboard.store.js'
import { useAuthStore } from '../store/auth.store.js'
import BaseButton from '../components/base/BaseButton.vue'

const router = useRouter()
const leaderboardStore = useLeaderboardStore()
const authStore = useAuthStore()
const isLoading = ref(false)

onMounted(async () => {
  isLoading.value = true
  try {
    await leaderboardStore.fetchAll()
    await Promise.all(
      leaderboardStore.leaderboards.map(b => leaderboardStore.fetchBoardScores(b.id))
    )
  } catch (error) {
    console.error('Failed to load leaderboards:', error)
  } finally {
    isLoading.value = false
  }
})

function summary(boardId) {
  return leaderboardStore.boardSummary(boardId, authStore.user?.id)
}

const bestRank = computed(() => {
  let best = null
  for (const board of leaderboardStore.leaderboards) {
    const s = summary(board.id)
    if (s.userRank != null && (best === null || s.userRank < best)) {
      best = s.userRank
    }
  }
  return best
})

const totalMembers = computed(() =>
  leaderboardStore.leaderboards.reduce((sum, b) => sum + (summary(b.id).memberCount || 0), 0)
)

function goToBoard(id) {
  router.push(`/leaderboards/${id}`)
}

function createBoard() {
  router.push('/leaderboards/new')
}
</script>

<template>
  <div class="leaderboards-view">
    <div class="container">
      <header class="page-header">
        <div class="header-content">
          <p class="header-label">Where You Stand</p>
          <h1 class="page-title">Leaderboards</h1>
          <p class="page-subtitle">
            In {{ leaderboardStore.leaderboards.length }} {{ leaderboardStore.leaderboards.length === 1 ? 'board' : 'boards' }}, {{ totalMembers }} {{ totalMembers === 1 ? 'person' : 'people' }}{{ bestRank != null ? ` · best rank #${bestRank}` : '' }}
          </p>
        </div>
        <BaseButton @click="createBoard">
          + New board
        </BaseButton>
      </header>

      <div v-if="isLoading" class="loading-state">
        Loading boards...
      </div>

      <div v-else-if="leaderboardStore.leaderboards.length === 0" class="empty-state">
        <div class="empty-icon">🏆</div>
        <h2 class="empty-title">No boards yet</h2>
        <p class="empty-subtitle">
          Create one or ask a friend for an invite link.
        </p>
        <BaseButton size="lg" @click="createBoard">
          Create your first board
        </BaseButton>
      </div>

      <div v-else class="boards-grid">
        <div
          v-for="board in leaderboardStore.leaderboards"
          :key="board.id"
          class="board-card"
          @click="goToBoard(board.id)"
        >
          <div class="board-header">
            <div class="board-rank-chip" :class="{ highlight: summary(board.id).userRank === 1 }">
              Rank<br>
              {{ summary(board.id).userRank != null ? `#${summary(board.id).userRank}` : '—' }}
            </div>
            <div class="board-members">
              <div
                v-for="(member, idx) in summary(board.id).topMembers"
                :key="idx"
                class="member-avatar"
                :style="{ backgroundColor: `hsl(${idx * 120}, 60%, 50%)` }"
              >
                {{ member.name?.charAt(0) || '?' }}
              </div>
              <span v-if="summary(board.id).memberCount > 3" class="member-count">
                +{{ summary(board.id).memberCount - 3 }}
              </span>
            </div>
          </div>

          <div class="board-body">
            <h3 class="board-name">{{ board.name }}</h3>
            <p class="board-meta">{{ summary(board.id).memberCount }} {{ summary(board.id).memberCount === 1 ? 'member' : 'members' }}</p>
          </div>

          <div class="board-footer">
            <div class="board-stat">
              <span class="stat-label">Lead</span>
              <span class="stat-value">{{ summary(board.id).leaderName ?? '—' }}</span>
            </div>
            <div class="board-stat">
              <span class="stat-label">{{ summary(board.id).userPoints ?? 0 }} pts</span>
              <span
                class="stat-value"
                :class="{
                  positive: (summary(board.id).pointsFromLead ?? 0) > 0,
                  negative: (summary(board.id).pointsFromLead ?? 0) < 0
                }"
              >
                {{ (summary(board.id).pointsFromLead ?? 0) > 0 ? '+' : '' }}{{ summary(board.id).pointsFromLead ?? 0 }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.leaderboards-view {
  min-height: 100vh;
  padding: var(--space-8) var(--space-6);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--space-6);
  margin-bottom: var(--space-10);
}

.header-content {
  flex: 1;
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

.loading-state {
  padding: var(--space-12);
  text-align: center;
  color: var(--color-text-secondary);
}

.empty-state {
  padding: var(--space-16) var(--space-8);
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-6);
}

.empty-icon {
  font-size: 64px;
  opacity: 0.5;
}

.empty-title {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.empty-subtitle {
  font-size: var(--text-lg);
  color: var(--color-text-secondary);
  max-width: 400px;
}

.boards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--space-6);
}

.board-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
  padding: var(--space-6);
  background: var(--color-surface);
  border: var(--border-hairline);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

[data-theme="dark"] .board-card {
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

[data-theme="dark"] .board-card:hover {
  border-color: var(--color-border-hover);
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

[data-theme="light"] .board-card {
  border-radius: 0;
}

[data-theme="light"] .board-card:hover {
  border-width: 2px;
  border-color: var(--color-border-hover);
  padding: calc(var(--space-6) - 1px);
}

.board-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.board-rank-chip {
  padding: var(--space-2) var(--space-3);
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  color: var(--color-text-secondary);
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
  text-align: center;
  line-height: 1.4;
}

[data-theme="dark"] .board-rank-chip {
  border-radius: var(--radius-sm);
}

.board-rank-chip.highlight {
  background: var(--color-accent);
  border-color: var(--color-accent);
  color: var(--color-accent-text);
}

[data-theme="dark"] .board-rank-chip.highlight {
  box-shadow: var(--shadow-accent);
}

.board-members {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.member-avatar {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  color: white;
  border-radius: var(--radius-pill);
  border: 2px solid var(--color-bg);
}

[data-theme="light"] .member-avatar {
  border-radius: 0;
}

.member-count {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-disabled);
  margin-left: var(--space-1);
}

.board-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.board-name {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.board-meta {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.board-footer {
  display: flex;
  justify-content: space-between;
  padding-top: var(--space-4);
  border-top: var(--border-hairline);
}

.board-stat {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.stat-label {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.stat-value {
  font-family: var(--font-mono);
  font-size: var(--text-base);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.stat-value.positive {
  color: var(--color-correct);
}

.stat-value.negative {
  color: var(--color-danger);
}

@media (max-width: 768px) {
  .leaderboards-view {
    padding: var(--space-6) var(--space-4);
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
  }

  .page-title {
    font-size: var(--text-4xl);
  }

  .boards-grid {
    grid-template-columns: 1fr;
  }
}
</style>
