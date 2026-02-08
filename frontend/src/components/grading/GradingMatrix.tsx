import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { Bot, Target, GraduationCap } from "lucide-react";
import { GradingMatrixRow } from "./GradingMatrixRow";
import { GradingMatrixProps } from "./types";
import { DEFAULT_GRADE_SCHEMA } from "./gradeSchema";
import { useTranslation } from "react-i18next";

export function GradingMatrix({
  projects,
  rubrics,
  assignment,
  classroomId,
  assignmentId,
  gradeSchema = DEFAULT_GRADE_SCHEMA,
}: GradingMatrixProps) {
  const { t } = useTranslation("assignment");
  const { t: tc } = useTranslation("common");

  if (projects.length === 0) {
    return (
      <Card className="border-border/50">
        <CardContent className="p-12 text-center">
          <Target className="w-12 h-12 mx-auto text-muted-foreground/50 mb-4" />
          <p className="text-muted-foreground">{t("grading.dashboard.noProjectsMatch")}</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="border-border/40 overflow-hidden bg-card/50 backdrop-blur-sm">
      <CardContent className="p-0">
        {/* Header row */}
        <div className="flex bg-muted/30 border-b border-border/60">
          {/* Fixed header for team column */}
          <div className="sticky left-0 z-20 px-4 py-3 min-w-[220px] max-w-[220px] bg-muted/50 backdrop-blur-sm border-r border-border/40">
            <span className="text-[11px] font-semibold uppercase tracking-widest text-muted-foreground">
              {t("grading.matrix.team")}
            </span>
          </div>

          {/* Scrollable rubric headers */}
          <div className="flex flex-1 overflow-visible">
            {rubrics.map((rubric) => (
              <div
                key={rubric.id}
                className="px-3 py-3 min-w-[120px] border-r border-border/40"
              >
                <p className="text-[11px] font-semibold uppercase tracking-widest text-muted-foreground truncate">
                  {rubric.name}
                </p>
                <p className="text-[10px] text-muted-foreground/50 font-mono mt-1 tabular-nums">
                  max {rubric.maxScore}
                </p>
              </div>
            ))}

            {/* Auto-grading header */}
            {assignment.gradingJUnitAutoGradingActive && (
              <div className="px-3 py-3 min-w-[110px] border-r border-border/40 bg-muted/20">
                <div className="flex items-center gap-1.5">
                  <Bot className="h-3 w-3 text-emerald-500" />
                  <span className="text-[11px] font-semibold uppercase tracking-widest text-muted-foreground">
                    Auto
                  </span>
                </div>
              </div>
            )}

            {/* Total header */}
            <div className="px-4 py-3 min-w-[100px] bg-muted/40 border-r border-border/40">
              <span className="text-[11px] font-semibold uppercase tracking-widest text-muted-foreground">
                {tc("total")}
              </span>
            </div>

            {/* Grade header */}
            <div className="px-4 py-3 min-w-[90px] bg-muted/50">
              <div className="flex items-center gap-1.5">
                <GraduationCap className="h-3.5 w-3.5 text-primary" />
                <span className="text-[11px] font-semibold uppercase tracking-widest text-muted-foreground">
                  {t("grading.grade")}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Data rows */}
        <div className={cn("max-h-[600px] overflow-auto")}>
          {projects.map((project, index) => (
            <GradingMatrixRow
              key={project.id}
              project={project}
              rubrics={rubrics}
              assignment={assignment}
              classroomId={classroomId}
              assignmentId={assignmentId}
              rowIndex={index}
              gradeSchema={gradeSchema}
            />
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
