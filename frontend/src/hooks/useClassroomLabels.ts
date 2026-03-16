import { useTranslation } from "react-i18next";
import { Role, Status } from "@/types/classroom";

export function useRoleLabels(): Record<Role, string> {
  const { t } = useTranslation("classroom");
  return {
    [Role.Owner]: t("members.role.owner"),
    [Role.Moderator]: t("members.role.moderator"),
    [Role.Student]: t("members.role.student"),
  };
}

export function useInvitationStatusLabels(): Record<Status, string> {
  const { t } = useTranslation("classroom");
  return {
    [Status.Pending]: t("invite.status.pending"),
    [Status.Accepted]: t("invite.status.accepted"),
    [Status.Rejected]: t("invite.status.rejected"),
    [Status.Revoked]: t("invite.status.revoked"),
    [Status.Failed]: t("invite.status.failed"),
  };
}
