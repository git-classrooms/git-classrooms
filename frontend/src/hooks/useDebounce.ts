import { useEffect, useState } from "react";

export const useDebounce = <T>(value: T, durationMs: number) => {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const id = setTimeout(() => setDebouncedValue(value), durationMs);
    return () => clearTimeout(id);
  });

  return debouncedValue;
};
