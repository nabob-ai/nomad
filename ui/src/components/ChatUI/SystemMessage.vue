<template>
  <div class="system-message">
    <div class="system-content">
      {{ content }}
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  content: string;
}

defineProps<Props>();
</script>

<style scoped>
.system-message {
  @apply flex justify-center my-4;
}

.system-content {
  @apply px-4 py-2 bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 text-xs rounded-full;
}
</style>
