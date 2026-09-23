<template>
  <Teleport to="body">
    <Transition name="nomad-modal">
      <div
        v-if="openValue"
        class="nomad-modal-overlay"
        @click.self="handleOverlayClick"
      >
        <div class="nomad-modal" role="dialog" aria-modal="true">
          <div v-if="titleValue || closable" class="nomad-modal-header">
            <h2 v-if="titleValue" class="nomad-modal-title">{{ titleValue }}</h2>
            <button
              v-if="closable"
              class="nomad-modal-close"
              aria-label="Close"
              @click="handleClose"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div class="nomad-modal-content">
            <slot />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { DataValue, UIActionEvent, PropertyValue } from '@/types/ui-protocol';
import { createUIActionEvent } from '@/types/ui-protocol';
import { useNomadDataBinding } from '@/composables/useNomadDataBinding';
import { useUIAction } from '@/composables/useUIAction';

/**
 * NomadModal Component
 *
 * Modal dialog with overlay, title, and close functionality.
 * Uses useNomadDataBinding for reactive data binding with the data model.
 * Emits action events both as Vue events and to the Control channel.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  /** Open property - can be literal or path reference */
  open?: PropertyValue | boolean;
  /** Title property - can be literal or path reference */
  title?: PropertyValue | string;
  closable?: boolean;
  closeAction?: string;
  /** Path for two-way binding (legacy support) */
  openPath?: string;
}

const props = withDefaults(defineProps<Props>(), {
  closable: true,
});

const emit = defineEmits<{
  action: [event: UIActionEvent];
  'update:value': [path: string, value: DataValue];
}>();

// Setup UI action emitter
const { emitFullAction } = useUIAction({
  surfaceId: props.surfaceId,
  componentId: props.componentId,
});

// Convert props to PropertyValue format for binding
function toPropertyValue(value: PropertyValue | string | boolean | undefined): PropertyValue | undefined {
  if (value === undefined) return undefined;
  if (typeof value === 'string') return { literalString: value };
  if (typeof value === 'boolean') return { literalBoolean: value };
  return value as PropertyValue;
}

// Use data binding for open
const { value: boundOpen, updateValue, path: openPath } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.open),
  defaultValue: false,
});

// Use data binding for title
const { value: boundTitle } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.title),
  defaultValue: '',
});

// Computed values for template
const openValue = computed(() => Boolean(boundOpen.value));
const titleValue = computed(() => boundTitle.value ? String(boundTitle.value) : undefined);

function handleClose() {
  // Update via data binding to close the modal
  updateValue(false);

  // Also emit for legacy support
  const path = openPath ?? props.openPath;
  if (path) {
    emit('update:value', path, false);
  }

  // Emit action event
  if (props.componentId && props.surfaceId) {
    const action = props.closeAction ?? 'close';
    const event = createUIActionEvent(
      props.surfaceId,
      props.componentId,
      action,
      {},
    );
    emit('action', event);
    emitFullAction(event);
  }
}

function handleOverlayClick() {
  if (props.closable) {
    handleClose();
  }
}
</script>

<style scoped>
.nomad-modal-overlay {
  @apply fixed inset-0 flex items-center justify-center;
  z-index: var(--nomad-z-modal-backdrop, 1040);
  background-color: var(--nomad-modal-backdrop, rgba(0, 0, 0, 0.5));
}

.nomad-modal {
  @apply max-w-lg w-full mx-4 max-h-[90vh] overflow-hidden flex flex-col;
  background-color: var(--nomad-surface, #ffffff);
  border-radius: var(--nomad-radius-lg, 0.5rem);
  box-shadow: var(--nomad-shadow-xl, 0 20px 25px -5px rgb(0 0 0 / 0.1));
}

.dark .nomad-modal {
  background-color: var(--nomad-surface, #334155);
}

.nomad-modal-header {
  @apply flex items-center justify-between px-4 py-3;
  border-bottom: 1px solid var(--nomad-border, #e5e7eb);
}

.dark .nomad-modal-header {
  border-bottom-color: var(--nomad-border, #475569);
}

.nomad-modal-title {
  @apply text-lg font-semibold m-0;
  color: var(--nomad-text, #111827);
}

.dark .nomad-modal-title {
  color: var(--nomad-text, #f1f5f9);
}

.nomad-modal-close {
  @apply p-1 rounded transition-colors;
  color: var(--nomad-text-secondary, #6b7280);
}

.nomad-modal-close:hover {
  color: var(--nomad-text, #111827);
}

.dark .nomad-modal-close {
  color: var(--nomad-text-secondary, #cbd5e1);
}

.dark .nomad-modal-close:hover {
  color: var(--nomad-text, #f1f5f9);
}

.nomad-modal-content {
  @apply overflow-y-auto;
  padding: var(--nomad-spacing-md, 1rem);
}

/* Transition styles */
.nomad-modal-enter-active,
.nomad-modal-leave-active {
  transition-duration: var(--nomad-transition-normal, 200ms);
  transition-timing-function: var(--nomad-transition-easing, cubic-bezier(0.4, 0, 0.2, 1));
  transition-property: opacity;
}

.nomad-modal-enter-from,
.nomad-modal-leave-to {
  @apply opacity-0;
}
</style>
