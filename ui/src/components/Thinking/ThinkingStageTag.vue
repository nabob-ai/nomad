<template>
  <span :class="['thinking-stage-tag', stageClass]">
    <component :is="iconComponent" class="stage-icon" />
    {{ displayStage }}
  </span>
</template>

<script setup lang="ts">
import { computed, h, type FunctionalComponent } from 'vue'
import type { ThinkAloudEvent } from '@/types/thinking'

const props = defineProps<{
  event: ThinkAloudEvent
}>()

// 图标组件
const ReasoningIcon: FunctionalComponent = () =>
  h('svg', { class: 'w-3 h-3', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
    h('path', {
      'stroke-linecap': 'round',
      'stroke-linejoin': 'round',
      'stroke-width': '2',
      d: 'M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z'
    })
  ])

const ToolCallIcon: FunctionalComponent = () =>
  h('svg', { class: 'w-3 h-3', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
    h('path', {
      'stroke-linecap': 'round',
      'stroke-linejoin': 'round',
      'stroke-width': '2',
      d: 'M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z'
    })
  ])

const ToolResultIcon: FunctionalComponent = () =>
  h('svg', { class: 'w-3 h-3', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
    h('path', {
      'stroke-linecap': 'round',
      'stroke-linejoin': 'round',
      'stroke-width': '2',
      d: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'
    })
  ])

const ApprovalIcon: FunctionalComponent = () =>
  h('svg', { class: 'w-3 h-3', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
    h('path', {
      'stroke-linecap': 'round',
      'stroke-linejoin': 'round',
      'stroke-width': '2',
      d: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z'
    })
  ])

// 根据事件类型获取图标组件
const iconComponent = computed(() => {
  switch (props.event.type) {
    case 'tool_call':
      return ToolCallIcon
    case 'tool_result':
      return ToolResultIcon
    case 'approval_required':
      return ApprovalIcon
    default:
      return ReasoningIcon
  }
})

// 获取阶段样式类
const stageClass = computed(() => {
  if (props.event.type === 'approval_required' || props.event.stage === '人工介入') {
    return 'stage-approval'
  }
  if (props.event.type === 'tool_call') {
    return 'stage-tool-call'
  }
  if (props.event.type === 'tool_result') {
    return 'stage-tool-result'
  }
  return 'stage-reasoning'
})

// 获取显示的阶段名称
const displayStage = computed(() => {
  if (props.event.stage) return props.event.stage

  switch (props.event.type) {
    case 'tool_call':
      return '工具调用'
    case 'tool_result':
      return '执行结果'
    case 'approval_required':
      return '人工介入'
    default:
      return '任务规划'
  }
})
</script>

<style scoped>
.thinking-stage-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 4px;
  white-space: nowrap;
  flex-shrink: 0;
}

.stage-icon {
  width: 12px;
  height: 12px;
}

.stage-reasoning {
  background: #f1f5f9;
  color: #64748b;
}

.stage-tool-call {
  background: #e0e7ff;
  color: #4338ca;
}

.stage-tool-result {
  background: #d1fae5;
  color: #047857;
}

.stage-approval {
  background: #fef3c7;
  color: #b45309;
}
</style>
