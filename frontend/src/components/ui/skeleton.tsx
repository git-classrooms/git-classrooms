import { cn } from "@/lib/utils";

function Skeleton({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={cn("animate-pulse rounded-md bg-muted", className)} {...props} />;
}

interface TableSkeletonProps {
  columns?: number;
  rows?: number;
  showHeader?: boolean;
}

function TableSkeleton({ columns = 4, rows = 5, showHeader = true }: TableSkeletonProps) {
  return (
    <div className="w-full space-y-4">
      {showHeader && (
        <div className="flex items-center justify-between">
          <Skeleton className="h-10 w-[250px]" />
          <Skeleton className="h-10 w-[100px]" />
        </div>
      )}
      <div className="rounded-md border">
        <div className="border-b">
          <div className="flex h-12 items-center gap-4 px-4">
            {Array.from({ length: columns }).map((_, i) => (
              <Skeleton key={i} className="h-4 flex-1" />
            ))}
          </div>
        </div>
        <div>
          {Array.from({ length: rows }).map((_, rowIndex) => (
            <div key={rowIndex} className="flex h-16 items-center gap-4 border-b px-4 last:border-0">
              {Array.from({ length: columns }).map((_, colIndex) => (
                <Skeleton key={colIndex} className="h-4 flex-1" />
              ))}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

interface CardSkeletonProps {
  showHeader?: boolean;
  showFooter?: boolean;
  lines?: number;
}

function CardSkeleton({ showHeader = true, showFooter = false, lines = 3 }: CardSkeletonProps) {
  return (
    <div className="rounded-lg border bg-card p-6 shadow-xs">
      {showHeader && (
        <div className="mb-4 space-y-2">
          <Skeleton className="h-6 w-[200px]" />
          <Skeleton className="h-4 w-[300px]" />
        </div>
      )}
      <div className="space-y-3">
        {Array.from({ length: lines }).map((_, i) => (
          <Skeleton key={i} className="h-4 w-full" style={{ width: `${100 - i * 10}%` }} />
        ))}
      </div>
      {showFooter && (
        <div className="mt-6 flex justify-end gap-2">
          <Skeleton className="h-10 w-[100px]" />
          <Skeleton className="h-10 w-[100px]" />
        </div>
      )}
    </div>
  );
}

interface FormSkeletonProps {
  fields?: number;
  showSubmit?: boolean;
}

function FormSkeleton({ fields = 3, showSubmit = true }: FormSkeletonProps) {
  return (
    <div className="space-y-6">
      {Array.from({ length: fields }).map((_, i) => (
        <div key={i} className="space-y-2">
          <Skeleton className="h-4 w-[100px]" />
          <Skeleton className="h-10 w-full" />
        </div>
      ))}
      {showSubmit && (
        <div className="flex justify-end gap-2 pt-4">
          <Skeleton className="h-10 w-[100px]" />
        </div>
      )}
    </div>
  );
}

export { Skeleton, TableSkeleton, CardSkeleton, FormSkeleton };
