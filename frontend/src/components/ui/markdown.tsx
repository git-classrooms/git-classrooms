import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { cn } from "@/lib/utils";

interface MarkdownProps {
  children: string;
  className?: string;
  /** Compact mode for smaller text and tighter spacing */
  compact?: boolean;
  /** Inline mode strips block elements, keeping only inline formatting */
  inline?: boolean;
}

/**
 * Markdown component for rendering GitHub-flavored markdown content.
 *
 * Modes:
 * - Default: Full markdown rendering with all features
 * - Compact: Smaller typography for rubric descriptions, feedback
 * - Inline: Only inline elements (bold, italic, code, links) - for single-line descriptions
 */
export function Markdown({ children, className, compact = false, inline = false }: MarkdownProps) {
  if (inline) {
    return (
      <span className={className}>
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          allowedElements={["p", "strong", "em", "del", "code", "a"]}
          unwrapDisallowed
          components={{
            p: ({ children }) => <>{children}</>,
            strong: ({ children }) => <strong className="font-semibold">{children}</strong>,
            em: ({ children }) => <em className="italic">{children}</em>,
            del: ({ children }) => <del className="line-through text-muted-foreground">{children}</del>,
            code: ({ children }) => (
              <code className="px-1 py-0.5 rounded bg-muted font-mono text-[0.9em]">{children}</code>
            ),
            a: ({ href, children }) => (
              <a
                href={href}
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary hover:underline underline-offset-2"
              >
                {children}
              </a>
            ),
          }}
        >
          {children}
        </ReactMarkdown>
      </span>
    );
  }

  return (
    <div
      className={cn(
        "markdown-content text-foreground",
        compact && "markdown-compact",
        className
      )}
    >
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          // Headings
          h1: ({ children }) => (
            <h1 className={cn(
              "font-display font-bold tracking-tight mb-4 pb-2 border-b border-border/50",
              compact ? "text-lg" : "text-2xl"
            )}>
              {children}
            </h1>
          ),
          h2: ({ children }) => (
            <h2 className={cn(
              "font-display font-semibold tracking-tight mt-6 mb-3",
              compact ? "text-base" : "text-xl"
            )}>
              {children}
            </h2>
          ),
          h3: ({ children }) => (
            <h3 className={cn(
              "font-display font-semibold mt-5 mb-2",
              compact ? "text-sm" : "text-lg"
            )}>
              {children}
            </h3>
          ),
          h4: ({ children }) => (
            <h4 className={cn(
              "font-display font-medium mt-4 mb-2",
              compact ? "text-sm" : "text-base"
            )}>
              {children}
            </h4>
          ),
          h5: ({ children }) => (
            <h5 className={cn(
              "font-display font-medium mt-3 mb-1",
              compact ? "text-xs" : "text-sm"
            )}>
              {children}
            </h5>
          ),
          h6: ({ children }) => (
            <h6 className={cn(
              "font-display font-medium text-muted-foreground mt-3 mb-1",
              compact ? "text-xs" : "text-sm"
            )}>
              {children}
            </h6>
          ),

          // Paragraph
          p: ({ children }) => (
            <p className={cn(
              "leading-relaxed mb-3 last:mb-0",
              compact ? "text-xs" : "text-sm"
            )}>
              {children}
            </p>
          ),

          // Lists
          ul: ({ children }) => (
            <ul className={cn(
              "list-disc pl-5 mb-3 space-y-1",
              compact ? "text-xs" : "text-sm"
            )}>
              {children}
            </ul>
          ),
          ol: ({ children }) => (
            <ol className={cn(
              "list-decimal pl-5 mb-3 space-y-1",
              compact ? "text-xs" : "text-sm"
            )}>
              {children}
            </ol>
          ),
          li: ({ children }) => (
            <li className="leading-relaxed">{children}</li>
          ),

          // Task list items (GFM)
          input: ({ checked }) => (
            <input
              type="checkbox"
              checked={checked}
              disabled
              className={cn(
                "mr-2 rounded border-border",
                "accent-primary cursor-default"
              )}
            />
          ),

          // Blockquote
          blockquote: ({ children }) => (
            <blockquote className={cn(
              "border-l-2 border-primary/50 pl-4 my-3 text-muted-foreground italic",
              compact ? "text-xs" : "text-sm"
            )}>
              {children}
            </blockquote>
          ),

          // Code blocks
          pre: ({ children }) => (
            <pre className={cn(
              "overflow-x-auto rounded-lg bg-muted/50 border border-border/50 mb-3",
              compact ? "p-3 text-xs" : "p-4 text-sm"
            )}>
              {children}
            </pre>
          ),
          code: ({ className: codeClassName, children, ...props }) => {
            const isCodeBlock = codeClassName?.includes("language-");

            if (isCodeBlock) {
              return (
                <code
                  className={cn(
                    "font-mono leading-relaxed",
                    compact ? "text-xs" : "text-sm",
                    codeClassName
                  )}
                  {...props}
                >
                  {children}
                </code>
              );
            }

            // Inline code
            return (
              <code
                className={cn(
                  "px-1.5 py-0.5 rounded bg-muted font-mono",
                  compact ? "text-[0.85em]" : "text-[0.9em]"
                )}
                {...props}
              >
                {children}
              </code>
            );
          },

          // Inline elements
          strong: ({ children }) => <strong className="font-semibold">{children}</strong>,
          em: ({ children }) => <em className="italic">{children}</em>,
          del: ({ children }) => <del className="line-through text-muted-foreground">{children}</del>,

          // Links
          a: ({ href, children }) => (
            <a
              href={href}
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary hover:underline underline-offset-2 transition-colors"
            >
              {children}
            </a>
          ),

          // Horizontal rule
          hr: () => <hr className="my-6 border-border/50" />,

          // Tables (GFM)
          table: ({ children }) => (
            <div className="overflow-x-auto mb-3 rounded-lg border border-border/50">
              <table className={cn(
                "w-full border-collapse",
                compact ? "text-xs" : "text-sm"
              )}>
                {children}
              </table>
            </div>
          ),
          thead: ({ children }) => (
            <thead className="bg-muted/50 border-b border-border/50">
              {children}
            </thead>
          ),
          tbody: ({ children }) => <tbody>{children}</tbody>,
          tr: ({ children }) => (
            <tr className="border-b border-border/30 last:border-0">{children}</tr>
          ),
          th: ({ children }) => (
            <th className={cn(
              "text-left font-semibold",
              compact ? "px-2 py-1.5" : "px-3 py-2"
            )}>
              {children}
            </th>
          ),
          td: ({ children }) => (
            <td className={cn(
              compact ? "px-2 py-1.5" : "px-3 py-2"
            )}>
              {children}
            </td>
          ),

          // Images
          img: ({ src, alt }) => (
            <img
              src={src}
              alt={alt}
              className="max-w-full h-auto rounded-lg my-3 border border-border/50"
            />
          ),
        }}
      >
        {children}
      </ReactMarkdown>
    </div>
  );
}

export default Markdown;
