import { useMemo } from "react";
import Fuse from "fuse.js";

export interface SearchDropdownItem {
  label: string;
  value: string | number;
}

export function useFuzzyDropdown<T>(
  items: T[],
  labelFn: (item: T) => string,
  valueFn: (item: T) => string,
  threshold = 0.3
) {
  // Memoize dropdown items
  const dropdownItems = useMemo(
    () =>
      items.map((item) => ({
        label: labelFn(item),
        value: valueFn(item),
      })),
    [items, labelFn, valueFn]
  );

  // Memoize Fuse index
  const fuse = useMemo(
    () =>
      new Fuse(dropdownItems, {
        keys: ["label"],
        threshold,
      }),
    [dropdownItems, threshold]
  );

  // Filter function
  const filterFunction = useMemo(
    () => (query: string) =>
      query ? fuse.search(query).map((r) => r.item) : dropdownItems,
    [fuse, dropdownItems]
  );

  return { dropdownItems, filterFunction };
}