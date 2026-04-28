import clsx from "clsx";
import Text from "./Text";
import Stack from "../layout/Stack";

type ListItem = {
    label: string;
    onClick: () => void;
    disabled?: boolean;
    selected?: boolean;
};

type ActionListProps = {
    title? : string;
    items: ListItem[];
};

export default function ActionList({ title, items }: ActionListProps) {
    return (
        <Stack space="sm">
        {title && <Text>{title}</Text>}
            <ul>
                {items.map((item) => (
                    <li key={item.label}>
                        <button
                            type="button"
                            disabled={item.disabled}
                            onClick={item.onClick}
                            className={clsx(
                                "flex w-full items-center gap-2 rounded px-2 py-2 text-left text-sm text-gray-900",
                                "hover:bg-gray-100",
                                item.disabled && "cursor-not-allowed opacity-50",
                                !item.disabled && "cursor-pointer",
                                item.selected && "bg-gray-200 font-semibold"
                            )}
                        >
                            {item.label}
                            </button>
                        </li>
                    )
                )}
            </ul>
        </Stack>
    );
}