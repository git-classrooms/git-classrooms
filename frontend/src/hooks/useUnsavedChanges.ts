import { useBlocker } from "@tanstack/react-router";

interface UseUnsavedChangesOptions {
  blocking: boolean;
  message?: string;
}

export function useUnsavedChanges({
  blocking,
  message = "You have unsaved changes. Are you sure you want to leave?",
}: UseUnsavedChangesOptions) {
  return useBlocker({
    shouldBlockFn: () => !window.confirm(message),
    enableBeforeUnload: true,
    disabled: !blocking,
  });
}
