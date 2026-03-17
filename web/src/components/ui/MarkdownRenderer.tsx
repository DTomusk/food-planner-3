import ReactMarkdown from "react-markdown";

export default function MarkdownRenderer({ content }: { content: string }) {
  return (
    <div className="space-y-4 text-base">
      {/* TODO: consider getting prose plugin */}
      <ReactMarkdown
        components={{
          h1: ({node, ...props}) => <h1 className="text-4xl font-bold" {...props} />,
          h2: ({node, ...props}) => <h2 className="text-3xl font-semibold" {...props} />,
          h3: ({node, ...props}) => <h3 className="text-2xl font-semibold" {...props} />,
          p: ({node, ...props}) => <p className="leading-7" {...props} />,
          ul: ({node, ...props}) => <ul className="list-disc ml-6 space-y-1" {...props} />,
          ol: ({node, ...props}) => <ol className="list-decimal ml-6 space-y-1" {...props} />,
          a: ({node, ...props}) => <a className="text-primary-600 hover:underline" {...props} />,
          blockquote: ({node, ...props}) => <blockquote className="border-l-4 border-gray-300 pl-4 italic" {...props} />,
          code: ({node, inline, className, ...props}: { node?: any; inline?: boolean; className?: string } & React.HTMLAttributes<HTMLElement>) =>
            inline ? (
              <code className="bg-gray-100 px-1 rounded" {...props} />
            ) : (
              <pre className="bg-gray-100 p-3 rounded overflow-x-auto" {...props} />
            ),
        }}>
        {content}
      </ReactMarkdown>
    </div>
  );
}