import { classroomQueryOptions } from "@/api/classroom";
import { isStudent } from "@/lib/utils";
import { UserClassroomResponse } from "@/swagger-client";
import { notFound, Outlet, redirect } from "@tanstack/react-router";
import { createFileRoute } from "@tanstack/react-router";
import { AxiosError } from "axios";

export const Route = createFileRoute("/_auth/classrooms/$classroomId")({
  loader: async ({ context: { queryClient }, params: { classroomId } }) => {
    let userClassroom: UserClassroomResponse;
    try {
      userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(classroomId));
    } catch (e) {
      if (e instanceof AxiosError && e.response?.status === 404) {
        throw notFound();
      }
      throw e;
    }

    if (isStudent(userClassroom) && !userClassroom.team) {
      throw redirect({
        to: "/classrooms/$classroomId/teams/join",
        params: { classroomId },
      });
    }
  },
  component: Outlet,
});
