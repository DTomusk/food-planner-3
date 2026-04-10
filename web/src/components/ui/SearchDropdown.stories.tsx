import type { Meta, StoryObj } from "@storybook/react-vite";
import SearchDropdown from "./SearchDropdown";
import { useState } from "react";
import { expect, fn, userEvent, within } from "storybook/test";

type SearchDropdownItem = { label: string; value: string };

const meta = {
    title: "UI/SearchDropdown",
    component: SearchDropdown,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    args: {
        items: [
            { label: "Apple", value: "apple" },
            { label: "Banana", value: "banana" },
            { label: "Cherry", value: "cherry" },
        ],
        onSelect: fn(),
        selectedItem: null,
    },
    argTypes: {
        items: { control: "object" },
        onSelect: { action: "selected" },
        selectedItem: { control: "object" },
    },
} satisfies Meta<typeof SearchDropdown>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    render: (args) => {
        const [selectedItem, setSelectedItem] = useState<SearchDropdownItem | null>(null);
        return (
            <SearchDropdown 
                {...args} 
                selectedItem={selectedItem} 
                onSelect={(item) => {
                    setSelectedItem(item);
                    args.onSelect(item);
                }} 
            />
        );
    }
};

export const EmptyItems: Story = {
    args: {
        items: [],
    },
};

export const FiltersFromSelectedItemOnOpen: Story = {
    render: (args) => {
        const [selectedItem, setSelectedItem] = useState<SearchDropdownItem | null>({ label: "Banana", value: "banana" });

        return (
            <SearchDropdown
                {...args}
                selectedItem={selectedItem}
                onSelect={(item) => {
                    setSelectedItem(item);
                    args.onSelect(item);
                }}
            />
        );
    },
    play: async ({ canvasElement }) => {
        const canvas = within(canvasElement);
        const input = canvas.getByRole("combobox");

        await userEvent.click(input);

        await expect(canvas.findByText("Banana")).resolves.toBeInTheDocument();
        await expect(canvas.queryByText("Apple")).not.toBeInTheDocument();
        await expect(canvas.queryByText("Cherry")).not.toBeInTheDocument();
    },
};