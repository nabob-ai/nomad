<template>
  <div class="nomad-datetime-input">
    <label v-if="labelValue" :for="inputId" class="nomad-datetime-label">
      {{ labelValue }}
    </label>
    <input
      :id="inputId"
      :value="currentValue"
      class="nomad-datetime-field"
      :type="inputType"
      :disabled="disabledValue"
      :min="min"
      :max="max"
      @change="handleChange"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { DataValue, DateTimeInputType, PropertyValue } from '@/types/ui-protocol';
import { useNomadDataBinding } from '@/composables/useNomadDataBinding';

/**
 * NomadDateTimeInput Component
 *
 * Date/time input with type variants and two-way binding support.
 * Uses useNomadDataBinding for reactive data binding with the data model.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  /** Value property - can be literal or path reference */
  value?: PropertyValue | string;
  /** Label property - can be literal or path reference */
  label?: PropertyValue | string;
  type?: DateTimeInputType;
  /** Disabled property - can be literal or path reference */
  disabled?: PropertyValue | boolean;
  min?: string;
  max?: string;
  /** Path for two-way binding (legacy support) */
  valuePath?: string;
}

const props = withDefaults(defineProps<Props>(), {
  type: 'date',
});

const emit = defineEmits<{
  'update:value': [path: string, value: DataValue];
}>();

// Generate unique ID for label association
const inputId = computed(() => `nomad-datetime-${props.componentId ?? Math.random().toString(36).slice(2)}`);

// Map type to HTML input type
const inputType = computed(() => {
  const typeMap: Record<DateTimeInputType, string> = {
    date: 'date',
    time: 'time',
    datetime: 'datetime-local',
  };
  return typeMap[props.type] ?? 'date';
});

// Convert props to PropertyValue format for binding
function toPropertyValue(value: PropertyValue | string | boolean | undefined): PropertyValue | undefined {
  if (value === undefined) return undefined;
  if (typeof value === 'string') return { literalString: value };
  if (typeof value === 'boolean') return { literalBoolean: value };
  return value as PropertyValue;
}

// Use data binding for value
const { value: boundValue, updateValue, path: valuePath } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.value),
  defaultValue: '',
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
const currentValue = computed(() => String(boundValue.value ?? ''));
const labelValue = computed(() => boundLabel.value ? String(boundLabel.value) : undefined);
const disabledValue = computed(() => Boolean(boundDisabled.value));

function handleChange(event: Event) {
  const target = event.target as HTMLInputElement;
  const newValue = target.value;

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
.nomad-datetime-input {
  @apply flex flex-col gap-1;
}

.nomad-datetime-label {
  @apply text-sm font-medium text-gray-700 dark:text-gray-300;
}

.nomad-datetime-field {
  @apply w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg;
  @apply bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100;
  @apply focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent;
  @apply disabled:opacity-50 disabled:cursor-not-allowed disabled:bg-gray-100 dark:disabled:bg-gray-700;
}
</style>
