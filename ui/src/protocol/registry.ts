/**
 * Nomad UI Protocol - 组件注册表
 *
 * 管理可用组件类型的白名单，实现安全的组件注册和查找机制。
 *
 * @module protocol/registry
 */

/**
 * 组件构造函数类型
 * 用于自定义组件的注册
 */
export type ComponentConstructor = new () => HTMLElement;

/**
 * 注册表错误代码
 */
export const RegistryErrorCodes = {
  REGISTRY_FROZEN: 'REGISTRY_FROZEN',
  INVALID_TYPE_NAME: 'INVALID_TYPE_NAME',
  COMPONENT_ALREADY_REGISTERED: 'COMPONENT_ALREADY_REGISTERED',
} as const;

export type RegistryErrorCode = typeof RegistryErrorCodes[keyof typeof RegistryErrorCodes];

/**
 * 注册表错误
 */
export class RegistryError extends Error {
  constructor(
    message: string,
    public code: RegistryErrorCode,
    public details?: Record<string, unknown>,
  ) {
    super(message);
    this.name = 'RegistryError';
  }
}

/**
 * 验证组件类型名称是否有效
 * 类型名称必须以字母开头，只能包含字母和数字
 *
 * @param typeName - 要验证的类型名称
 * @returns 是否有效
 */
export function isValidTypeName(typeName: string): boolean {
  return /^[a-zA-Z][a-zA-Z0-9]*$/.test(typeName);
}

/**
 * 组件注册表
 *
 * 管理可用组件类型的白名单，提供组件注册、查找和冻结功能。
 * 在生产模式下，注册表可以被冻结以防止运行时修改。
 *
 * @example
 * ```typescript
 * const registry = new ComponentRegistry();
 * registry.register('MyComponent', MyComponentClass);
 * registry.freeze(); // 生产模式下冻结
 *
 * if (registry.has('MyComponent')) {
 *   const constructor = registry.get('MyComponent');
 * }
 * ```
 */
export class ComponentRegistry {
  /** 组件映射表 */
  private components: Map<string, ComponentConstructor> = new Map();

  /** 是否已冻结 */
  private frozen: boolean = false;

  /**
   * 注册组件
   *
   * @param typeName - 组件类型名称（必须以字母开头，只能包含字母和数字）
   * @param constructor - 组件构造函数
   * @param tagName - 可选的自定义元素标签名（用于 Web Components）
   * @throws {RegistryError} 如果注册表已冻结、类型名称无效或组件已注册
   */
  register(typeName: string, constructor: ComponentConstructor, _tagName?: string): void {
    // 检查是否已冻结
    if (this.frozen) {
      throw new RegistryError(
        'Cannot register component: registry is frozen in production mode',
        RegistryErrorCodes.REGISTRY_FROZEN,
        { typeName },
      );
    }

    // 验证类型名称
    if (!isValidTypeName(typeName)) {
      throw new RegistryError(
        `Invalid component type name: "${typeName}". Type name must start with a letter and contain only alphanumeric characters.`,
        RegistryErrorCodes.INVALID_TYPE_NAME,
        { typeName },
      );
    }

    // 检查是否已注册
    if (this.components.has(typeName)) {
      // 如果已注册相同的构造函数，静默返回（幂等性）
      if (this.components.get(typeName) === constructor) {
        return;
      }
      // 否则记录警告并返回（不覆盖）
      console.warn(`Component "${typeName}" is already registered. Skipping registration.`);
      return;
    }

    // 注册组件
    this.components.set(typeName, constructor);
  }

  /**
   * 获取已注册的组件构造函数
   *
   * @param typeName - 组件类型名称
   * @returns 组件构造函数，如果未注册则返回 undefined
   */
  get(typeName: string): ComponentConstructor | undefined {
    return this.components.get(typeName);
  }

  /**
   * 检查组件是否已注册
   *
   * @param typeName - 组件类型名称
   * @returns 是否已注册
   */
  has(typeName: string): boolean {
    return this.components.has(typeName);
  }

  /**
   * 冻结注册表
   *
   * 冻结后，任何注册新组件的尝试都会抛出错误。
   * 用于生产模式下防止运行时修改。
   */
  freeze(): void {
    this.frozen = true;
  }

  /**
   * 检查注册表是否已冻结
   *
   * @returns 是否已冻结
   */
  isFrozen(): boolean {
    return this.frozen;
  }

  /**
   * 获取所有已注册的组件类型名称
   *
   * @returns 组件类型名称数组
   */
  getRegisteredTypes(): string[] {
    return Array.from(this.components.keys());
  }

  /**
   * 获取已注册组件的数量
   *
   * @returns 组件数量
   */
  size(): number {
    return this.components.size;
  }
}

/**
 * 创建新的组件注册表实例
 *
 * @returns 新的 ComponentRegistry 实例
 */
export function createComponentRegistry(): ComponentRegistry {
  return new ComponentRegistry();
}
