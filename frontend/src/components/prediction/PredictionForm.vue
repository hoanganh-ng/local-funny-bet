<script setup>
import { ref, computed } from 'vue'
import OutcomePicker from './OutcomePicker.vue'
import BaseButton from '../base/BaseButton.vue'
import MatchStatusBadge from '../match/MatchStatusBadge.vue'
import { predictionService } from '../../services/prediction.service.js'

const props = defineProps({
  match: {
    type: Object,
    required: true
  },
  currentPrediction: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['success'])

const selectedOutcome = ref(props.currentPrediction?.outcome || null)
const isSubmitting = ref(false)
const error = ref(null)

const isLocked = computed(() => {
  if (props.match.status === 'live' || props.match.status === 'finished') return true
  return new Date(props.match.kickoffAt) <= new Date()
})

const showScore = computed(() =>
  props.match.status === 'finished' || props.match.status === 'live'
)

const formattedTime = computed(() => {
  const d = new Date(props.match.kickoffAt)
  const months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec']
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${months[d.getMonth()]} ${d.getDate()} · ${hh}:${mm}`
})

const canSubmit = computed(() =>
  selectedOutcome.value && !isSubmitting.value && !isLocked.value
)

const submitPrediction = async () => {
  if (!canSubmit.value) return
  isSubmitting.value = true
  error.value = null
  try {
    await predictionService.upsert(props.match.id, selectedOutcome.value)
    emit('success')
  } catch (err) {
    error.value = err.message || 'Failed to save prediction'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="prediction-card">
    <!-- Match info -->
    <div class="card-match">
      <div class="match-header">
        <span class="match-time">{{ formattedTime }}</span>
        <MatchStatusBadge :status="match.status" />
      </div>

      <div class="match-teams">
        <div class="team">
          <span v-if="match.homeTeamCode" class="team-code">{{ match.homeTeamCode }}</span>
          <span class="team-name">{{ match.homeTeam }}</span>
          <span v-if="showScore" class="score">{{ match.homeScore ?? '—' }}</span>
        </div>
        <div class="teams-sep">
          <span v-if="!showScore" class="vs-label">vs</span>
          <span v-else class="score-sep">—</span>
        </div>
        <div class="team team--away">
          <span v-if="showScore" class="score">{{ match.awayScore ?? '—' }}</span>
          <span class="team-name">{{ match.awayTeam }}</span>
          <span v-if="match.awayTeamCode" class="team-code">{{ match.awayTeamCode }}</span>
        </div>
      </div>
    </div>

    <div class="card-divider" />

    <!-- Picker row -->
    <div class="card-picker">
      <div class="picker-row">
        <OutcomePicker
          v-model="selectedOutcome"
          :locked="isLocked"
          class="picker"
        />
        <div class="picker-submit">
          <BaseButton
            v-if="!isLocked"
            variant="primary"
            :disabled="!canSubmit"
            :loading="isSubmitting"
            @click="submitPrediction"
          >
            Submit
          </BaseButton>
          <span v-else class="locked-label">Locked</span>
        </div>
      </div>
      <p v-if="error" class="error-message">{{ error }}</p>
    </div>
  </div>
</template>

<style scoped>
.prediction-card {
  display: flex;
  flex-direction: column;
  background: var(--color-surface);
  border: var(--border-hairline);
  transition: all var(--duration-fast) var(--ease-default);
}

[data-theme="dark"] .prediction-card {
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

[data-theme="dark"] .prediction-card:hover {
  border-color: var(--color-border-hover);
  box-shadow: var(--shadow-md);
}

[data-theme="light"] .prediction-card {
  border-radius: 0;
}

/* ── Match info section ── */
.card-match {
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

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

.match-teams {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  gap: var(--space-3);
}

.team {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.team--away {
  flex-direction: row-reverse;
}

.team-code {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  padding: var(--space-1) var(--space-2);
  background: rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
}

[data-theme="dark"] .team-code {
  border-radius: var(--radius-sm);
}

[data-theme="light"] .team-code {
  background: rgba(0, 0, 0, 0.06);
}

.team-name {
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
}

.score {
  font-size: var(--text-2xl);
  font-weight: var(--font-black);
  color: var(--color-text-primary);
  font-family: var(--font-mono);
}

.teams-sep {
  text-align: center;
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
  color: var(--color-border-strong);
  font-weight: var(--font-light);
}

/* ── Divider ── */
.card-divider {
  height: 1px;
  background: var(--color-border);
}

/* ── Picker section ── */
.card-picker {
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.picker-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.picker {
  flex: 1;
}

.picker-submit {
  flex-shrink: 0;
}

.locked-label {
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-text-disabled);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.error-message {
  font-size: var(--text-sm);
  color: var(--color-danger);
}

/* ── Mobile: stack picker row ── */
@media (max-width: 767px) {
  .picker-row {
    flex-direction: column;
    align-items: stretch;
  }

  .picker-submit {
    display: flex;
    justify-content: stretch;
  }

  .picker-submit .btn {
    width: 100%;
  }
}
</style>
