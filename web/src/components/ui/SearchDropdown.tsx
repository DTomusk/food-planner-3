import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from "@headlessui/react";
import clsx from "clsx";
import { Check, ChevronDown } from "lucide-react";
import { useState } from "react";
import Text from "./Text";

type SearchDropdownItem = {
    label: string
    value: string
}

type SearchDropdownProps = {
    items: SearchDropdownItem[];
    onSelect: (item: SearchDropdownItem | null) => void;
    selectedItem?: SearchDropdownItem | null;
    maxResults?: number;
}

export default function SearchDropdown({ items, onSelect, selectedItem, maxResults }: SearchDropdownProps) {
    const [query, setQuery] = useState("");

    const defaultFilter = (query: string) => 
        items.filter((item) => item.label.toLowerCase().includes(query.toLowerCase()));

    const filteredItems = query === "" ? items : defaultFilter(query);
    const limitedItems = filteredItems.slice(0, maxResults ?? filteredItems.length);

    return (
         <div className="mx-auto w-full max-w-xs">
            <Combobox value={selectedItem} onChange={onSelect} onClose={() => setQuery("")} by="value">
                <div className="relative">
                    <ComboboxInput 
                        displayValue={(item?: SearchDropdownItem) => item?.label ?? ""}
                        onChange={(event) => setQuery(event.target.value)}
                        className={clsx(
                            "w-full border border-gray-300 bg-white py-2 pl-3 pr-10 text-sm leading-5 text-gray-900 focus:border-blue-500 focus:outline-none focus-ring",
                        )}
                    />
                    <ComboboxButton className="group absolute inset-y-0 right-0 px-2.5">
                        <ChevronDown className="size-4 fill-white/60 group-data-hover:fill-white" />
                    </ComboboxButton>
                </div>

                <ComboboxOptions
                    anchor="bottom"
                    transition
                    className={clsx(
                        'w-(--input-width) border border-black bg-white p-1 [--anchor-gap:--spacing(1)] empty:invisible',
                        'transition duration-100 ease-in data-leave:data-closed:opacity-0'
                    )}    
                >
                    {limitedItems.map((item) => (
                        <ComboboxOption 
                        key={item.value} 
                        value={item}
                        className={
                            clsx(
                            "group flex cursor-pointer items-center gap-2 rounded-sm px-3 py-1.5 select-none",
                            "data-focus:bg-primary-100"
                            )
                        }
                        >
                        <Check className="invisible size-4 group-data-[selected]:visible" /> 
                        <Text variant="body">{item.label}</Text>
                        </ComboboxOption>
                    ))}
                </ComboboxOptions>
            </Combobox>
            
         </div>
    );
}