/**
 * Nomad UI Protocol - 标准组件白名单
 *
 * 定义和注册所有标准组件类型，确保只有受信任的 UI 组件才能被渲染。
 *
 * @module protocol/standard-components
 */

import { ComponentRegistry, createComponentRegistry, type ComponentConstructor } from './registry';

/**
 * 标准布局组件类型
 * 需求: 9.1
 */
export const LAYOUT_COMPONENTS = [
  'Row',
  'Column',
  'Card',
  'List',
  'Tabs',
  'Modal',
  'Divider',
] as const;

/**
 * 标准内容组件类型
 * 需求: 9.2
 */
export const CONTENT_COMPONENTS = [
  'Text',
  'Image',
  'Icon',
  'Video',
  'AudioPlayer',
] as const;

/**
 * 标准输入组件类型
 * 需求: 9.3
 */
export const INPUT_COMPONENTS = [
  'Button',
  'TextField',
  'Checkbox',
  'Select',
  'DateTimeInput',
  'Slider',
  'MultipleChoice',
] as const;

/**
 * 自定义组件类型
 * 需求: 2.3
 */
export const CUSTOM_COMPONENT = 'Custom' as const;

/**
 * 所有标准组件类型
 * 需求: 2.1, 2.3
 */
export const STANDARD_COMPONENTS = [
  ...LAYOUT_COMPONENTS,
  ...CONTENT_COMPONENTS,
  ...INPUT_COMPONENTS,
  CUSTOM_COMPONENT,
] as const;

/**
 * 标准组件类型
 */
export type StandardComponentType = typeof STANDARD_COMPONENTS[number];

/**
 * 布局组件类型
 */
export type LayoutComponentType = typeof LAYOUT_COMPONENTS[number];

/**
 * 内容组件类型
 */
export type ContentComponentType = typeof CONTENT_COMPONENTS[number];

/**
 * 输入组件类型
 */
export type InputComponentType = typeof INPUT_COMPONENTS[number];

/**
 * 判断是否为标准组件类型
 *
 * @param type - 组件类型名称
 * @returns 是否为标准组件类型
 */
export function isStandardComponent(type: string): type is StandardComponentType {
  return (STANDARD_COMPONENTS as readonly string[]).includes(type);
}

/**
 * 判断是否为布局组件类型
 *
 * @param type - 组件类型名称
 * @returns 是否为布局组件类型
 */
export function isLayoutComponent(type: string): type is LayoutComponentType {
  return (LAYOUT_COMPONENTS as readonly string[]).includes(type);
}

/**
 * 判断是否为内容组件类型
 *
 * @param type - 组件类型名称
 * @returns 是否为内容组件类型
 */
export function isContentComponent(type: string): type is ContentComponentType {
  return (CONTENT_COMPONENTS as readonly string[]).includes(type);
}

/**
 * 判断是否为输入组件类型
 *
 * @param type - 组件类型名称
 * @returns 是否为输入组件类型
 */
export function isInputComponent(type: string): type is InputComponentType {
  return (INPUT_COMPONENTS as readonly string[]).includes(type);
}

/**
 * 创建占位符组件构造函数
 *
 * 在实际的 Vue 渲染器实现之前，使用占位符构造函数注册标准组件。
 * 这些占位符将在 Vue 组件实现后被替换。
 *
 * @param typeName - 组件类型名称
 * @returns 占位符组件构造函数
 */
function createPlaceholderConstructor(typeName: string): ComponentConstructor {
  // 创建一个简单的 HTMLElement 子类作为占位符
  const PlaceholderComponent = class extends HTMLElement {
    static componentType = typeName;

    connectedCallback() {
      this.textContent = `[${typeName} placeholder]`;
    }
  };

  // 设置类名以便调试
  Object.defineProperty(PlaceholderComponent, 'name', { value: `${typeName}Placeholder` });

  return PlaceholderComponent as unknown as ComponentConstructor;
}

/**
 * 注册所有标准组件到注册表
 *
 * @param registry - 组件注册表
 * @param constructorMap - 可选的组件构造函数映射，用于提供实际的组件实现
 */
export function registerStandardComponents(
  registry: ComponentRegistry,
  constructorMap?: Partial<Record<StandardComponentType, ComponentConstructor>>,
): void {
  for (const componentType of STANDARD_COMPONENTS) {
    const constructor = constructorMap?.[componentType] ?? createPlaceholderConstructor(componentType);
    registry.register(componentType, constructor);
  }
}

/**
 * 创建并初始化标准组件注册表
 *
 * @param constructorMap - 可选的组件构造函数映射
 * @param freeze - 是否在初始化后冻结注册表（生产模式）
 * @returns 初始化后的组件注册表
 */
export function createStandardRegistry(
  constructorMap?: Partial<Record<StandardComponentType, ComponentConstructor>>,
  freeze: boolean = false,
): ComponentRegistry {
  const registry = createComponentRegistry();
  registerStandardComponents(registry, constructorMap);

  if (freeze) {
    registry.freeze();
  }

  return registry;
}

/**
 * 默认的标准组件注册表实例
 *
 * 注意：在生产环境中，应该使用 createStandardRegistry() 创建新实例，
 * 并在初始化完成后调用 freeze() 方法。
 */
let defaultRegistry: ComponentRegistry | null = null;

/**
 * 获取默认的标准组件注册表
 *
 * @returns 默认的组件注册表实例
 */
export function getDefaultRegistry(): ComponentRegistry {
  if (!defaultRegistry) {
    defaultRegistry = createStandardRegistry();
  }
  return defaultRegistry;
}

/**
 * 重置默认注册表（仅用于测试）
 */
export function resetDefaultRegistry(): void {
  defaultRegistry = null;
}
