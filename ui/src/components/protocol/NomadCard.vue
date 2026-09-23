<template>
  <div
    class="nomad-card"
    :class="{ 'nomad-card-clickable': clickable }"
    @click="handleClick"
  >
    <div v-if="title || subtitle" class="nomad-card-header">
      <h3 v-if="title" class="nomad-card-title">{{ title }}</h3>
      <p v-if="subtitle" class="nomad-card-subtitle">{{ subtitle }}</p>
    </div>
    <div class="nomad-card-content">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { UIActionEvent } from '@/types/ui-protocol';
import { createUIActionEvent } from '@/types/ui-protocol';
import { useUIAction } from '@/composables/useUIAction';

/**
 * NomadCard Component
 *
 * Card container with optional title, subtitle, and click action.
 * Emits action events both as Vue events and to the Control channel.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  title?: string;
  subtitle?: string;
  clickable?: boolean;
  action?: string;
}

const props = withDefaults(defineProps<Props>(), {
  clickable: false,
});

const emit = defineEmits<{
  action: [event: UIActionEvent];
}>();

// Setup UI action emitter
const { emitFullAction } = useUIAction({
  surfaceId: props.surfaceId,
  componentId: props.componentId,
});

function handleClick() {
  if (props.clickable && props.action && props.componentId && props.surfaceId) {
    const event = createUIActionEvent(
      props.surfaceId,
      props.componentId,
      props.action,
      {},
    );
    emit('action', event);
    emitFullAction(event);
  }
}
</script>

<style scoped>
.nomad-card {
  @apply overflow-hidden;
  background-color: var(--nomad-surface, #ffffff);
  border: 1px solid var(--nomad-border, #e5e7eb);
  border-radius: var(--nomad-radius-lg, 0.5rem);
  box-shadow: var(--nomad-shadow-sm, 0 1px 2px 0 rgb(0 0 0 / 0.05));
}

.dark .nomad-card {
  background-color: var(--nomad-surface, #334155);
  border-color: var(--nomad-border, #475569);
}

.nomad-card-clickable {
  @apply cursor-pointer transition-shadow;
}

.nomad-card-clickable:hover {
  box-shadow: var(--nomad-shadow-md, 0 4px 6px -1px rgb(0 0 0 / 0.1));
}

.nomad-card-header {
  @apply px-4 py-3;
  border-bottom: 1px solid var(--nomad-border, #e5e7eb);
}

.dark .nomad-card-header {
  border-bottom-color: var(--nomad-border, #475569);
}

.nomad-card-title {
  @apply text-base font-semibold m-0;
  color: var(--nomad-text, #111827);
}

.dark .nomad-card-title {
  color: var(--nomad-text, #f1f5f9);
}

.nomad-card-subtitle {
  @apply text-sm mt-1 m-0;
  color: var(--nomad-text-secondary, #6b7280);
}

.dark .nomad-card-subtitle {
  color: var(--nomad-text-secondary, #cbd5e1);
}

.nomad-card-content {
  padding: var(--nomad-spacing-md, 1rem);
}
</style>
