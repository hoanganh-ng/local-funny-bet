<script setup>
import { computed } from 'vue'

const props = defineProps({
  src: { type: String, default: '' },
  name: { type: String, required: true },
  size: {
    type: String,
    default: 'md',
    validator: v => ['sm', 'md', 'lg', 'xl'].includes(v)
  }
})

const initials = computed(() => {
  const names = props.name.trim().split(' ')
  if (names.length >= 2) {
    return (names[0][0] + names[names.length - 1][0]).toUpperCase()
  }
  return names[0].slice(0, 2).toUpperCase()
})

const imageError = computed(() => !props.src)
</script>

<template>
  <div
    class="avatar"
    :class="`avatar--${size}`"
    :aria-label="`Avatar for ${name}`"
  >
    <img
      v-if="src && !imageError"
      :src="src"
      :alt="name"
      class="avatar__image"
      @error="imageError = true"
    />
    <span v-else class="avatar__initials">
      {{ initials }}
    </span>
  </div>
</template>

<style scoped>
.avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-family: var(--font-body);
  font-weight: var(--font-semibold);
  overflow: hidden;
  user-select: none;
}

.avatar__image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar__initials {
  line-height: 1;
}

/* ─── Size variants ─── */
.avatar--sm {
  width: 32px;
  height: 32px;
  font-size: var(--text-xs);
}

.avatar--md {
  width: 40px;
  height: 40px;
  font-size: var(--text-sm);
}

.avatar--lg {
  width: 56px;
  height: 56px;
  font-size: var(--text-base);
}

.avatar--xl {
  width: 80px;
  height: 80px;
  font-size: var(--text-lg);
}

/* ─── Dark theme ─── */
[data-theme="dark"] .avatar {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
  border-radius: var(--radius-full);
}

/* ─── Light theme ─── */
[data-theme="light"] .avatar {
  background: #000;
  border: 2px solid #000;
  color: #fff;
  border-radius: 0;
}
</style>
