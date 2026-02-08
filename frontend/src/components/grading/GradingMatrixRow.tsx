import { Button } from "@/components/ui/button";
import { StatusBadge } from "@/components/ui/status-badge";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { ExternalLink, Bot } from "lucide-react";
import { GradingMatrixCell } from "./GradingMatrixCell";
import { GradingMatrixRowProps, isProjectGraded, getScoreIndicatorColor } from "./types";
import { calculateGradeFromScore, getGradeColor, getGradeBgColor, DEFAULT_GRADE_SCHEMA } from "./gradeSchema";

export function GradingMatrixRow({
  project,
  rubrics,
  assignment,
  classroomId,
  assignmentId,
  rowIndex,
  gradeSchema = DEFAULT_GRADE_SCHEMA,
}: GradingMatrixRowProps) {
  const isGraded = isProjectGraded(project, rubrics.length);
  const totalScore = project.gradingResult?.score ?? 0;
  const maxTotalScore = project.gradingResult?.maxScore ?? 0;
  const autoScore = project.gradingResult?.autogradingScore ?? 0;
  const autoMaxScore = project.gradingResult?.autogradingMaxScore ?? 0;
  const totalPercentage = maxTotalScore > 0 ? (totalScore / maxTotalScore) * 100 : 0;
  const gradeResult = calculateGradeFromScore(totalScore, maxTotalScore, gradeSchema);

  return (
    <div
      className={cn(
        "flex border-b border-border/30 transition-colors",
        "hover:bg-muted/20",
        "animate-in fade-in slide-in-from-left-1",
        rowIndex % 2 === 0 ? "bg-transparent" : "bg-muted/10"
      )}
      style={{ animationDelay: `${rowIndex * 20}ms` }}
    >
      {/* Fixed team column */}
      <div className="sticky left-0 z-10 flex items-center gap-3 px-4 py-2.5 min-w-[220px] max-w-[220px] bg-card/95 backdrop-blur-sm border-r border-border/40">
        <div className="w-9 h-9 rounded-lg bg-muted/80 flex items-center justify-center shrink-0 ring-1 ring-border/50">
          <span className="text-xs font-mono font-bold text-muted-foreground">
            {project.team.name.substring(0, 2).toUpperCase()}
          </span>
        </div>
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium truncate text-foreground/90">{project.team.name}</p>
          <StatusBadge
            variant={isGraded ? "success" : "warning"}
            size="sm"
            showDot
            className="mt-0.5"
          >
            {isGraded ? "Graded" : "Pending"}
          </StatusBadge>
        </div>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="ghost" size="icon" className="h-7 w-7 shrink-0 opacity-50 hover:opacity-100" asChild>
              <a href={project.webUrl} target="_blank" rel="noopener noreferrer">
                <ExternalLink className="h-3.5 w-3.5" />
              </a>
            </Button>
          </TooltipTrigger>
          <TooltipContent>View on GitLab</TooltipContent>
        </Tooltip>
      </div>

      {/* Scrollable rubric cells */}
      <div className="flex flex-1 overflow-visible">
        {rubrics.map((rubric) => (
          <GradingMatrixCell
            key={rubric.id}
            project={project}
            rubric={rubric}
            allRubrics={rubrics}
            classroomId={classroomId}
            assignmentId={assignmentId}
          />
        ))}

        {/* Auto-grading column */}
        {assignment.gradingJUnitAutoGradingActive && (
          <div className="relative flex items-center justify-center gap-1.5 px-3 py-2.5 min-w-[110px] border-r border-border/40 bg-muted/10">
            <div
              className={cn(
                "absolute left-0 top-2 bottom-2 w-1 rounded-full",
                getScoreIndicatorColor(autoScore, autoMaxScore)
              )}
            />
            <Bot className="h-3.5 w-3.5 text-emerald-500/70 shrink-0" />
            <span className="font-mono text-sm tabular-nums">
              <span className="text-foreground font-medium">{autoScore}</span>
              <span className="text-muted-foreground/50">/</span>
              <span className="text-muted-foreground/50">{autoMaxScore}</span>
            </span>
          </div>
        )}

        {/* Total column */}
        <div className="relative flex items-center justify-center px-4 py-2.5 min-w-[100px] bg-muted/20 border-r border-border/40">
          <div
            className={cn(
              "absolute left-0 top-2 bottom-2 w-1.5 rounded-full",
              getScoreIndicatorColor(totalScore, maxTotalScore)
            )}
          />
          <div className="text-center">
            <span className={cn(
              "font-mono text-base font-bold tabular-nums",
              totalPercentage >= 80 ? "text-emerald-400" :
              totalPercentage >= 50 ? "text-amber-400" :
              totalPercentage > 0 ? "text-rose-400" : "text-muted-foreground"
            )}>
              {totalScore}
            </span>
            <span className="text-muted-foreground/40 font-mono text-sm">/{maxTotalScore}</span>
          </div>
        </div>

        {/* Grade column */}
        <div className="flex items-center justify-center px-4 py-2.5 min-w-[90px] bg-muted/30">
          <Tooltip>
            <TooltipTrigger asChild>
              <div className={cn(
                "px-2.5 py-1 rounded-md text-center cursor-default",
                getGradeBgColor(gradeResult.grade)
              )}>
                <span className={cn(
                  "font-mono text-base font-bold",
                  getGradeColor(gradeResult.grade)
                )}>
                  {gradeResult.grade}
                </span>
              </div>
            </TooltipTrigger>
            <TooltipContent>
              <p className="font-medium capitalize">{gradeResult.label}</p>
              <p className="text-xs text-muted-foreground">
                {totalPercentage.toFixed(1)}% · min {gradeResult.minPercentage}%
              </p>
            </TooltipContent>
          </Tooltip>
        </div>
      </div>
    </div>
  );
}
