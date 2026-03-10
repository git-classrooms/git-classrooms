"use client";

import * as React from "react";
import { useNavigate } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { BookOpen, GraduationCap, Home, Settings, Users } from "lucide-react";

import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
  CommandShortcut,
} from "@/components/ui/command";
import { classroomsQueryOptions } from "@/api/classroom";
import { useTranslation } from "react-i18next";

interface CommandPaletteProps {
  isAuthenticated: boolean;
}

export function CommandPalette({ isAuthenticated }: CommandPaletteProps) {
  const { t } = useTranslation();
  const [open, setOpen] = React.useState(false);
  const navigate = useNavigate();

  const { data: classrooms } = useQuery({
    ...classroomsQueryOptions(),
    enabled: isAuthenticated && open,
  });

  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((open) => !open);
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  const runCommand = React.useCallback(
    (command: () => void) => {
      setOpen(false);
      command();
    },
    [setOpen],
  );

  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <CommandInput placeholder={t("commandPalette.placeholder")} />
      <CommandList>
        <CommandEmpty>{t("commandPalette.noResults")}</CommandEmpty>

        <CommandGroup heading={t("commandPalette.navigation")}>
          <CommandItem onSelect={() => runCommand(() => navigate({ to: "/" }))}>
            <Home className="mr-2 h-4 w-4" />
            <span>{t("commandPalette.home")}</span>
          </CommandItem>
          {isAuthenticated && (
            <CommandItem onSelect={() => runCommand(() => navigate({ to: "/classrooms" }))}>
              <GraduationCap className="mr-2 h-4 w-4" />
              <span>{t("commandPalette.allClassrooms")}</span>
              <CommandShortcut>⌘C</CommandShortcut>
            </CommandItem>
          )}
        </CommandGroup>

        {isAuthenticated && classrooms && classrooms.length > 0 && (
          <>
            <CommandSeparator />
            <CommandGroup heading={t("commandPalette.recentClassrooms")}>
              {classrooms.slice(0, 5).map((classroom) => (
                <CommandItem
                  key={classroom.classroom.id}
                  onSelect={() =>
                    runCommand(() =>
                      navigate({
                        to: "/classrooms/$classroomId",
                        params: { classroomId: classroom.classroom.id },
                        search: { tab: "assignments" },
                      }),
                    )
                  }
                >
                  <BookOpen className="mr-2 h-4 w-4" />
                  <span>{classroom.classroom.name}</span>
                </CommandItem>
              ))}
            </CommandGroup>

            <CommandSeparator />
            <CommandGroup heading={t("commandPalette.quickActions")}>
              <CommandItem
                onSelect={() =>
                  runCommand(() =>
                    navigate({
                      to: "/classrooms/create",
                    }),
                  )
                }
              >
                <GraduationCap className="mr-2 h-4 w-4" />
                <span>{t("commandPalette.createClassroom")}</span>
              </CommandItem>
              {classrooms.slice(0, 3).map((classroom) => (
                <React.Fragment key={`actions-${classroom.classroom.id}`}>
                  <CommandItem
                    onSelect={() =>
                      runCommand(() =>
                        navigate({
                          to: "/classrooms/$classroomId/members",
                          params: { classroomId: classroom.classroom.id },
                        }),
                      )
                    }
                  >
                    <Users className="mr-2 h-4 w-4" />
                    <span>{t("commandPalette.membersOf", { name: classroom.classroom.name })}</span>
                  </CommandItem>
                  <CommandItem
                    onSelect={() =>
                      runCommand(() =>
                        navigate({
                          to: "/classrooms/$classroomId/settings",
                          params: { classroomId: classroom.classroom.id },
                        }),
                      )
                    }
                  >
                    <Settings className="mr-2 h-4 w-4" />
                    <span>{t("commandPalette.settingsOf", { name: classroom.classroom.name })}</span>
                  </CommandItem>
                </React.Fragment>
              ))}
            </CommandGroup>
          </>
        )}
      </CommandList>
    </CommandDialog>
  );
}
