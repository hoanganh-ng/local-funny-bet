<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import MatchStatusBadge from '../match/MatchStatusBadge.vue'
import OutcomePicker from '../prediction/OutcomePicker.vue'
import { predictionService } from '../../services/prediction.service.js'

const props = defineProps({
  match: {
    type: Object,
    required: true
  }
})

const prediction = ref(props.match.userPrediction || null)
const isSaving = ref(false)

const showScore = computed(() => {
  return props.match.status === 'finished' || props.match.status === 'live'
})

const isLocked = computed(() => {
  if (props.match.status === 'live' || props.match.status === 'finished') {
    return true
  }
  return new Date(props.match.kickoffAt) <= new Date()
})

// Fix 1: "Jun 12 · 02:00" format
const formattedTime = computed(() => {
  const d = new Date(props.match.kickoffAt)
  const months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec']
  const month = months[d.getMonth()]
  const day = d.getDate()
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${month} ${day} · ${hh}:${mm}`
})

const timeUntilKickoff = ref(null)
const showCountdown = computed(() => {
  if (!timeUntilKickoff.value) return false
  return props.match.status === 'scheduled' && timeUntilKickoff.value < 4 * 60 * 60 * 1000
})

const formatCountdown = computed(() => {
  if (!timeUntilKickoff.value) return ''
  const hours = Math.floor(timeUntilKickoff.value / (60 * 60 * 1000))
  const minutes = Math.floor((timeUntilKickoff.value % (60 * 60 * 1000)) / (60 * 1000))
  return `${hours}h ${minutes}m`
})

let countdownInterval = null

onMounted(() => {
  updateCountdown()
  countdownInterval = setInterval(updateCountdown, 60000)
})

onUnmounted(() => {
  if (countdownInterval) clearInterval(countdownInterval)
})

function updateCountdown() {
  const kickoff = new Date(props.match.kickoffAt)
  const now = new Date()
  timeUntilKickoff.value = kickoff - now
}

async function handlePredictionChange(newPrediction) {
  if (isLocked.value || isSaving.value) return

  prediction.value = newPrediction
  isSaving.value = true

  try {
    await predictionService.upsert(props.match.id, newPrediction)
  } catch (error) {
    console.error('Failed to save prediction:', error)
    prediction.value = props.match.userPrediction || null
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div class="match-card">
    <!-- Header: time + badge -->
    <div class="match-header">
      <span class="match-time">{{ formattedTime }}</span>
      <MatchStatusBadge :status="match.status" />
    </div>

    <!-- Teams and score -->
    <div class="match-teams">
      <div class="team home-team">
        <!-- Fix 2: team code badge from API field -->
        <span v-if="match.homeTeamCode" class="team-code">{{ match.homeTeamCode }}</span>
        <span class="team-name">{{ match.homeTeam }}</span>
        <span v-if="showScore" class="score">{{ match.homeScore ?? '-' }}</span>
      </div>

      <!-- Fix 3: compact VS divider -->
      <div class="score-divider">
        <span v-if="!showScore" class="vs-label">vs</span>
        <span v-else class="score-sep">-</span>
      </div>

      <div class="team away-team">
        <span v-if="match.awayTeamCode" class="team-code">{{ match.awayTeamCode }}</span>
        <span class="team-name">{{ match.awayTeam }}</span>
        <span v-if="showScore" class="score">{{ match.awayScore ?? '-' }}</span>
      </div>
    </div>

    <!-- Prediction picker -->
    <OutcomePicker
      v-model="prediction"
      :locked="isLocked"
      @update:model-value="handlePredictionChange"
    />

    <!-- Lock countdown -->
    <div v-if="showCountdown" class="countdown-bar">
      <span class="countdown-icon">🔒</span>
      <span class="countdown-text">Locks in {{ formatCountdown }}</span>
    </div>

    <!-- Locked message -->
    <div v-else-if="isLocked && !match.result" class="locked-message">
      Locked at kickoff · result settles soon
    </div>
  </div>
</template>

<style scoped>
.match-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-4);
  background: var(--color-surface);
  border: var(--border-hairline);
  transition: all var(--duration-fast) var(--ease-default);
}

[data-theme="dark"] .match-card {
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

[data-theme="dark"] .match-card:hover {
  border-color: var(--color-border-hover);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

[data-theme="light"] .match-card {
  border-radius: 0;
}

/* ── Header ── */
.match-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.match-time {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  letter-spacing: var(--tracking-wide);
  text-transform: uppercase;
}

/* ── Teams — Fix 3: tighter gap ── */
.match-teams {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.team {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.home-team,
.away-team {
  flex-direction: row;
}

/* Fix 2: team code badge */
.team-code {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  padding: var(--space-1) var(--space-2);
  background: rgba(255, 255, 255, 0.06);
  border-radius: var(--radius-sm);
  flex-shrink: 0;
}

[data-theme="light"] .team-code {
  background: rgba(0, 0, 0, 0.06);
}

.team-name {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  flex: 1;
}

.score {
  font-size: var(--text-4xl);
  font-weight: var(--font-black);
  color: var(--color-text-primary);
  font-family: var(--font-mono);
  min-width: var(--space-12);
  text-align: center;
  letter-spacing: var(--tracking-tight);
}

/* Fix 3: no extra vertical padding on divider */
.score-divider {
  text-align: center;
  margin: 0;
}

.vs-label {
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-text-disabled);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.score-sep {
  font-size: var(--text-xl);
  color: var(--color-border);
  font-weight: var(--font-light);
}

/* ── Countdown ── */
.countdown-bar {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  background: var(--color-bg-elevated);
  border: var(--border-hairline);
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-warning);
}

[data-theme="dark"] .countdown-bar {
  border-radius: var(--radius-sm);
  background: var(--color-warning-bg);
  border-color: var(--color-warning-border);
}

.countdown-icon {
  font-size: var(--text-sm);
}

.countdown-text {
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
  font-weight: var(--font-medium);
}

/* ── Locked message ── */
.locked-message {
  text-align: center;
  padding: var(--space-2);
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-disabled);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

/* Fix 4: desktop button sizing — delegated to OutcomePicker via CSS custom props */
@media (min-width: 769px) {
  .match-card :deep(.outcome-btn) {
    height: var(--space-12);
    padding-top: 0;
    padding-bottom: 0;
  }

  .match-card :deep(.outcome-sublabel) {
    font-size: var(--text-sm);
  }
}
</style>
