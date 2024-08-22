import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/classrooms/$classroomId/members/')({
  component: () => <div>Hello /_auth/classrooms/$classroomId/members/!</div>
})