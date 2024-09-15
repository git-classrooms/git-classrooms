import { ReactNode } from "react";
import { Table, TableBody, TableRow } from "@/components/ui/table.tsx";

interface ListProps<T> {
  items: T[];
  renderItem: (item: T) => ReactNode;
}

const List = <T,>({ items, renderItem }: ListProps<T>) => {
  return (
    <Table>
      <TableBody>
        {items.map((item, index) => (
          <TableRow key={index}>{renderItem(item)}</TableRow>
        ))}
      </TableBody>
    </Table>
  );
};

export default List;
