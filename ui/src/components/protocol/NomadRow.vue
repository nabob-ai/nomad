<template>
  <div
    class="nomad-row"
    :class="[alignClass, wrapClass]"
    :style="{ gap: gapStyle }"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Alignment } from '@/types/ui-protocol';

/**
 * NomadRow Component
 *
 * Horizontal layout container with configurable gap and alignment.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  gap?: number;
  align?: Alignment;
  wrap?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  gap: 8,
  align: 'start',
  wrap: false,
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

const wrapClass = computed(() => props.wrap ? 'flex-wrap' : 'flex-nowrap');

const gapStyle = computed(() => `${props.gap}px`);
</script>

<style scoped>
.nomad-row {
  @apply flex flex-row;
}
</style>
