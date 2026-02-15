import { useEffect, useCallback } from "react";
import { useBlocker } from "@tanstack/react-router";

interface UseUnsavedChangesOptions {
  isDirty: boolean;
  message?: string;
}

export function useUnsavedChanges({
  isDirty,
  message = "You have unsaved changes. Are you sure you want to leave?",
}: UseUnsavedChangesOptions) {
  useBlocker(() => window.confirm(message), isDirty);

  const handleBeforeUnload = useCallback(
    (e: BeforeUnloadEvent) => {
      if (isDirty) {
        e.preventDefault();
        e.returnValue = message;
        return message;
      }
    },
    [isDirty, message],
  );

  useEffect(() => {
    window.addEventListener("beforeunload", handleBeforeUnload);
    return () => window.removeEventListener("beforeunload", handleBeforeUnload);
  }, [handleBeforeUnload]);

  useEffect(() => {
    if (!isDirty) return;

    let skipNextPop = false;

    window.history.pushState(window.history.state, "", window.location.href);

    const handlePopState = () => {
      if (skipNextPop) {
        skipNextPop = false;
        return;
      }

      if (window.confirm(message)) {
        skipNextPop = true;
        window.history.back();
      } else {
        window.history.pushState(window.history.state, "", window.location.href);
      }
    };

    window.addEventListener("popstate", handlePopState);
    return () => window.removeEventListener("popstate", handlePopState);
  }, [isDirty, message]);
}
