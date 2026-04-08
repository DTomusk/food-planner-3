type labelProps = {
  htmlFor: string;
  children: React.ReactNode;
};

export default function Label({
  htmlFor,
  children,
}: labelProps) {
  return (
    <label
        htmlFor={htmlFor}
        className="text-sm font-medium text-gray-700"
      >
        {children}
      </label>
  );
}
