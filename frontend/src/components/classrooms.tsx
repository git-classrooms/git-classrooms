import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowRight as ArrowRight, Plus, SearchCode } from "lucide-react";
import { useState } from "react";
import { UserClassroomResponse } from "@/swagger-client";
import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import {
  ColumnFiltersState,
  SortingState,
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { Input } from "@/components/ui/input";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Link } from "@tanstack/react-router";

const classroomColumnHelper = createColumnHelper<UserClassroomResponse>();
const classroomColumns = [
  classroomColumnHelper.accessor((row) => row.classroom.name, {
    id: "name",
    header: "Name",
    cell: ({ row: { original: classroom } }) => (
      <div className="cursor-default flex">
        <div className="pr-2">
          <Avatar>
            <AvatarFallback className="bg-gray-200 text-black text-lg">
              {classroom.classroom.name.charAt(0)}
            </AvatarFallback>
          </Avatar>
        </div>
        <div>
          <div className="font-medium">{classroom.classroom.name}</div>
          <div className="text-sm text-muted-foreground md:inline">
            {classroom.assignmentsCount} Assignment{classroom.assignmentsCount === 1 ? "" : "s"}
          </div>
        </div>
      </div>
    ),
  }),

  classroomColumnHelper.display({
    id: "actions",
    header: "Actions",
    cell: ({ row: { original: classroom } }) => (
      <div className="p-2 flex justify-end align-middle">
        <Button variant="ghost" size="icon" asChild>
          <a href={classroom.webUrl} target="_blank" rel="noreferrer">
            <SearchCode className="h-6 w-6 text-gray-600 dark:text-white" />
          </a>
        </Button>
        <Button variant="ghost" size="icon" asChild>
          <Link
            to="/classrooms/$classroomId"
            search={{ tab: "assignments" }}
            params={{ classroomId: classroom.classroom.id }}
          >
            <ArrowRight className="h-6 w-6 text-gray-600 dark:text-white" />
          </Link>
        </Button>
      </div>
    ),
  }),
];

function ClassroomTable({ classrooms }: { classrooms: UserClassroomResponse[] }) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);

  const table = useReactTable({
    data: classrooms,
    columns: classroomColumns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    state: {
      sorting,
      columnFilters,
    },
  });

  return (
    <>
      <Input
        placeholder="Filter classrooms..."
        value={(table.getColumn("team")?.getFilterValue() as string) ?? ""}
        onChange={(event) => table.getColumn("team")?.setFilterValue(event.target.value)}
        className="max-w-sm"
      />
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup, i) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                return (
                  <TableHead className={i === 0 ? "w-full" : undefined} key={header.id}>
                    {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                  </TableHead>
                );
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id} data-state={row.getIsSelected() && "selected"}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={classroomColumns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </>
  );
}

export function OwnedClassroomTable({
  classrooms,
  showAll,
}: {
  classrooms: UserClassroomResponse[];
  showAll: boolean;
}) {
  return (
    <Card>
      <CardHeader className="md:flex md:flex-row md:items-center justify-between space-y-0 pb-2 mb-4">
        <div className="mb-4 md:mb-0">
          <CardTitle className="mb-1">Managed Classrooms</CardTitle>
          <CardDescription>Classrooms which are managed by you</CardDescription>
        </div>
        <div className="flex gap-2">
          {showAll ? (
            <>
              <Button asChild variant="outline">
                <Link to="/classrooms">View all</Link>
              </Button>
              <Button asChild variant="outline" size="icon">
                <Link to="/classrooms/create">
                  <Plus />
                </Link>
              </Button>
            </>
          ) : (
            <Button asChild variant="outline">
              <Link to="/classrooms/create">
                <Plus className="h-4 w-4 mr-2" /> Create classroom
              </Link>
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <ClassroomTable classrooms={classrooms} />
      </CardContent>
    </Card>
  );
}

export function JoinedClassroomTable({ classrooms }: { classrooms: UserClassroomResponse[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Joined Classrooms</CardTitle>
        <CardDescription>Classroom of which you are a member</CardDescription>
      </CardHeader>
      <CardContent>
        <ClassroomTable classrooms={classrooms} />
      </CardContent>
    </Card>
  );
}
