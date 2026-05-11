<script setup>
import { ref, computed } from 'vue'
import OutcomePicker from './OutcomePicker.vue'
import BaseButton from '../base/BaseButton.vue'
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
  const kickoffTime = new Date(props.match.kickoffAt)
  return kickoffTime < new Date()
})

const canSubmit = computed(() => {
  return selectedOutcome.value && !isSubmitting.value && !isLocked.value
})

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
  <div class="prediction-form">
    <div v-if="isLocked" class="locked-state">
      <OutcomePicker
        :model-value="currentPrediction?.outcome"
        :disabled="true"
      />
      <p class="locked-label">Locked</p>
    </div>
    <div v-else class="active-form">
      <OutcomePicker
        v-model="selectedOutcome"
        :disabled="isSubmitting"
      />
      <BaseButton
        :disabled="!canSubmit"
        :loading="isSubmitting"
        @click="submitPrediction"
      >
        Submit Prediction
      </BaseButton>
      <p v-if="error" class="error-message">{{ error }}</p>
    </div>
  </div>
</template>

<style scoped>
.prediction-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.locked-state,
.active-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.locked-label {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  text-align: center;
}

.error-message {
  font-size: var(--text-sm);
  color: var(--color-danger);
  text-align: center;
}
</style>
