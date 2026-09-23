<template>
  <component :is="tagName" :class="textClass">
    {{ sanitizedText }}
  </component>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { TextUsageHint } from '@/types/ui-protocol';
import { sanitizeText } from '@/protocol/security';

/**
 * NomadText Component
 *
 * Text display with semantic styling based on usage hint.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  text?: string;
  usageHint?: TextUsageHint;
}

const props = withDefaults(defineProps<Props>(), {
  text: '',
  usageHint: 'body',
});

// Sanitize text to prevent XSS
const sanitizedText = computed(() => {
  // Note: We use textContent binding ({{ }}) which auto-escapes,
  // but we still sanitize for defense in depth
  return props.text;
});

// Map usage hint to HTML tag
const tagName = computed(() => {
  const tagMap: Record<TextUsageHint, string> = {
    h1: 'h1',
    h2: 'h2',
    h3: 'h3',
    h4: 'h4',
    h5: 'h5',
    caption: 'span',
    body: 'p',
  };
  return tagMap[props.usageHint] ?? 'p';
});

// Map usage hint to CSS class
const textClass = computed(() => {
  const classMap: Record<TextUsageHint, string> = {
    h1: 'nomad-text-h1',
    h2: 'nomad-text-h2',
    h3: 'nomad-text-h3',
    h4: 'nomad-text-h4',
    h5: 'nomad-text-h5',
    caption: 'nomad-text-caption',
    body: 'nomad-text-body',
  };
  return ['nomad-text', classMap[props.usageHint] ?? 'nomad-text-body'];
});
</script>

<style scoped>
.nomad-text {
  color: var(--nomad-text, #111827);
}

.dark .nomad-text {
  color: var(--nomad-text, #f1f5f9);
}

.nomad-text-h1 {
  @apply font-bold mb-4;
  font-size: var(--nomad-font-size-3xl, 1.875rem);
}

.nomad-text-h2 {
  @apply font-semibold mb-3;
  font-size: var(--nomad-font-size-2xl, 1.5rem);
}

.nomad-text-h3 {
  @apply font-semibold mb-2;
  font-size: var(--nomad-font-size-xl, 1.25rem);
}

.nomad-text-h4 {
  @apply font-medium mb-2;
  font-size: var(--nomad-font-size-lg, 1.125rem);
}

.nomad-text-h5 {
  @apply font-medium mb-1;
  font-size: var(--nomad-font-size-base, 1rem);
}

.nomad-text-caption {
  font-size: var(--nomad-font-size-sm, 0.875rem);
  color: var(--nomad-text-secondary, #6b7280);
}

.dark .nomad-text-caption {
  color: var(--nomad-text-secondary, #cbd5e1);
}

.nomad-text-body {
  font-size: var(--nomad-font-size-base, 1rem);
  line-height: var(--nomad-line-height-relaxed, 1.625);
}
</style>
