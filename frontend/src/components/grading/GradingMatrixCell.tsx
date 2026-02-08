import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { useState, useEffect, useRef, KeyboardEvent } from "react";
import { GradingFeedbackPopover } from "./GradingFeedbackPopover";
import { GradingMatrixCellProps, getScoreIndicatorColor } from "./types";
import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { useGradeProject } from "@/api/grading";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";

const QUICK_SCORES = [0, 25, 50, 75, 100];

export function GradingMatrixCell({
  project,
  rubric,
  allRubrics,
  classroomId,
  assignmentId,
}: GradingMatrixCellProps) {
  const { t } = useTranslation("assignment");
  const rubricResult = project.gradingResult?.rubricResults?.[rubric.name];
  const initialScore = rubricResult?.score ?? 0;
  const initialFeedback = rubricResult?.feedback ?? "";

  const [score, setScore] = useState(initialScore);
  const [feedback, setFeedback] = useState(initialFeedback);
  const [isEditing, setIsEditing] = useState(false);
  const [showQuickScores, setShowQuickScores] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const { mutateAsync, isPending } = useGradeProject(classroomId, assignmentId, project.id);

  useEffect(() => {
    setScore(rubricResult?.score ?? 0);
    setFeedback(rubricResult?.feedback ?? "");
  }, [rubricResult]);

  const saveGrade = async (newScore: number, newFeedback?: string) => {
    const allGrades = allRubrics.map((r) => {
      const existingResult = project.gradingManualResults?.find((gr) => gr.rubric.id === r.id);
      const existingRubricResult = project.gradingResult?.rubricResults?.[r.name];

      if (r.id === rubric.id) {
        return {
          id: r.id,
          score: newScore,
          feedback: newFeedback ?? feedback,
        };
      } else if (existingResult) {
        return {
          id: r.id,
          score: existingResult.score,
          feedback: existingResult.feedback ?? "",
        };
      } else if (existingRubricResult) {
        return {
          id: r.id,
          score: existingRubricResult.score,
          feedback: existingRubricResult.feedback ?? "",
        };
      } else {
        return {
          id: r.id,
          score: 0,
          feedback: "",
        };
      }
    });

    try {
      await mutateAsync({ gradingManualRubrics: allGrades });
      setScore(newScore);
      if (newFeedback !== undefined) setFeedback(newFeedback);
    } catch {
      toast.error(t("grading.dashboard.gradeSaveFailed"));
      setScore(initialScore);
      setFeedback(initialFeedback);
    }
  };

  const handleScoreChange = (value: string) => {
    const num = value === "" ? 0 : Number(value);
    const clamped = Math.min(Math.max(num, 0), rubric.maxScore);
    setScore(clamped);
  };

  const handleBlur = () => {
    setIsEditing(false);
    if (score !== initialScore) {
      saveGrade(score);
    }
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      e.preventDefault();
      inputRef.current?.blur();
    }
    if (e.key === "Escape") {
      setScore(initialScore);
      setIsEditing(false);
    }
  };

  const handleQuickScore = (percentage: number) => {
    const newScore = Math.round((rubric.maxScore * percentage) / 100);
    setScore(newScore);
    saveGrade(newScore);
    setShowQuickScores(false);
  };

  const handleFeedbackChange = (newFeedback: string) => {
    setFeedback(newFeedback);
    saveGrade(score, newFeedback);
  };

  const indicatorColor = getScoreIndicatorColor(score, rubric.maxScore);
  const hasFeedback = feedback.trim().length > 0;

  return (
    <div
      className={cn(
        "relative flex items-center justify-between min-w-[120px] border-r border-border/40 transition-all group",
        "hover:bg-muted/40",
        isPending && "opacity-50 pointer-events-none"
      )}
    >
      {/* Left indicator bar */}
      <div
        className={cn(
          "absolute left-0 top-2 bottom-2 w-1 rounded-full transition-colors",
          indicatorColor
        )}
      />

      <Popover open={showQuickScores} onOpenChange={setShowQuickScores}>
        <PopoverTrigger asChild>
          <button
            className={cn(
              "flex-1 flex items-center justify-center gap-0.5 py-3 px-3 pl-4",
              "font-mono text-sm cursor-pointer transition-colors",
              "hover:bg-muted/30 rounded-sm mx-1",
              isEditing && "hidden"
            )}
            onClick={() => setShowQuickScores(true)}
            onDoubleClick={() => {
              setIsEditing(true);
              setTimeout(() => inputRef.current?.focus(), 0);
            }}
            disabled={isPending}
          >
            <span className="text-foreground font-semibold tabular-nums">{score}</span>
            <span className="text-muted-foreground/60">/</span>
            <span className="text-muted-foreground/60 tabular-nums">{rubric.maxScore}</span>
          </button>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-2" align="start" sideOffset={4}>
          <div className="flex gap-1">
            {QUICK_SCORES.map((pct) => (
              <Button
                key={pct}
                variant={Math.round((rubric.maxScore * pct) / 100) === score ? "default" : "outline"}
                size="sm"
                className="h-8 px-3 text-xs font-mono tabular-nums"
                onClick={() => handleQuickScore(pct)}
                disabled={isPending}
              >
                {pct}%
              </Button>
            ))}
          </div>
        </PopoverContent>
      </Popover>

      {isEditing && (
        <div className="flex-1 flex items-center justify-center py-2 px-2 pl-4">
          <Input
            ref={inputRef}
            type="number"
            min={0}
            max={rubric.maxScore}
            value={score}
            onChange={(e) => handleScoreChange(e.target.value)}
            onBlur={handleBlur}
            onKeyDown={handleKeyDown}
            className="h-7 w-16 text-center font-mono text-sm tabular-nums"
            disabled={isPending}
          />
        </div>
      )}

      <div className={cn(
        "pr-2 opacity-40 group-hover:opacity-100 transition-opacity",
        hasFeedback && "opacity-100"
      )}>
        <GradingFeedbackPopover
          feedback={feedback}
          onFeedbackChange={handleFeedbackChange}
          disabled={isPending}
        />
      </div>
    </div>
  );
}
