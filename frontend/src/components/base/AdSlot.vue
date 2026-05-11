<script setup>
defineProps({
  position: {
    type: String,
    required: true,
    validator: v => ['top', 'sidebar', 'between-matches'].includes(v)
  }
})

const adsEnabled = import.meta.env.VITE_ADS_ENABLED === 'true'
</script>

<template>
  <div v-if="adsEnabled" class="ad-slot" :data-position="position">
    <!-- AdSense injected here in global version -->
  </div>
</template>

<style scoped>
.ad-slot {
  display: v-bind("adsEnabled ? 'block' : 'none'");
  min-height: 100px;
  background: var(--color-surface);
  border: 1px dashed var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-disabled);
  font-family: var(--font-mono);
  font-size: var(--text-xs);
  text-transform: uppercase;
  letter-spacing: var(--tracking-wide);
}

.ad-slot::after {
  content: 'Ad Space';
}

/* ─── Dark theme ─── */
[data-theme="dark"] .ad-slot {
  border-radius: var(--radius-md);
}

/* ─── Light theme ─── */
[data-theme="light"] .ad-slot {
  border-radius: 0;
  border-style: solid;
}

/* ─── Position variants ─── */
.ad-slot[data-position="top"] {
  min-height: 90px;
}

.ad-slot[data-position="sidebar"] {
  min-height: 250px;
}

.ad-slot[data-position="between-matches"] {
  min-height: 120px;
}
</style>
