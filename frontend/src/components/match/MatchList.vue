<script setup>
import { computed, ref, onMounted, onUnmounted, watch } from 'vue'
import MatchCard from '../common/MatchCard.vue'
import BaseSpinner from '../base/BaseSpinner.vue'

const props = defineProps({
  matches:     { type: Array,   required: true },
  hasMore:     { type: Boolean, default: false },
  loadingMore: { type: Boolean, default: false },
})

const emit = defineEmits(['load-more'])

const sentinel = ref(null)
let observer = null

onMounted(() => {
  observer = new IntersectionObserver(
    (entries) => {
      if (entries[0].isIntersecting && props.hasMore && !props.loadingMore) {
        emit('load-more')
      }
    },
    { threshold: 0.1 }
  )
  if (sentinel.value) observer.observe(sentinel.value)
})

onUnmounted(() => {
  observer?.disconnect()
})

watch(() => props.hasMore, (newVal) => {
  if (newVal && sentinel.value) observer?.observe(sentinel.value)
})

// Filter out finished matches
const activeMatches = computed(() =>
  props.matches.filter(m => m.status !== 'finished')
)

// Live group
const liveMatches = computed(() =>
  activeMatches.value.filter(m => m.status === 'live')
)

const liveIds = computed(() => new Set(liveMatches.value.map(m => m.id)))

// Non-live scheduled matches
const scheduledMatches = computed(() =>
  activeMatches.value.filter(m => !liveIds.value.has(m.id))
)

// Today boundaries (recalculated once at render — fine for a list)
const todayStart = new Date()
todayStart.setHours(0, 0, 0, 0)
const todayEnd = new Date(todayStart)
todayEnd.setDate(todayEnd.getDate() + 1)

// Today's non-live scheduled matches
const todayMatches = computed(() =>
  scheduledMatches.value.filter(m => {
    const k = new Date(m.kickoffAt)
    return k >= todayStart && k < todayEnd
  })
)

// Format a Date as "YYYY-MM-DD" for grouping key
function dateKey(dateStr) {
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

// Future date groups (not today, not live)
const futureDateGroups = computed(() => {
  const future = scheduledMatches.value.filter(m => new Date(m.kickoffAt) >= todayEnd)

  const map = {}
  for (const m of future) {
    const key = dateKey(m.kickoffAt)
    if (!map[key]) map[key] = []
    map[key].push(m)
  }

  return Object.keys(map)
    .sort()
    .map(key => ({ key, matches: map[key] }))
})

// "Today · Tuesday, May 13" header
const todayHeader = computed(() => {
  const d = new Date()
  const dayName = d.toLocaleDateString('en-US', { weekday: 'long' })
  const month = d.toLocaleDateString('en-US', { month: 'long' })
  const day = d.getDate()
  return `Today · ${dayName}, ${month} ${day}`
})

// "Saturday, Jun 6" header for a date key "YYYY-MM-DD"
function futureDateHeader(key) {
  const [y, mo, d] = key.split('-').map(Number)
  const date = new Date(y, mo - 1, d)
  const dayName = date.toLocaleDateString('en-US', { weekday: 'long' })
  const month = date.toLocaleDateString('en-US', { month: 'short' })
  return `${dayName}, ${month} ${d}`
}

const hasAny = computed(() =>
  liveMatches.value.length > 0 ||
  todayMatches.value.length > 0 ||
  futureDateGroups.value.length > 0
)
</script>

<template>
  <div class="match-list">
    <div v-if="!hasAny" class="empty-state">
      No upcoming matches
    </div>

    <template v-else>
      <!-- LIVE group -->
      <section v-if="liveMatches.length > 0" class="group group--live">
        <h2 class="group-header">Live now</h2>
        <div class="group-cards">
          <MatchCard
            v-for="match in liveMatches"
            :key="match.id"
            :match="match"
          />
        </div>
      </section>

      <!-- TODAY group -->
      <section v-if="todayMatches.length > 0" class="group">
        <h2 class="group-header">{{ todayHeader }}</h2>
        <div class="group-cards">
          <MatchCard
            v-for="match in todayMatches"
            :key="match.id"
            :match="match"
          />
        </div>
      </section>

      <!-- DATE groups -->
      <section
        v-for="group in futureDateGroups"
        :key="group.key"
        class="group"
      >
        <h2 class="group-header">{{ futureDateHeader(group.key) }}</h2>
        <div class="group-cards">
          <MatchCard
            v-for="match in group.matches"
            :key="match.id"
            :match="match"
          />
        </div>
      </section>
    </template>

    <div ref="sentinel" class="scroll-sentinel" />

    <div v-if="loadingMore" class="loading-row">
      <BaseSpinner size="md" />
    </div>

    <div v-if="!hasMore && matches.length > 0" class="end-of-list">
      All matches loaded
    </div>
  </div>
</template>

<style scoped>
.match-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
}

.empty-state {
  padding: var(--space-8);
  text-align: center;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

/* ── Group ── */
.group {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.group-header {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  padding: var(--space-4) 0 var(--space-2);
  border-bottom: 1px solid var(--color-border);
  margin: 0;
}

.group-cards {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

/* ── Live group accent border ── */
.group--live {
  border-left: 2px solid var(--color-status-live);
  padding-left: var(--space-4);
  animation: live-border-pulse 1.5s ease-in-out infinite;
}

[data-theme="light"] .group--live {
  border-left-color: var(--color-status-live);
  animation: none;
}

@keyframes live-border-pulse {
  0%, 100% { border-left-color: var(--color-status-live); }
  50%       { border-left-color: var(--color-status-live-dim); }
}

/* ── Scroll sentinel ── */
.scroll-sentinel {
  height: 1px;
}

/* ── Loading row ── */
.loading-row {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-6) 0;
  min-height: var(--space-16);
}

/* ── End of list ── */
.end-of-list {
  text-align: center;
  padding: var(--space-6) 0;
  color: var(--color-text-disabled);
  font-size: var(--text-sm);
}
</style>
