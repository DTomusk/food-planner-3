import { useEffect } from "react";
import { useBlocker } from "react-router-dom";

type UseUnsavedChangesOptions = {
  blockRouteChange?: boolean;
  routeChangeMessage?: string;
};

export function useUnsavedChanges(
  shouldWarn: boolean,
  {
    blockRouteChange = false,
    routeChangeMessage = "You have unsaved changes. Leave this page?",
  }: UseUnsavedChangesOptions = {},
) {
  const blocker = useBlocker(blockRouteChange && shouldWarn);

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

  useEffect(() => {
    if (blocker.state !== "blocked") return;

    if (window.confirm(routeChangeMessage)) {
      blocker.proceed();
      return;
    }

    blocker.reset();
  }, [blocker, routeChangeMessage]);
}