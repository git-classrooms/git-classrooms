import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/dashboard")({
  component: Dashboard,
});

function Dashboard() {
  return (
    <div className="p-2">
      <h3>Dashboard</h3>
    </div>
  );
}
