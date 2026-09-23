/**
 * Nomad UI Protocol - 安全性工具函数
 *
 * 提供 XSS 防护和 URL 验证功能，确保协议的安全性。
 *
 * @module protocol/security
 */

/**
 * 安全错误代码
 */
export const SecurityErrorCodes = {
  XSS_DETECTED: 'XSS_DETECTED',
  INVALID_URL: 'INVALID_URL',
} as const;

export type SecurityErrorCode = typeof SecurityErrorCodes[keyof typeof SecurityErrorCodes];

/**
 * 安全错误
 */
export class SecurityError extends Error {
  constructor(
    message: string,
    public code: SecurityErrorCode,
    public details?: Record<string, unknown>,
  ) {
    super(message);
    this.name = 'SecurityError';
  }
}

/**
 * 允许的 URL 方案白名单
 */
export const ALLOWED_URL_SCHEMES = ['https:', 'http:', 'data:'] as const;

export type AllowedUrlScheme = typeof ALLOWED_URL_SCHEMES[number];

/**
 * 安全的 data URL MIME 类型前缀
 * 只允许图片、音频、视频等安全的媒体类型
 */
const SAFE_DATA_URL_MIME_PREFIXES = [
  'image/',
  'audio/',
  'video/',
  'font/',
  'application/pdf',
  'application/json',
] as const;

/**
 * 检查 data URL 是否安全
 *
 * @param url - data URL
 * @returns 是否安全
 */
function isDataUrlSafe(url: string): boolean {
  // data:[<mediatype>][;base64],<data>
  const dataUrlMatch = url.match(/^data:([^;,]*)/i);
  if (!dataUrlMatch) {
    return false;
  }

  const mimeType = dataUrlMatch[1]?.toLowerCase() || '';

  // 空 MIME 类型默认为 text/plain，不安全
  if (!mimeType) {
    return false;
  }

  // 检查是否是安全的 MIME 类型
  return SAFE_DATA_URL_MIME_PREFIXES.some(prefix => mimeType.startsWith(prefix));
}

/**
 * 危险的 HTML 字符映射表
 */
const HTML_ESCAPE_MAP: Record<string, string> = {
  '&': '&amp;',
  '<': '&lt;',
  '>': '&gt;',
  '"': '&quot;',
  '\'': '&#x27;',
  '/': '&#x2F;',
  '`': '&#x60;',
  '=': '&#x3D;',
};

/**
 * 清理文本内容以防止 XSS 攻击
 *
 * 将所有可能导致 XSS 攻击的 HTML 特殊字符转义为安全的 HTML 实体。
 * 这包括：& < > " ' / ` =
 *
 * @param text - 要清理的文本
 * @returns 清理后的安全文本
 *
 * @example
 * ```typescript
 * sanitizeText('<script>alert("xss")</script>')
 * // 返回: '&lt;script&gt;alert(&quot;xss&quot;)&lt;&#x2F;script&gt;'
 *
 * sanitizeText('Hello, World!')
 * // 返回: 'Hello, World!'
 * ```
 */
export function sanitizeText(text: string): string {
  if (typeof text !== 'string') {
    return '';
  }

  return text.replace(/[&<>"'`=/]/g, char => HTML_ESCAPE_MAP[char] || char);
}

/**
 * 检测文本是否包含潜在的 XSS 攻击代码
 *
 * @param text - 要检测的文本
 * @returns 是否包含潜在的 XSS 攻击代码
 */
export function containsXSS(text: string): boolean {
  if (typeof text !== 'string') {
    return false;
  }

  // 检测常见的 XSS 攻击模式
  const xssPatterns = [
    /<script\b[^>]*>/i,
    /<\/script>/i,
    /javascript:/i,
    /on\w+\s*=/i, // onclick=, onerror=, etc.
    /<iframe\b[^>]*>/i,
    /<object\b[^>]*>/i,
    /<embed\b[^>]*>/i,
    /<link\b[^>]*>/i,
    /<style\b[^>]*>/i,
    /expression\s*\(/i, // CSS expression
    /url\s*\(\s*["']?\s*javascript:/i,
  ];

  return xssPatterns.some(pattern => pattern.test(text));
}

/**
 * 验证 URL 是否安全
 *
 * 检查 URL 的方案是否在允许列表中（https、http、data）。
 * 拒绝不安全的 URL 方案（如 javascript:）。
 * 对于 data: URLs，只允许安全的 MIME 类型（如 image/*）。
 *
 * @param url - 要验证的 URL
 * @returns 是否为安全的 URL
 *
 * @example
 * ```typescript
 * validateUrl('https://example.com/image.png')
 * // 返回: true
 *
 * validateUrl('javascript:alert("xss")')
 * // 返回: false
 *
 * validateUrl('data:image/png;base64,iVBORw0KGgo=')
 * // 返回: true
 * ```
 */
export function validateUrl(url: string): boolean {
  if (typeof url !== 'string' || url.trim() === '') {
    return false;
  }

  try {
    const parsed = new URL(url);

    // 检查协议是否在允许列表中
    if (!(ALLOWED_URL_SCHEMES as readonly string[]).includes(parsed.protocol)) {
      return false;
    }

    // 对于 data: URLs，需要额外检查 MIME 类型
    if (parsed.protocol === 'data:') {
      return isDataUrlSafe(url);
    }

    return true;
  }
  catch {
    // 如果 URL 解析失败，检查是否是相对 URL
    // 相对 URL 是安全的，因为它们会继承当前页面的协议
    if (url.startsWith('/') || url.startsWith('./') || url.startsWith('../')) {
      return true;
    }

    // 检查是否以不安全的协议开头
    const lowerUrl = url.toLowerCase().trim();
    if (lowerUrl.startsWith('javascript:') || lowerUrl.startsWith('vbscript:')) {
      return false;
    }

    return false;
  }
}

/**
 * 获取 URL 的协议
 *
 * @param url - URL 字符串
 * @returns 协议字符串（如 'https:'），如果无法解析则返回 null
 */
export function getUrlScheme(url: string): string | null {
  if (typeof url !== 'string' || url.trim() === '') {
    return null;
  }

  try {
    const parsed = new URL(url);
    return parsed.protocol;
  }
  catch {
    // 尝试从字符串中提取协议
    const match = url.match(/^([a-z][a-z0-9+.-]*):\/?\/?/i);
    return match && match[1] ? `${match[1].toLowerCase()}:` : null;
  }
}

/**
 * 清理 URL，移除潜在的危险部分
 *
 * @param url - 要清理的 URL
 * @returns 清理后的 URL，如果 URL 不安全则返回空字符串
 */
export function sanitizeUrl(url: string): string {
  if (!validateUrl(url)) {
    return '';
  }
  return url;
}
