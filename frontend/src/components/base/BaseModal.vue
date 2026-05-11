<script setup>
import { watch } from 'vue'

const props = defineProps({
  open: { type: Boolean, required: true },
  title: { type: String, default: '' }
})

const emit = defineEmits(['close'])

watch(() => props.open, (isOpen) => {
  if (isOpen) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})

function closeModal() {
  emit('close')
}

function handleBackdropClick(event) {
  if (event.target === event.currentTarget) {
    closeModal()
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="open" class="modal-backdrop" @click="handleBackdropClick">
        <div class="modal" role="dialog" aria-modal="true" :aria-label="title">
          <div v-if="title || $slots.header" class="modal__header">
            <slot name="header">
              <h2 class="modal__title">{{ title }}</h2>
            </slot>
            <button
              class="modal__close"
              aria-label="Close modal"
              @click="closeModal"
            >
              ✕
            </button>
          </div>
          <div class="modal__body">
            <slot />
          </div>
          <div v-if="$slots.footer" class="modal__footer">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-4);
}

.modal {
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
}

.modal__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: var(--space-4);
  margin-bottom: var(--space-6);
}

.modal__title {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-text-primary);
}

.modal__close {
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  font-size: var(--text-2xl);
  cursor: pointer;
  padding: var(--space-2);
  line-height: 1;
  transition: color var(--duration-fast);
}

.modal__close:hover {
  color: var(--color-text-primary);
}

.modal__body {
  color: var(--color-text-primary);
}

.modal__footer {
  margin-top: var(--space-6);
  padding-top: var(--space-4);
  display: flex;
  gap: var(--space-3);
  justify-content: flex-end;
}

/* ─── Dark theme ─── */
[data-theme="dark"] .modal-backdrop {
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(4px);
}

[data-theme="dark"] .modal {
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  padding: var(--space-8);
}

[data-theme="dark"] .modal__header {
  border-bottom: 1px solid var(--color-border);
}

[data-theme="dark"] .modal__footer {
  border-top: 1px solid var(--color-border);
}

/* ─── Light theme ─── */
[data-theme="light"] .modal-backdrop {
  background: rgba(0, 0, 0, 0.6);
}

[data-theme="light"] .modal {
  background: var(--color-bg);
  border: var(--border-medium);
  border-radius: 0;
  box-shadow: none;
  padding: var(--space-12);
}

[data-theme="light"] .modal__header {
  border-bottom: var(--border-medium);
}

[data-theme="light"] .modal__footer {
  border-top: var(--border-thin);
}

[data-theme="light"] .modal__title {
  letter-spacing: var(--tracking-tight);
}

/* ─── Transitions ─── */
.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--duration-base) var(--ease-default);
}

.modal-enter-active .modal,
.modal-leave-active .modal {
  transition: transform var(--duration-base) var(--ease-default);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal,
.modal-leave-to .modal {
  transform: scale(0.95) translateY(20px);
}

[data-theme="light"] .modal-enter-from .modal,
[data-theme="light"] .modal-leave-to .modal {
  transform: scale(1) translateY(10px);
}
</style>
