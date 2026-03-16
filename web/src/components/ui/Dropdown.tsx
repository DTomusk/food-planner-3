import { Menu, MenuButton, MenuHeading, MenuItem, MenuItems, MenuSection } from "@headlessui/react"
import { type ReactNode } from "react"
import clsx from "clsx"

type DropdownItem = {
  label: string
  onClick?: () => void
  danger?: boolean
  disabled?: boolean
  icon?: ReactNode
}

type DropdownSection = {
  title?: string
  items: DropdownItem[]
}

type DropdownProps = {
  button: ReactNode
  sections: DropdownSection[]
}

export default function Dropdown({ button, sections }: DropdownProps) {
  return (
    <Menu as="div" className="relative inline-block text-left">
      <MenuButton as="div">
        {button}
      </MenuButton>

      <MenuItems className="absolute right-0 mt-2 w-56 border bg-white py-1 shadow-lg focus:outline-none overflow-scroll max-h-60">
        {sections.map((section, sectionIndex) => (
            <MenuSection key={sectionIndex}>
                {section.title && <MenuHeading className="px-4 py-2 text-xs font-semibold text-gray-500 uppercase select-none">{section.title}</MenuHeading>}
                {section.items.map((item, itemIndex) => (
                    <MenuItem
                      as="button"
                      key={itemIndex}
                      type="button"
                      disabled={item.disabled}
                      onClick={item.onClick}
                      className={clsx(
                        "flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-gray-900",
                        "data-focus:bg-gray-100",
                        "data-disabled:cursor-not-allowed data-disabled:opacity-50",
                        item.danger && "text-red-600",
                        !item.disabled && "cursor-pointer"
                      )}
                    >
                            {item.icon}
                            <span>{item.label}</span>
                    </MenuItem>
                ))}
            </MenuSection>
        ))}
      </MenuItems>
    </Menu>
  )
}