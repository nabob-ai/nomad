import { ref, onMounted, onUnmounted, type Ref } from "vue";

/**
 * 定义 useScroll 的显式返回值类型
 * 显式指明类型可以彻底消除 PNPM 的 TS2742 "不可移植" 推导错误
 */
export interface UseScrollReturn<T> {
  target: Ref<T | null>;
  isAtBottom: Ref<boolean>;
  isAtTop: Ref<boolean>;
  scrollToTop: (smooth?: boolean) => void;
  scrollToBottom: (smooth?: boolean) => void;
  refresh: () => void;
}

/**
 * useScroll
 * 监听元素滚动并提供滚动控制方法
 */
export function useScroll<T extends HTMLElement>(): UseScrollReturn<T> {
  const target = ref<T | null>(null);
  const isAtBottom = ref(true);
  const isAtTop = ref(true);

  const updateState = () => {
    if (!target.value) return;
    const { scrollTop, scrollHeight, clientHeight } = target.value;
    isAtTop.value = scrollTop <= 0;
    isAtBottom.value = scrollTop + clientHeight >= scrollHeight - 10;
  };

  const scrollToTop = (smooth = true) => {
    if (!target.value) return;
    target.value.scrollTo({ top: 0, behavior: smooth ? "smooth" : "auto" });
  };

  const scrollToBottom = (smooth = true) => {
    if (!target.value) return;
    target.value.scrollTo({
      top: target.value.scrollHeight,
      behavior: smooth ? "smooth" : "auto",
    });
  };

  onMounted(() => {
    if (!target.value) return;
    target.value.addEventListener("scroll", updateState);
  });

  onUnmounted(() => {
    if (!target.value) return;
    target.value.removeEventListener("scroll", updateState);
  });

  return {
    target,
    isAtBottom,
    isAtTop,
    scrollToTop,
    scrollToBottom,
    refresh: updateState,
  } as unknown as UseScrollReturn<T>;
}
