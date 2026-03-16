import type { Meta, StoryObj } from "@storybook/react-vite";
import { MemoryRouter } from "react-router-dom";
import BackLink from "./BackLink";

const meta = {
    title: "UI/BackLink",
    component: BackLink,
    tags: ["autodocs"],
    parameters: { layout: "centered" },
    decorators: [
        (Story) => (
            <MemoryRouter initialEntries={["/recipes/123/edit"]}>
                <Story />
            </MemoryRouter>
        ),
    ],
} satisfies Meta<typeof BackLink>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};