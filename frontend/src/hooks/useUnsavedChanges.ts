import { useBlocker } from "@tanstack/react-router";

interface UseUnsavedChangesOptions {
  isDirty: boolean;
  message?: string;
}

export function useUnsavedChanges({
  isDirty,
  message = "You have unsaved changes. Are you sure you want to leave?",
}: UseUnsavedChangesOptions) {
  return useBlocker({
    shouldBlockFn: () => !window.confirm(message),
    enableBeforeUnload: true,
    disabled: !isDirty,
  });
}
