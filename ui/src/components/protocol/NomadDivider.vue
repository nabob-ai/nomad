<template>
  <hr
    class="nomad-divider"
    :class="orientationClass"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { DividerOrientation } from '@/types/ui-protocol';

/**
 * NomadDivider Component
 *
 * Visual separator with horizontal or vertical orientation.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  orientation?: DividerOrientation;
}

const props = withDefaults(defineProps<Props>(), {
  orientation: 'horizontal',
});

const orientationClass = computed(() => {
  return props.orientation === 'vertical'
    ? 'nomad-divider-vertical'
    : 'nomad-divider-horizontal';
});
</script>

<style scoped>
.nomad-divider {
  @apply border-0 bg-gray-200 dark:bg-gray-700;
}

.nomad-divider-horizontal {
  @apply w-full h-px my-2;
}

.nomad-divider-vertical {
  @apply h-full w-px mx-2;
}
</style>
