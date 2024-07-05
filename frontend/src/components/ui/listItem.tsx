import { ReactElement } from "react";
import { TableCell } from "@/components/ui/table.tsx";

interface ListItemProps {
  leftContent?: ReactElement;
  rightContent?: ReactElement;
}

const ListItem: React.FC<ListItemProps> = ({leftContent, rightContent }) => {
  return (
    <>
      <TableCell className="p-2">
        {leftContent}
      </TableCell>
      <TableCell className="p-2 flex justify-end align-middle">
        {rightContent}
      </TableCell>
    </>
  );
};

export default ListItem;
