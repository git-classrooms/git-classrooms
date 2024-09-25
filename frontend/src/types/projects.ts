import { reversed } from "./utils";

export const Status = {
  Pending: "pending",
  Creating: "creating",
  Accepted: "accepted",
  Failed: "failed",
} as const;

export type Status = (typeof Status)[keyof typeof Status];

const ReversedStatus = reversed(Status);
export const getStatus = (status: Status) => ReversedStatus[status];

export const getStatusProps = (status: Status) => {
  switch (status) {
    case Status.Pending:
      return { color: { primary: "bg-blue-500", secondary: "bg-blue-400" }, name: "Pending" };
    case Status.Creating:
      return { color: { primary: "bg-yellow-300", secondary: "bg-yellow-200" }, name: "Creating" };
    case Status.Accepted:
      return { color: { primary: "bg-emerald-500", secondary: "bg-emerald-400" }, name: "Accepted" };
    case Status.Failed:
      return { color: { primary: "bg-red-600", secondary: "bg-red-500" }, name: "Failed" };
  }
};
