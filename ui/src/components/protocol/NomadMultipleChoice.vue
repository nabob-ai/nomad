<template>
  <div class="nomad-multiple-choice">
    <label v-if="labelValue" class="nomad-multiple-choice-label">
      {{ labelValue }}
    </label>
    <div class="nomad-multiple-choice-options" :class="directionClass">
      <label
        v-for="option in normalizedOptions"
        :key="option.value"
        class="nomad-multiple-choice-option"
        :class="{ 'nomad-multiple-choice-option-disabled': option.disabled || disabledValue }"
      >
        <input
          type="checkbox"
          class="nomad-multiple-choice-input"
          :value="option.value"
          :checked="isSelected(option.value)"
          :disabled="option.disabled || disabledValue"
          @change="handleChange(option.value, $event)"
        />
        <span class="nomad-multiple-choice-box">
          <svg v-if="isSelected(option.value)" class="nomad-multiple-choice-check" viewBox="0 0 12 12">
            <path d="M3.5 6L5.5 8L8.5 4" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </span>
        <div class="nomad-multiple-choice-content">
          <span class="nomad-multiple-choice-option-label">{{ option.label }}</span>
          <span v-if="option.description" class="nomad-multiple-choice-description">
            {{ option.description }}
          </span>
        </div>
      </label>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { DataValue, MultipleChoiceOption, PropertyValue } from '@/types/ui-protocol';
import { useNomadDataBinding } from '@/composables/useNomadDataBinding';

/**
 * NomadMultipleChoice Component
 *
 * Multiple selection with checkboxes and two-way binding support.
 * Uses useNomadDataBinding for reactive data binding with the data model.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  /** Value property - can be literal or path reference */
  value?: PropertyValue | string[];
  /** Options - can be array or path reference */
  options?: PropertyValue | MultipleChoiceOption[] | string[];
  /** Label property - can be literal or path reference */
  label?: PropertyValue | string;
  /** Disabled property - can be literal or path reference */
  disabled?: PropertyValue | boolean;
  direction?: 'horizontal' | 'vertical';
  /** Path for two-way binding (legacy support) */
  valuePath?: string;
}

const props = withDefaults(defineProps<Props>(), {
  options: () => [],
  direction: 'vertical',
});

const emit = defineEmits<{
  'update:value': [path: string, value: DataValue];
}>();

// Convert props to PropertyValue format for binding
function toPropertyValue(value: PropertyValue | string | boolean | string[] | MultipleChoiceOption[] | undefined): PropertyValue | undefined {
  if (value === undefined) return undefined;
  if (typeof value === 'string') return { literalString: value };
  if (typeof value === 'boolean') return { literalBoolean: value };
  if (Array.isArray(value)) return undefined; // Arrays are handled directly
  return value as PropertyValue;
}

// Use data binding for value
const { value: boundValue, updateValue, path: valuePath } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.value),
  defaultValue: [],
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

// Use data binding for options (if it's a path reference)
const { value: boundOptions } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.options),
  defaultValue: [],
});

// Normalize options to MultipleChoiceOption format
const normalizedOptions = computed((): MultipleChoiceOption[] => {
  // If options is bound from data model
  const optionsSource = boundOptions.value ?? props.options;
  if (!optionsSource || !Array.isArray(optionsSource)) return [];

  return (optionsSource as (MultipleChoiceOption | string)[]).map((opt) => {
    if (typeof opt === 'string') {
      return { value: opt, label: opt };
    }
    return opt as MultipleChoiceOption;
  });
});

// Direction class
const directionClass = computed(() => {
  return props.direction === 'horizontal'
    ? 'nomad-multiple-choice-horizontal'
    : 'nomad-multiple-choice-vertical';
});

// Computed values for template
const currentValue = computed((): string[] => {
  if (Array.isArray(boundValue.value)) {
    return boundValue.value as string[];
  }
  return [];
});
const labelValue = computed(() => boundLabel.value ? String(boundLabel.value) : undefined);
const disabledValue = computed(() => Boolean(boundDisabled.value));

function isSelected(value: string): boolean {
  return currentValue.value.includes(value);
}

function handleChange(value: string, event: Event) {
  const target = event.target as HTMLInputElement;
  let newValue: string[];

  if (target.checked) {
    newValue = [...currentValue.value, value];
  }
  else {
    newValue = currentValue.value.filter(v => v !== value);
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
.nomad-multiple-choice {
  @apply flex flex-col gap-2;
}

.nomad-multiple-choice-label {
  @apply text-sm font-medium text-gray-700 dark:text-gray-300;
}

.nomad-multiple-choice-options {
  @apply flex gap-3;
}

.nomad-multiple-choice-vertical {
  @apply flex-col;
}

.nomad-multiple-choice-horizontal {
  @apply flex-row flex-wrap;
}

.nomad-multiple-choice-option {
  @apply flex items-start gap-2 cursor-pointer;
}

.nomad-multiple-choice-option-disabled {
  @apply opacity-50 cursor-not-allowed;
}

.nomad-multiple-choice-input {
  @apply sr-only;
}

.nomad-multiple-choice-box {
  @apply w-5 h-5 mt-0.5 border-2 border-gray-300 dark:border-gray-600 rounded flex items-center justify-center flex-shrink-0;
  @apply bg-white dark:bg-gray-800 transition-colors;
}

.nomad-multiple-choice-input:checked + .nomad-multiple-choice-box {
  @apply bg-blue-600 border-blue-600;
}

.nomad-multiple-choice-input:focus + .nomad-multiple-choice-box {
  @apply ring-2 ring-blue-500 ring-offset-2;
}

.nomad-multiple-choice-check {
  @apply w-3 h-3 text-white;
}

.nomad-multiple-choice-content {
  @apply flex flex-col;
}

.nomad-multiple-choice-option-label {
  @apply text-sm text-gray-700 dark:text-gray-300;
}

.nomad-multiple-choice-description {
  @apply text-xs text-gray-500 dark:text-gray-400;
}
</style>
