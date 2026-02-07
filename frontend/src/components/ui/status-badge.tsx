import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const statusBadgeVariants = cva(
  "inline-flex items-center gap-1.5 px-2.5 py-0.5 text-xs font-medium font-mono tracking-wide uppercase rounded-full border transition-colors",
  {
    variants: {
      variant: {
        success:
          "bg-success/15 border-success/30 text-success",
        warning:
          "bg-warning/15 border-warning/30 text-warning",
        destructive:
          "bg-destructive/15 border-destructive/30 text-destructive",
        info:
          "bg-info/15 border-info/30 text-info",
        neutral:
          "bg-muted border-muted-foreground/30 text-muted-foreground",
      },
      size: {
        sm: "text-[10px] px-2 py-0.5",
        default: "text-xs px-2.5 py-0.5",
      },
    },
    defaultVariants: {
      variant: "neutral",
      size: "default",
    },
  }
);

export interface StatusBadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof statusBadgeVariants> {
  showDot?: boolean;
}

function StatusBadge({
  className,
  variant,
  size,
  showDot = false,
  children,
  ...props
}: StatusBadgeProps) {
  return (
    <span className={cn(statusBadgeVariants({ variant, size }), className)} {...props}>
      {showDot && (
        <span className="w-1.5 h-1.5 rounded-full bg-current" />
      )}
      {children}
    </span>
  );
}

export { StatusBadge, statusBadgeVariants };
