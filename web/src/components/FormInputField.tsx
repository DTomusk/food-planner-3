import type { FieldError, UseFormRegisterReturn } from "react-hook-form"

type InputFieldProps = {
  label: string
  register: UseFormRegisterReturn
  error?: FieldError
  id: string
  placeholder?: string
}

export function InputField({ label, register, error, id, placeholder }: InputFieldProps) {
  return (
    <div className="mb-4">
      <label htmlFor={id} className="block font-medium mb-1">{label}</label>
      <input
        id={id}
        {...register}
        placeholder={placeholder}
        className="border rounded px-3 py-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      {error && <p className="text-red-600 text-sm mt-1">{error.message}</p>}
    </div>
  )
}
