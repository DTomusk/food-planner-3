import type { Meta, StoryObj } from "@storybook/react-vite";
import MarkdownRenderer from "./MarkdownRenderer";

const meta = {
    title: "UI/MarkdownRenderer",
    component: MarkdownRenderer,
    tags: ["autodocs"],
    args: {
        content: `# Recipe Notes\n\nUse **fresh herbs** for the best flavor.\n\n- 1 tsp salt\n- 1 tbsp olive oil\n- Optional: lemon zest`,
    },
    argTypes: {
        content: { control: "text" },
    },
} satisfies Meta<typeof MarkdownRenderer>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithChecklist: Story = {
    args: {
        content: `## Prep Checklist\n\n1. Wash produce\n2. Preheat oven to 400F\n3. Set aside garnish`,
    },
};