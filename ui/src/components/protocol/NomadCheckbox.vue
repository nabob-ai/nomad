<template>
  <label class="nomad-checkbox" :class="{ 'nomad-checkbox-disabled': disabledValue }">
    <input
      type="checkbox"
      class="nomad-checkbox-input"
      :checked="checkedValue"
      :disabled="disabledValue"
      @change="handleChange"
    />
    <span class="nomad-checkbox-box">
      <svg v-if="checkedValue" class="nomad-checkbox-check" viewBox="0 0 12 12">
        <path d="M3.5 6L5.5 8L8.5 4" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round" />
      </svg>
    </span>
    <span v-if="labelValue" class="nomad-checkbox-label">{{ labelValue }}</span>
  </label>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { DataValue, PropertyValue } from '@/types/ui-protocol';
import { useNomadDataBinding } from '@/composables/useNomadDataBinding';

/**
 * NomadCheckbox Component
 *
 * Checkbox with label and two-way binding support.
 * Uses useNomadDataBinding for reactive data binding with the data model.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  /** Checked property - can be literal or path reference */
  checked?: PropertyValue | boolean;
  /** Label property - can be literal or path reference */
  label?: PropertyValue | string;
  /** Disabled property - can be literal or path reference */
  disabled?: PropertyValue | boolean;
  /** Path for two-way binding (legacy support) */
  checkedPath?: string;
}

const props = withDefaults(defineProps<Props>(), {});

const emit = defineEmits<{
  'update:value': [path: string, value: DataValue];
}>();

// Convert props to PropertyValue format for binding
function toPropertyValue(value: PropertyValue | string | boolean | undefined): PropertyValue | undefined {
  if (value === undefined) return undefined;
  if (typeof value === 'string') return { literalString: value };
  if (typeof value === 'boolean') return { literalBoolean: value };
  return value as PropertyValue;
}

// Use data binding for checked
const { value: boundChecked, updateValue, path: checkedPath } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.checked),
  defaultValue: false,
});

// Use data binding for label
const { value: boundLabel } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.label),
  defaultValue: '',
});

// Use data binding for disabled
const { value: boundDisabled } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.disabled),
  defaultValue: false,
});

// Computed values for template
const checkedValue = computed(() => Boolean(boundChecked.value));
const labelValue = computed(() => boundLabel.value ? String(boundLabel.value) : undefined);
const disabledValue = computed(() => Boolean(boundDisabled.value));

function handleChange(event: Event) {
  const target = event.target as HTMLInputElement;
  const newValue = target.checked;

  // Update via data binding
  updateValue(newValue);

  // Also emit for legacy support
  const path = checkedPath ?? props.checkedPath;
  if (path) {
    emit('update:value', path, newValue);
  }
}
</script>

<style scoped>
.nomad-checkbox {
  @apply inline-flex items-center cursor-pointer;
  gap: var(--nomad-spacing-sm, 0.5rem);
}

.nomad-checkbox-disabled {
  @apply opacity-50 cursor-not-allowed;
}

.nomad-checkbox-input {
  @apply sr-only;
}

.nomad-checkbox-box {
  @apply w-5 h-5 flex items-center justify-center transition-colors;
  border: 2px solid var(--nomad-border, #e5e7eb);
  border-radius: var(--nomad-radius-sm, 0.25rem);
  background-color: var(--nomad-surface, #ffffff);
}

.dark .nomad-checkbox-box {
  border-color: var(--nomad-border, #475569);
  background-color: var(--nomad-surface, #334155);
}

.nomad-checkbox-input:checked + .nomad-checkbox-box {
  background-color: var(--nomad-primary, #3b82f6);
  border-color: var(--nomad-primary, #3b82f6);
}

.nomad-checkbox-input:focus + .nomad-checkbox-box {
  @apply ring-2 ring-offset-2;
  --tw-ring-color: var(--nomad-border-focus, #3b82f6);
}

.nomad-checkbox-check {
  @apply w-3 h-3;
  color: var(--nomad-primary-contrast, #ffffff);
}

.nomad-checkbox-label {
  font-size: var(--nomad-font-size-sm, 0.875rem);
  color: var(--nomad-text-secondary, #6b7280);
}

.dark .nomad-checkbox-label {
  color: var(--nomad-text-secondary, #cbd5e1);
}
</style>
