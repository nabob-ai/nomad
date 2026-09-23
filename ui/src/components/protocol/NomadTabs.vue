<template>
  <div class="nomad-tabs">
    <div class="nomad-tabs-header" role="tablist">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        class="nomad-tab-button"
        :class="{ 'nomad-tab-active': activeTabValue === tab.id }"
        :disabled="tab.disabled"
        role="tab"
        :aria-selected="activeTabValue === tab.id"
        @click="handleTabClick(tab.id)"
      >
        <span v-if="tab.icon" class="nomad-tab-icon">{{ tab.icon }}</span>
        <span class="nomad-tab-label">{{ tab.label }}</span>
      </button>
    </div>
    <div class="nomad-tabs-content" role="tabpanel">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { DataValue, UIActionEvent, PropertyValue } from '@/types/ui-protocol';
import { createUIActionEvent } from '@/types/ui-protocol';
import { useNomadDataBinding } from '@/composables/useNomadDataBinding';
import { useUIAction } from '@/composables/useUIAction';

/**
 * NomadTabs Component
 *
 * Tabbed interface with header buttons and content panel.
 * Uses useNomadDataBinding for reactive data binding with the data model.
 * Emits action events both as Vue events and to the Control channel.
 */

interface TabItem {
  id: string;
  label: string;
  icon?: string;
  disabled?: boolean;
}

interface Props {
  componentId?: string;
  surfaceId?: string;
  /** Active tab property - can be literal or path reference */
  activeTab?: PropertyValue | string;
  tabs?: TabItem[];
  /** Path for two-way binding (legacy support) */
  activeTabPath?: string;
}

const props = withDefaults(defineProps<Props>(), {
  tabs: () => [],
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
function toPropertyValue(value: PropertyValue | string | undefined): PropertyValue | undefined {
  if (value === undefined) return undefined;
  if (typeof value === 'string') return { literalString: value };
  return value as PropertyValue;
}

// Use data binding for activeTab
const { value: boundActiveTab, updateValue, path: activeTabPath } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.activeTab),
  defaultValue: '',
});

// Computed value for template
const activeTabValue = computed(() => String(boundActiveTab.value ?? ''));

function handleTabClick(tabId: string) {
  // Update via data binding
  updateValue(tabId);

  // Also emit for legacy support
  const path = activeTabPath ?? props.activeTabPath;
  if (path) {
    emit('update:value', path, tabId);
  }

  // Emit action event
  if (props.componentId && props.surfaceId) {
    const event = createUIActionEvent(
      props.surfaceId,
      props.componentId,
      'tab-change',
      { tabId },
    );
    emit('action', event);
    emitFullAction(event);
  }
}
</script>

<style scoped>
.nomad-tabs {
  @apply flex flex-col;
}

.nomad-tabs-header {
  @apply flex border-b border-gray-200 dark:border-gray-700;
}

.nomad-tab-button {
  @apply px-4 py-2 text-sm font-medium text-gray-600 dark:text-gray-400 border-b-2 border-transparent transition-colors;
  @apply hover:text-gray-900 dark:hover:text-gray-100;
  @apply disabled:opacity-50 disabled:cursor-not-allowed;
}

.nomad-tab-active {
  @apply text-blue-600 dark:text-blue-400 border-blue-600 dark:border-blue-400;
}

.nomad-tab-icon {
  @apply mr-2;
}

.nomad-tabs-content {
  @apply p-4;
}
</style>
