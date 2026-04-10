import { useEffect } from "react";

export function useUnsavedChanges(shouldWarn: boolean) {
  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (!shouldWarn) return;

      e.preventDefault();
      e.returnValue = ""; // required for Chrome
    };

    window.addEventListener("beforeunload", handleBeforeUnload);

    return () => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
    };
  }, [shouldWarn]);
}