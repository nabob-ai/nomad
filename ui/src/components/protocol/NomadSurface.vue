<template>
  <div
    ref="containerRef"
    class="nomad-surface"
    :style="surfaceStyles"
    :data-surface-id="surfaceId"
    :data-deleted="isDeleted"
  >
    <template v-if="isDeleted">
      <div class="nomad-surface-deleted">
        <slot name="deleted">
          <!-- Surface deleted, component will be unmounted -->
        </slot>
      </div>
    </template>
    <template v-else-if="surface?.componentTree">
      <NomadComponentRenderer
        :node="surface.componentTree"
        :surface-id="surfaceId"
        @action="handleAction"
        @update:value="handleValueUpdate"
      />
    </template>
    <div v-else-if="loading || isStreamingMode" class="nomad-surface-loading">
      <slot name="loading">
        <div class="nomad-surface-spinner" />
      </slot>
    </div>
    <div v-else class="nomad-surface-empty">
      <slot name="empty">
        <span class="nomad-surface-empty-text">No content</span>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, provide, type Ref } from 'vue';
import type { Surface, UIActionEvent, DataValue, DataMap } from '@/types/ui-protocol';
import { MessageProcessor } from '@/protocol/message-processor';
import { DATA_BINDING_CONTEXT_KEY, type DataBindingContext } from '@/composables/useNomadDataBinding';
import { UI_ACTION_CONTEXT_KEY, createUIActionContext, type UIActionContext } from '@/composables/useUIAction';
import { useStreamingState } from '@/composables/useStreamingState';
import { convertProtocolStyles, useNomadTheme } from '@/composables/useNomadTheme';
import NomadComponentRenderer from './NomadComponentRenderer.vue';

/**
 * NomadSurface Component
 *
 * Renders a UI surface based on the Nomad UI Protocol.
 * Integrates with MessageProcessor to manage surface state.
 * Provides data binding context to child components.
 * Supports streaming rendering with state preservation.
 */

interface Props {
  /** Surface ID to render */
  surfaceId: string;
  /** Optional MessageProcessor instance (creates one if not provided) */
  processor?: MessageProcessor;
  /** Show loading state */
  loading?: boolean;
  /** Whether to preserve scroll position during updates */
  preserveScroll?: boolean;
  /** Whether to preserve input focus during updates */
  preserveFocus?: boolean;
  /** Optional callback for sending UI actions to the Control channel */
  onControlAction?: (event: UIActionEvent) => void;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  preserveScroll: true,
  preserveFocus: true,
});

const emit = defineEmits<{
  /** Emitted when a user interacts with a component */
  action: [event: UIActionEvent];
  /** Emitted when the surface is updated */
  'surface-update': [surface: Surface];
  /** Emitted when the surface is deleted */
  'surface-delete': [surfaceId: string];
  /** Emitted when a UI action should be sent to the Control channel */
  'control-action': [event: UIActionEvent];
}>();

// Container ref for state preservation
const containerRef = ref<HTMLElement | null>(null);

// Internal processor instance
const internalProcessor = ref<MessageProcessor | null>(null);

// Get the processor to use
const processorRef = computed(() => props.processor ?? internalProcessor.value);

// Surface state
const surface = ref<Surface | undefined>(undefined);

// Reactive data model for binding - use a ref that we update manually
const dataModel = ref<DataMap>({});

// Check if in streaming mode
const isStreamingMode = computed(() => {
  if (!processorRef.value) return false;
  return processorRef.value.isStreaming(props.surfaceId);
});

// Unsubscribe function
let unsubscribe: (() => void) | null = null;

// Unsubscribe from delete events
let unsubscribeDelete: (() => void) | null = null;

// Track if surface has been deleted
const isDeleted = ref(false);

// Setup theme system
const { isDark, applyTheme, removeTheme } = useNomadTheme();

// Provide theme context to child components
provide('nomad-theme-dark', isDark);

// Track applied styles for cleanup
let appliedStyles: Record<string, string> = {};

// Computed surface styles from CSS custom properties
// Only allows CSS custom properties (--*) for security
const surfaceStyles = computed(() => {
  if (!surface.value?.styles) {
    return {};
  }
  // Convert and validate protocol styles - only CSS custom properties allowed
  return convertProtocolStyles(surface.value.styles);
});

// Setup streaming state preservation
const {
  saveState,
  restoreState,
} = useStreamingState({
  surfaceId: props.surfaceId,
  containerRef,
  preserveScroll: props.preserveScroll,
  preserveFocus: props.preserveFocus,
});

// Provide processor and surfaceId to child components
provide('nomad-processor', processorRef);
provide('nomad-surface-id', props.surfaceId);

// Create stable data binding context (not computed, to maintain stable refs)
const dataBindingContext: DataBindingContext = {
  processor: processorRef as Ref<MessageProcessor | null>,
  surfaceId: props.surfaceId,
  dataModel,
};
provide(DATA_BINDING_CONTEXT_KEY, dataBindingContext);

// Connection status for UI action context
const isActionConnected = ref(true);

/**
 * Emit UI action to the Control channel
 * This function is called by child components via the UI action context
 */
function emitControlAction(event: UIActionEvent): void {
  // Emit Vue event for parent components
  emit('action', event);
  emit('control-action', event);

  // Call the optional callback for sending to Control channel
  if (props.onControlAction) {
    props.onControlAction(event);
  }
}

// Create and provide UI action context for child components
const uiActionContext: UIActionContext = createUIActionContext(
  emitControlAction,
  isActionConnected,
);
provide(UI_ACTION_CONTEXT_KEY, uiActionContext);

// Watch surface changes to update dataModel ref
watch(
  () => surface.value?.dataModel,
  (newDataModel) => {
    if (newDataModel) {
      dataModel.value = newDataModel;
    }
  },
  { immediate: true, deep: true },
);

// Watch surface styles and apply them to the container element
watch(
  surfaceStyles,
  (newStyles, oldStyles) => {
    if (!containerRef.value) return;

    // Remove old styles
    if (oldStyles && Object.keys(oldStyles).length > 0) {
      removeTheme(containerRef.value, oldStyles);
    }

    // Apply new styles
    if (newStyles && Object.keys(newStyles).length > 0) {
      applyTheme(containerRef.value, newStyles);
      appliedStyles = { ...newStyles };
    }
    else {
      appliedStyles = {};
    }
  },
  { immediate: true, deep: true },
);

/**
 * Handle action events from child components
 * Emits the action to both Vue events and the Control channel
 */
function handleAction(event: UIActionEvent) {
  emitControlAction(event);
}

/**
 * Handle value updates from input components (two-way binding)
 */
function handleValueUpdate(path: string, value: DataValue) {
  if (processorRef.value) {
    processorRef.value.setData(props.surfaceId, path, value);
  }
}

/**
 * Handle surface deletion
 * Cleans up resources and emits the surface-delete event
 */
function handleSurfaceDelete(surfaceId: string) {
  // Mark as deleted
  isDeleted.value = true;

  // Clear surface state
  surface.value = undefined;
  dataModel.value = {};

  // Emit deletion event
  emit('surface-delete', surfaceId);

  // Cleanup subscriptions
  cleanupSubscriptions();
}

/**
 * Cleanup all subscriptions
 */
function cleanupSubscriptions() {
  if (unsubscribe) {
    unsubscribe();
    unsubscribe = null;
  }
  if (unsubscribeDelete) {
    unsubscribeDelete();
    unsubscribeDelete = null;
  }
}

/**
 * Subscribe to surface changes with state preservation
 */
function subscribeSurface() {
  if (!processorRef.value) {
    return;
  }

  // Cleanup previous subscriptions
  cleanupSubscriptions();

  // Reset deleted state
  isDeleted.value = false;

  // Get initial surface state
  surface.value = processorRef.value.getSurface(props.surfaceId);

  // Subscribe to changes with state preservation
  unsubscribe = processorRef.value.subscribe(props.surfaceId, (updatedSurface) => {
    // Save state before update
    saveState();

    // Update surface
    surface.value = updatedSurface;
    emit('surface-update', updatedSurface);

    // Restore state after update
    restoreState();
  });

  // Subscribe to deletion events
  unsubscribeDelete = processorRef.value.subscribeToDelete(props.surfaceId, handleSurfaceDelete);
}

// Watch for processor changes
watch(
  () => props.processor,
  () => {
    subscribeSurface();
  },
);

// Watch for surfaceId changes
watch(
  () => props.surfaceId,
  () => {
    subscribeSurface();
  },
);

onMounted(() => {
  // Create internal processor if not provided
  if (!props.processor) {
    internalProcessor.value = new MessageProcessor();
  }
  subscribeSurface();
});

onUnmounted(() => {
  // Cleanup all subscriptions
  cleanupSubscriptions();

  // Remove applied theme styles
  if (containerRef.value && Object.keys(appliedStyles).length > 0) {
    removeTheme(containerRef.value, appliedStyles);
  }
});

// Expose processor for external access
defineExpose({
  processor: processorRef,
  surface,
  dataModel,
  isStreamingMode,
  isDeleted,
  isDark,
  surfaceStyles,
  saveState,
  restoreState,
  emitControlAction,
  isActionConnected,
  cleanupSubscriptions,
});
</script>

<style scoped>
.nomad-surface {
  @apply relative w-full;
  /* Apply theme font family */
  font-family: var(--nomad-font-family, inherit);
  color: var(--nomad-text, inherit);
}

.nomad-surface-loading {
  @apply flex items-center justify-center p-4;
}

.nomad-surface-spinner {
  @apply w-6 h-6 rounded-full animate-spin;
  border-width: 2px;
  border-style: solid;
  border-color: var(--nomad-border, #e5e7eb);
  border-top-color: var(--nomad-primary, #3b82f6);
}

.nomad-surface-empty {
  @apply flex items-center justify-center p-4;
  color: var(--nomad-text-muted, #9ca3af);
}

.nomad-surface-empty-text {
  @apply text-sm;
}

.nomad-surface-deleted {
  /* Hidden by default, parent can use slot to show custom content */
  @apply hidden;
}

.nomad-surface[data-deleted="true"] {
  /* Surface is marked as deleted */
  @apply pointer-events-none;
}
</style>
