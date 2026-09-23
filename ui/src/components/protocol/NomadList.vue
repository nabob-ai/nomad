<template>
  <div class="nomad-list" :class="{ 'nomad-list-dividers': dividers }">
    <slot />
  </div>
</template>

<script setup lang="ts">
/**
 * NomadList Component
 *
 * List container with optional dividers between items.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  dividers?: boolean;
}

withDefaults(defineProps<Props>(), {
  dividers: false,
});
</script>

<style scoped>
.nomad-list {
  @apply flex flex-col;
}

.nomad-list-dividers > :deep(*:not(:last-child)) {
  @apply border-b border-gray-200 dark:border-gray-700;
}
</style>
