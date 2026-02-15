import { useState, useEffect } from 'react'

/**
 * Hook that debounces a value by a specified delay.
 *
 * @param value - The value to debounce
 * @param delay - The debounce delay in milliseconds
 * @returns The debounced value
 */
export function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value)

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedValue(value)
    }, delay)

    return () => {
      clearTimeout(timer)
    }
  }, [value, delay])

  return debouncedValue
}

/**
 * Hook that provides a debounced callback.
 *
 * @param callback - The callback to debounce
 * @param delay - The debounce delay in milliseconds
 * @returns A debounced version of the callback
 */
export function useDebouncedCallback<T extends (...args: unknown[]) => unknown>(
  callback: T,
  delay: number
): T {
  const [timer, setTimer] = useState<ReturnType<typeof setTimeout> | null>(null)

  const debouncedCallback = ((...args: Parameters<T>) => {
    if (timer) {
      clearTimeout(timer)
    }
    const newTimer = setTimeout(() => {
      callback(...args)
    }, delay)
    setTimer(newTimer)
  }) as T

  useEffect(() => {
    return () => {
      if (timer) {
        clearTimeout(timer)
      }
    }
  }, [timer])

  return debouncedCallback
}
