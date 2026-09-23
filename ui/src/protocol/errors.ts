/**
 * Nomad UI Protocol - 错误类型和错误码
 *
 * 定义协议级别的错误类型和错误码，用于统一的错误处理。
 *
 * @module protocol/errors
 */

/**
 * 协议错误代码
 *
 * 定义所有可能的协议错误类型，用于错误分类和处理。
 */
export const ErrorCodes = {
  /** 无效的消息格式 */
  INVALID_MESSAGE: 'INVALID_MESSAGE',
  /** 未知的组件类型（不在白名单中） */
  UNKNOWN_COMPONENT: 'UNKNOWN_COMPONENT',
  /** 无效的 JSON Pointer 路径 */
  INVALID_PATH: 'INVALID_PATH',
  /** 组件树中检测到循环引用 */
  CIRCULAR_REFERENCE: 'CIRCULAR_REFERENCE',
  /** 注册表已冻结，无法注册新组件 */
  REGISTRY_FROZEN: 'REGISTRY_FROZEN',
  /** 无效的组件类型名称 */
  INVALID_TYPE_NAME: 'INVALID_TYPE_NAME',
  /** 检测到 XSS 攻击代码 */
  XSS_DETECTED: 'XSS_DETECTED',
  /** 无效或不安全的 URL */
  INVALID_URL: 'INVALID_URL',
  /** 无效的组件属性 */
  INVALID_PROPS: 'INVALID_PROPS',
  /** 无效的 Surface ID */
  INVALID_SURFACE_ID: 'INVALID_SURFACE_ID',
  /** 子组件引用不存在 */
  INVALID_CHILD_REFERENCE: 'INVALID_CHILD_REFERENCE',
} as const;

export type ErrorCode = typeof ErrorCodes[keyof typeof ErrorCodes];

/**
 * 协议错误
 *
 * 用于表示 Nomad UI Protocol 中的各种错误情况。
 * 包含错误代码和可选的详细信息，便于错误分类和调试。
 *
 * @example
 * ```typescript
 * throw new ProtocolError(
 *   'Unknown component type: InvalidComponent',
 *   ErrorCodes.UNKNOWN_COMPONENT,
 *   { typeName: 'InvalidComponent', surfaceId: 'surface-1' }
 * );
 * ```
 */
export class ProtocolError extends Error {
  /**
   * 创建协议错误实例
   *
   * @param message - 错误消息
   * @param code - 错误代码
   * @param details - 可选的错误详情
   */
  constructor(
    message: string,
    public code: ErrorCode,
    public details?: Record<string, unknown>,
  ) {
    super(message);
    this.name = 'ProtocolError';

    // 确保原型链正确（TypeScript 编译到 ES5 时需要）
    Object.setPrototypeOf(this, ProtocolError.prototype);
  }

  /**
   * 转换为 JSON 格式
   *
   * @returns JSON 表示
   */
  toJSON(): Record<string, unknown> {
    return {
      name: this.name,
      message: this.message,
      code: this.code,
      details: this.details,
    };
  }

  /**
   * 转换为字符串格式
   *
   * @returns 字符串表示
   */
  toString(): string {
    const detailsStr = this.details ? ` Details: ${JSON.stringify(this.details)}` : '';
    return `${this.name} [${this.code}]: ${this.message}${detailsStr}`;
  }
}

/**
 * 检查错误是否为协议错误
 *
 * @param error - 要检查的错误
 * @returns 是否为 ProtocolError 实例
 */
export function isProtocolError(error: unknown): error is ProtocolError {
  return error instanceof ProtocolError;
}

/**
 * 检查错误是否为特定类型的协议错误
 *
 * @param error - 要检查的错误
 * @param code - 期望的错误代码
 * @returns 是否为指定类型的 ProtocolError
 */
export function isProtocolErrorWithCode(error: unknown, code: ErrorCode): boolean {
  return isProtocolError(error) && error.code === code;
}

/**
 * 创建无效消息错误
 *
 * @param message - 错误消息
 * @param details - 错误详情
 * @returns ProtocolError 实例
 */
export function createInvalidMessageError(
  message: string,
  details?: Record<string, unknown>,
): ProtocolError {
  return new ProtocolError(message, ErrorCodes.INVALID_MESSAGE, details);
}

/**
 * 创建未知组件错误
 *
 * @param typeName - 未知的组件类型名称
 * @param details - 额外的错误详情
 * @returns ProtocolError 实例
 */
export function createUnknownComponentError(
  typeName: string,
  details?: Record<string, unknown>,
): ProtocolError {
  return new ProtocolError(
    `Unknown component type: ${typeName}`,
    ErrorCodes.UNKNOWN_COMPONENT,
    { typeName, ...details },
  );
}

/**
 * 创建无效路径错误
 *
 * @param path - 无效的路径
 * @param details - 额外的错误详情
 * @returns ProtocolError 实例
 */
export function createInvalidPathError(
  path: string,
  details?: Record<string, unknown>,
): ProtocolError {
  return new ProtocolError(
    `Invalid JSON Pointer path: ${path}`,
    ErrorCodes.INVALID_PATH,
    { path, ...details },
  );
}

/**
 * 创建循环引用错误
 *
 * @param componentId - 导致循环引用的组件 ID
 * @param details - 额外的错误详情
 * @returns ProtocolError 实例
 */
export function createCircularReferenceError(
  componentId: string,
  details?: Record<string, unknown>,
): ProtocolError {
  return new ProtocolError(
    `Circular reference detected for component: ${componentId}`,
    ErrorCodes.CIRCULAR_REFERENCE,
    { componentId, ...details },
  );
}

/**
 * 创建无效属性错误
 *
 * @param componentId - 组件 ID
 * @param propName - 无效的属性名称
 * @param details - 额外的错误详情
 * @returns ProtocolError 实例
 */
export function createInvalidPropsError(
  componentId: string,
  propName?: string,
  details?: Record<string, unknown>,
): ProtocolError {
  const propInfo = propName ? ` (property: ${propName})` : '';
  return new ProtocolError(
    `Invalid or malformed props for component: ${componentId}${propInfo}`,
    ErrorCodes.INVALID_PROPS,
    { componentId, propName, ...details },
  );
}

/**
 * 创建无效子组件引用错误
 *
 * @param parentId - 父组件 ID
 * @param childId - 不存在的子组件 ID
 * @param details - 额外的错误详情
 * @returns ProtocolError 实例
 */
export function createInvalidChildReferenceError(
  parentId: string,
  childId: string,
  details?: Record<string, unknown>,
): ProtocolError {
  return new ProtocolError(
    `Child component "${childId}" not found for parent "${parentId}"`,
    ErrorCodes.INVALID_CHILD_REFERENCE,
    { parentId, childId, ...details },
  );
}

/**
 * 错误处理策略
 *
 * 定义不同错误类型的处理方式。
 */
export const ErrorHandlingStrategy = {
  /** 记录警告，跳过无效消息 */
  [ErrorCodes.INVALID_MESSAGE]: 'warn_and_skip',
  /** 记录警告，跳过该组件 */
  [ErrorCodes.UNKNOWN_COMPONENT]: 'warn_and_skip',
  /** 返回 null/默认值 */
  [ErrorCodes.INVALID_PATH]: 'return_default',
  /** 抛出错误，中断渲染 */
  [ErrorCodes.CIRCULAR_REFERENCE]: 'throw',
  /** 抛出错误 */
  [ErrorCodes.REGISTRY_FROZEN]: 'throw',
  /** 抛出错误 */
  [ErrorCodes.INVALID_TYPE_NAME]: 'throw',
  /** 清理内容，记录警告 */
  [ErrorCodes.XSS_DETECTED]: 'sanitize_and_warn',
  /** 返回空字符串，记录警告 */
  [ErrorCodes.INVALID_URL]: 'return_empty_and_warn',
  /** 记录警告，跳过该组件 */
  [ErrorCodes.INVALID_PROPS]: 'warn_and_skip',
  /** 记录警告，跳过操作 */
  [ErrorCodes.INVALID_SURFACE_ID]: 'warn_and_skip',
  /** 记录警告，跳过该子组件 */
  [ErrorCodes.INVALID_CHILD_REFERENCE]: 'warn_and_skip',
} as const;

export type ErrorHandlingAction = typeof ErrorHandlingStrategy[keyof typeof ErrorHandlingStrategy];

/**
 * 获取错误的处理策略
 *
 * @param code - 错误代码
 * @returns 处理策略
 */
export function getErrorHandlingStrategy(code: ErrorCode): ErrorHandlingAction {
  return ErrorHandlingStrategy[code] || 'warn_and_skip';
}

/**
 * 根据错误处理策略处理错误
 *
 * @param error - 协议错误
 * @returns 是否应该继续执行（true）或中断（false）
 */
export function handleProtocolError(error: ProtocolError): boolean {
  const strategy = getErrorHandlingStrategy(error.code);

  switch (strategy) {
    case 'throw':
      throw error;

    case 'warn_and_skip':
      console.warn(`[ProtocolError] ${error.toString()}`);
      return true;

    case 'return_default':
      console.warn(`[ProtocolError] ${error.toString()}`);
      return true;

    case 'sanitize_and_warn':
      console.warn(`[ProtocolError] ${error.toString()}`);
      return true;

    case 'return_empty_and_warn':
      console.warn(`[ProtocolError] ${error.toString()}`);
      return true;

    default:
      console.warn(`[ProtocolError] ${error.toString()}`);
      return true;
  }
}
