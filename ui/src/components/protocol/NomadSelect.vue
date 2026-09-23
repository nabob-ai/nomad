<template>
  <div class="nomad-select">
    <label v-if="labelValue" :for="selectId" class="nomad-select-label">
      {{ labelValue }}
    </label>
    <select
      :id="selectId"
      :value="currentValue"
      class="nomad-select-input"
      :disabled="disabledValue"
      :multiple="multiple"
      @change="handleChange"
    >
      <option v-if="placeholderValue && !multiple" value="" disabled>
        {{ placeholderValue }}
      </option>
      <option
        v-for="option in normalizedOptions"
        :key="option.value"
        :value="option.value"
        :disabled="option.disabled"
      >
        {{ option.label }}
      </option>
    </select>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { DataValue, SelectOption, PropertyValue } from '@/types/ui-protocol';
import { useNomadDataBinding } from '@/composables/useNomadDataBinding';

/**
 * NomadSelect Component
 *
 * Dropdown select with options and two-way binding support.
 * Uses useNomadDataBinding for reactive data binding with the data model.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  /** Value property - can be literal or path reference */
  value?: PropertyValue | string | string[];
  /** Options - can be array or path reference */
  options?: PropertyValue | SelectOption[] | string[];
  /** Label property - can be literal or path reference */
  label?: PropertyValue | string;
  /** Placeholder property - can be literal or path reference */
  placeholder?: PropertyValue | string;
  /** Disabled property - can be literal or path reference */
  disabled?: PropertyValue | boolean;
  multiple?: boolean;
  /** Path for two-way binding (legacy support) */
  valuePath?: string;
}

const props = withDefaults(defineProps<Props>(), {
  options: () => [],
  multiple: false,
});

const emit = defineEmits<{
  'update:value': [path: string, value: DataValue];
}>();

// Generate unique ID for label association
const selectId = computed(() => `nomad-select-${props.componentId ?? Math.random().toString(36).slice(2)}`);

// Convert props to PropertyValue format for binding
function toPropertyValue(value: PropertyValue | string | boolean | string[] | SelectOption[] | undefined): PropertyValue | undefined {
  if (value === undefined) return undefined;
  if (typeof value === 'string') return { literalString: value };
  if (typeof value === 'boolean') return { literalBoolean: value };
  if (Array.isArray(value)) return undefined; // Arrays are handled directly
  return value as PropertyValue;
}

// Use data binding for value
const { value: boundValue, updateValue, path: valuePath } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.value),
  defaultValue: props.multiple ? [] : '',
});

// Use data binding for label
const { value: boundLabel } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.label),
  defaultValue: '',
});

// Use data binding for placeholder
const { value: boundPlaceholder } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.placeholder),
  defaultValue: '',
});

// Use data binding for disabled
const { value: boundDisabled } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.disabled),
  defaultValue: false,
});

// Use data binding for options (if it's a path reference)
const { value: boundOptions } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.options),
  defaultValue: [],
});

// Normalize options to SelectOption format
const normalizedOptions = computed((): SelectOption[] => {
  // If options is bound from data model
  const optionsSource = boundOptions.value ?? props.options;
  if (!optionsSource || !Array.isArray(optionsSource)) return [];

  return (optionsSource as (SelectOption | string)[]).map((opt) => {
    if (typeof opt === 'string') {
      return { value: opt, label: opt };
    }
    return opt as SelectOption;
  });
});

// Computed values for template
const currentValue = computed(() => {
  if (props.multiple) {
    return Array.isArray(boundValue.value) ? boundValue.value : [];
  }
  return String(boundValue.value ?? '');
});
const labelValue = computed(() => boundLabel.value ? String(boundLabel.value) : undefined);
const placeholderValue = computed(() => String(boundPlaceholder.value ?? ''));
const disabledValue = computed(() => Boolean(boundDisabled.value));

function handleChange(event: Event) {
  const target = event.target as HTMLSelectElement;
  let newValue: DataValue;

  if (props.multiple) {
    newValue = Array.from(target.selectedOptions).map(opt => opt.value);
  }
  else {
    newValue = target.value;
  }

  // Update via data binding
  updateValue(newValue);

  // Also emit for legacy support
  const path = valuePath ?? props.valuePath;
  if (path) {
    emit('update:value', path, newValue);
  }
}
</script>

<style scoped>
.nomad-select {
  @apply flex flex-col;
  gap: var(--nomad-spacing-xs, 0.25rem);
}

.nomad-select-label {
  @apply font-medium;
  font-size: var(--nomad-font-size-sm, 0.875rem);
  color: var(--nomad-text-secondary, #6b7280);
}

.dark .nomad-select-label {
  color: var(--nomad-text-secondary, #cbd5e1);
}

.nomad-select-input {
  @apply w-full focus:outline-none focus:ring-2 focus:border-transparent;
  padding: var(--nomad-spacing-sm, 0.5rem) var(--nomad-spacing-md, 1rem);
  font-size: var(--nomad-font-size-sm, 0.875rem);
  border: 1px solid var(--nomad-border, #e5e7eb);
  border-radius: var(--nomad-radius-lg, 0.5rem);
  background-color: var(--nomad-surface, #ffffff);
  color: var(--nomad-text, #111827);
  --tw-ring-color: var(--nomad-border-focus, #3b82f6);
}

.dark .nomad-select-input {
  border-color: var(--nomad-border, #475569);
  background-color: var(--nomad-surface, #334155);
  color: var(--nomad-text, #f1f5f9);
}

.nomad-select-input:disabled {
  @apply cursor-not-allowed opacity-50;
  background-color: var(--nomad-background-alt, #f3f4f6);
}

.dark .nomad-select-input:disabled {
  background-color: var(--nomad-background-alt, #0f172a);
}
</style>
