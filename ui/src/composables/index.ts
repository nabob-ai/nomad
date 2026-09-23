/**
 * Composables 统一导出
 *
 * @module composables
 */

// 核心功能
export { useChat } from "./useChat";
export { useChat as useChatV2 } from "./useChatV2";
export { useAgentLoop } from "./useAgentLoop";
export { useNomadClient } from "./useNomadClient";
export { useMessage, useMessageList } from "./useMessage";

// UI 相关
export { useDarkMode } from "./useDarkMode";
export type { Theme } from "./useDarkMode";
export { useScroll } from "./useScroll";
export { useNotification } from "./useNotification";

// Nomad UI Protocol
export {
  useNomadDataBinding,
  useWatchDataPath,
  createDataBindingContext,
  DATA_BINDING_CONTEXT_KEY,
} from "./useNomadDataBinding";
export type {
  DataBindingContext,
  UseNomadDataBindingOptions,
  UseNomadDataBindingReturn,
} from "./useNomadDataBinding";

// UI Action Events
export {
  useUIAction,
  createUIActionContext,
  UI_ACTION_CONTEXT_KEY,
} from "./useUIAction";
export type {
  UIActionContext,
  UIActionEmitter,
  UseUIActionOptions,
  UseUIActionReturn,
} from "./useUIAction";

// Streaming State Preservation
export {
  useStreamingState,
  useScrollPreserve,
  getPreservedState,
  clearPreservedState,
} from "./useStreamingState";
export type {
  PreservedState,
  UseStreamingStateOptions,
} from "./useStreamingState";

// Theme System
export {
  useNomadTheme,
  convertProtocolStyles,
  createThemePreset,
  mergeThemeVariables,
  LIGHT_THEME,
  DARK_THEME,
} from "./useNomadTheme";
export type {
  ThemeMode,
  NomadThemeVariables,
  NomadThemePreset,
  UseNomadThemeOptions,
  UseNomadThemeReturn,
} from "./useNomadTheme";

// 工具函数
export { useLocalStorage } from "./useLocalStorage";
export { useDebounce, useDebounceFn } from "./useDebounce";
export { useThrottle, useThrottleFn } from "./useThrottle";
export { useClipboard } from "./useClipboard";
export { useWebSocket } from "./useWebSocket";
