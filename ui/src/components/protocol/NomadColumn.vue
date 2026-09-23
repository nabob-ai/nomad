<template>
  <div
    class="nomad-column"
    :class="alignClass"
    :style="{ gap: gapStyle }"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Alignment } from '@/types/ui-protocol';

/**
 * NomadColumn Component
 *
 * Vertical layout container with configurable gap and alignment.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  gap?: number;
  align?: Alignment;
}

const props = withDefaults(defineProps<Props>(), {
  gap: 8,
  align: 'stretch',
});

const alignClass = computed(() => {
  const alignMap: Record<Alignment, string> = {
    start: 'items-start',
    center: 'items-center',
    end: 'items-end',
    stretch: 'items-stretch',
  };
  return alignMap[props.align];
});

const gapStyle = computed(() => `${props.gap}px`);
</script>

<style scoped>
.nomad-column {
  @apply flex flex-col;
}
</style>
