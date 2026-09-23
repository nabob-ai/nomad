<template>
  <button
    class="nomad-button"
    :class="[variantClass, { 'nomad-button-disabled': disabled }]"
    :disabled="disabled"
    @click="handleClick"
  >
    <span v-if="icon" class="nomad-button-icon">{{ icon }}</span>
    <span class="nomad-button-label">{{ label }}</span>
  </button>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue';
import type { ButtonVariant, UIActionEvent, PropertyValue, DataMap } from '@/types/ui-protocol';
import { createUIActionEvent } from '@/types/ui-protocol';
import { useUIAction } from '@/composables/useUIAction';
import { resolveActionContext } from '@/protocol/client-message';

/**
 * NomadButton Component
 *
 * Button with variants and action emission.
 * Emits action events both as Vue events and to the Control channel.
 * Supports actionContext for automatic form data collection.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  label?: string;
  action?: string;
  variant?: ButtonVariant;
  disabled?: boolean;
  icon?: string;
  actionContext?: Record<string, PropertyValue>;
}

const props = withDefaults(defineProps<Props>(), {
  label: '',
  action: 'click',
  variant: 'primary',
  disabled: false,
});

const emit = defineEmits<{
  action: [event: UIActionEvent];
}>();

// Inject data model from parent surface (if available)
const dataModel = inject<DataMap>('surfaceDataModel', {});

// Setup UI action emitter
const { emitFullAction } = useUIAction({
  surfaceId: props.surfaceId,
  componentId: props.componentId,
});

const variantClass = computed(() => {
  const variantMap: Record<ButtonVariant, string> = {
    primary: 'nomad-button-primary',
    secondary: 'nomad-button-secondary',
    text: 'nomad-button-text',
  };
  return variantMap[props.variant] ?? 'nomad-button-primary';
});

function handleClick() {
  if (!props.disabled && props.componentId && props.surfaceId && props.action) {
    // Resolve action context if provided
    const context = props.actionContext
      ? resolveActionContext(props.actionContext, dataModel)
      : {};

    const event = createUIActionEvent(
      props.surfaceId,
      props.componentId,
      props.action,
      context,
    );

    emit('action', event);
    emitFullAction(event);
  }
}
</script>

<style scoped>
.nomad-button {
  @apply inline-flex items-center justify-center text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
  padding: var(--nomad-spacing-sm, 0.5rem) var(--nomad-spacing-md, 1rem);
  border-radius: var(--nomad-radius-lg, 0.5rem);
}

.nomad-button-primary {
  background-color: var(--nomad-primary, #3b82f6);
  color: var(--nomad-primary-contrast, #ffffff);
}

.nomad-button-primary:hover:not(:disabled) {
  background-color: var(--nomad-primary-hover, #2563eb);
}

.nomad-button-primary:focus {
  --tw-ring-color: var(--nomad-primary, #3b82f6);
}

.nomad-button-secondary {
  background-color: var(--nomad-secondary-light, #f3f4f6);
  color: var(--nomad-text, #111827);
}

.nomad-button-secondary:hover:not(:disabled) {
  background-color: var(--nomad-border, #e5e7eb);
}

.nomad-button-secondary:focus {
  --tw-ring-color: var(--nomad-secondary, #6b7280);
}

.dark .nomad-button-secondary {
  background-color: var(--nomad-surface, #334155);
  color: var(--nomad-text, #f1f5f9);
}

.dark .nomad-button-secondary:hover:not(:disabled) {
  background-color: var(--nomad-surface-hover, #475569);
}

.nomad-button-text {
  background-color: transparent;
  color: var(--nomad-primary, #3b82f6);
}

.nomad-button-text:hover:not(:disabled) {
  background-color: var(--nomad-primary-light, #dbeafe);
}

.nomad-button-text:focus {
  --tw-ring-color: var(--nomad-primary, #3b82f6);
}

.dark .nomad-button-text {
  color: var(--nomad-primary, #60a5fa);
}

.dark .nomad-button-text:hover:not(:disabled) {
  background-color: var(--nomad-surface, #334155);
}

.nomad-button-disabled {
  @apply opacity-50 cursor-not-allowed;
}

.nomad-button-icon {
  @apply mr-2;
}
</style>
