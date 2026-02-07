"use client";

import * as React from "react";
import { cn } from "@/lib/utils";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface ResponsiveTableColumn<T> {
  header: string;
  accessorKey?: keyof T;
  cell?: (row: T) => React.ReactNode;
  hideOnMobile?: boolean;
  isPrimary?: boolean;
}

interface ResponsiveTableProps<T> {
  data: T[];
  columns: ResponsiveTableColumn<T>[];
  keyExtractor: (row: T) => string;
  onRowClick?: (row: T) => void;
  emptyMessage?: string;
  className?: string;
}

export function ResponsiveTable<T>({
  data,
  columns,
  keyExtractor,
  onRowClick,
  emptyMessage = "No data available.",
  className,
}: ResponsiveTableProps<T>) {
  if (data.length === 0) {
    return <div className="py-8 text-center text-muted-foreground">{emptyMessage}</div>;
  }

  const primaryColumn = columns.find((c) => c.isPrimary) || columns[0];
  const mobileColumns = columns.filter((c) => !c.hideOnMobile && c !== primaryColumn);

  return (
    <>
      {/* Desktop Table */}
      <div className={cn("hidden md:block", className)}>
        <div className="rounded-md border">
          <table className="w-full caption-bottom text-sm">
            <thead className="[&_tr]:border-b">
              <tr className="border-b transition-colors hover:bg-muted/50">
                {columns.map((column, index) => (
                  <th
                    key={index}
                    className="h-12 px-4 text-left align-middle font-medium text-muted-foreground"
                  >
                    {column.header}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="[&_tr:last-child]:border-0">
              {data.map((row) => (
                <tr
                  key={keyExtractor(row)}
                  onClick={() => onRowClick?.(row)}
                  className={cn(
                    "border-b transition-colors hover:bg-muted/50",
                    onRowClick && "cursor-pointer",
                  )}
                >
                  {columns.map((column, index) => (
                    <td key={index} className="p-4 align-middle">
                      {column.cell
                        ? column.cell(row)
                        : column.accessorKey
                          ? String(row[column.accessorKey] ?? "")
                          : null}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Mobile Card View */}
      <div className={cn("md:hidden space-y-3", className)}>
        {data.map((row) => (
          <Card
            key={keyExtractor(row)}
            onClick={() => onRowClick?.(row)}
            className={cn(onRowClick && "cursor-pointer hover:bg-muted/50 transition-colors")}
          >
            <CardHeader className="pb-2">
              <CardTitle className="text-base font-medium">
                {primaryColumn.cell
                  ? primaryColumn.cell(row)
                  : primaryColumn.accessorKey
                    ? String(row[primaryColumn.accessorKey] ?? "")
                    : null}
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-0">
              <dl className="grid grid-cols-2 gap-2 text-sm">
                {mobileColumns.map((column, index) => (
                  <div key={index} className="flex flex-col">
                    <dt className="text-muted-foreground text-xs">{column.header}</dt>
                    <dd className="font-medium">
                      {column.cell
                        ? column.cell(row)
                        : column.accessorKey
                          ? String(row[column.accessorKey] ?? "")
                          : null}
                    </dd>
                  </div>
                ))}
              </dl>
            </CardContent>
          </Card>
        ))}
      </div>
    </>
  );
}
