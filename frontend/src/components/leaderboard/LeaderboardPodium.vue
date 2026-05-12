<script setup>
import BaseAvatar from '../base/BaseAvatar.vue'
import { computed } from 'vue'

const props = defineProps({
  scores: {
    type: Array,
    default: () => []
  }
})

const topThree = computed(() => {
  return props.scores.slice(0, 3)
})

const first = computed(() => topThree.value[0] || null)
const second = computed(() => topThree.value[1] || null)
const third = computed(() => topThree.value[2] || null)

function getMedalIcon(rank) {
  const medals = { 1: '🥇', 2: '🥈', 3: '🥉' }
  return medals[rank] || rank
}
</script>

<template>
  <div class="leaderboard-podium">
    <div v-if="topThree.length === 0" class="empty-state">
      No predictions yet. First to predict gets the podium!
    </div>
    <div v-else class="podium-container">
      <!-- Second Place -->
      <div v-if="second" class="podium-place" :class="{ rank2: true }">
        <div class="podium-rank">
          <span class="medal-icon">{{ getMedalIcon(2) }}</span>
        </div>
        <BaseAvatar
          :src="second.user.avatarUrl"
          :name="second.user.name"
          size="lg"
        />
        <div class="podium-name">{{ second.user.name }}</div>
        <div class="podium-points">{{ second.points }} pts</div>
      </div>

      <!-- First Place -->
      <div v-if="first" class="podium-place" :class="{ rank1: true }">
        <div class="podium-rank">
          <span class="medal-icon">{{ getMedalIcon(1) }}</span>
        </div>
        <BaseAvatar
          :src="first.user.avatarUrl"
          :name="first.user.name"
          size="xl"
        />
        <div class="podium-name">{{ first.user.name }}</div>
        <div class="podium-points">{{ first.points }} pts</div>
      </div>

      <!-- Third Place -->
      <div v-if="third" class="podium-place" :class="{ rank3: true }">
        <div class="podium-rank">
          <span class="medal-icon">{{ getMedalIcon(3) }}</span>
        </div>
        <BaseAvatar
          :src="third.user.avatarUrl"
          :name="third.user.name"
          size="lg"
        />
        <div class="podium-name">{{ third.user.name }}</div>
        <div class="podium-points">{{ third.points }} pts</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.leaderboard-podium {
  padding: var(--space-8) var(--space-4);
  background: var(--color-surface);
  border: var(--border-hairline);
}

[data-theme="dark"] .leaderboard-podium {
  border-radius: var(--radius-lg);
}

.empty-state {
  padding: var(--space-8);
  text-align: center;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

.podium-container {
  display: flex;
  justify-content: center;
  align-items: flex-end;
  gap: var(--space-6);
  max-width: 600px;
  margin: 0 auto;
}

.podium-place {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-6);
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
  flex: 1;
  transition: all var(--duration-fast) var(--ease-default);
}

[data-theme="dark"] .podium-place {
  border-radius: var(--radius-lg);
}

/* First place special styling */
.podium-place.rank1 {
  order: 2;
  padding-top: var(--space-8);
  padding-bottom: var(--space-8);
  background: var(--color-accent);
  border-color: var(--color-accent);
}

[data-theme="dark"] .podium-place.rank1 {
  box-shadow: var(--shadow-accent);
}

[data-theme="light"] .podium-place.rank1 {
  background: var(--color-pick-bg);
  border-width: 2px;
  border-color: var(--color-border-strong);
}

/* Second place */
.podium-place.rank2 {
  order: 1;
}

/* Third place */
.podium-place.rank3 {
  order: 3;
}

.podium-rank {
  font-size: var(--text-2xl);
  line-height: 1;
}

.medal-icon {
  display: block;
}

.podium-name {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  text-align: center;
  word-break: break-word;
}

.podium-place.rank1 .podium-name {
  color: var(--color-accent-text);
  font-size: var(--text-base);
}

.podium-points {
  font-family: var(--font-mono);
  font-size: var(--text-xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.podium-place.rank1 .podium-points {
  color: var(--color-accent-text);
  font-size: var(--text-2xl);
}

@media (max-width: 768px) {
  .podium-container {
    gap: var(--space-4);
  }

  .podium-place {
    padding: var(--space-4);
    gap: var(--space-2);
  }

  .podium-place.rank1 {
    padding-top: var(--space-6);
    padding-bottom: var(--space-6);
  }

  .podium-rank {
    font-size: var(--text-xl);
  }

  .podium-name {
    font-size: var(--text-xs);
  }

  .podium-place.rank1 .podium-name {
    font-size: var(--text-sm);
  }

  .podium-points {
    font-size: var(--text-base);
  }

  .podium-place.rank1 .podium-points {
    font-size: var(--text-xl);
  }
}
</style>
