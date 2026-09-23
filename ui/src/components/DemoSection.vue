<template>
  <div class="demo-section">
    <div v-if="title" class="demo-section-title">{{ title }}</div>
    <div class="demo-section-content">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title?: string;
}

defineProps<Props>();
</script>

<style scoped>
.demo-section {
  @apply mb-6 last:mb-0;
}

.demo-section-title {
  @apply text-sm font-medium text-gray-700 dark:text-gray-300 mb-3;
}

.demo-section-content {
  @apply p-4 bg-gray-50 dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-700;
}
</style>
