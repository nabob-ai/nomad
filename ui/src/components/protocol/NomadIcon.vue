<template>
  <span
    class="nomad-icon"
    :style="iconStyle"
    role="img"
    :aria-label="name"
  >
    {{ name }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

/**
 * NomadIcon Component
 *
 * Icon display using icon font or emoji.
 * For production, integrate with an icon library like Heroicons or Material Icons.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  name?: string;
  size?: string | number;
  color?: string;
}

const props = withDefaults(defineProps<Props>(), {
  name: '',
  size: 24,
});

const iconStyle = computed(() => {
  const style: Record<string, string> = {};

  if (props.size) {
    const sizeValue = typeof props.size === 'number' ? `${props.size}px` : props.size;
    style.fontSize = sizeValue;
    style.width = sizeValue;
    style.height = sizeValue;
  }

  if (props.color) {
    style.color = props.color;
  }

  return style;
});
</script>

<style scoped>
.nomad-icon {
  @apply inline-flex items-center justify-center;
}
</style>
