import { Card, CardContent } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { StatusBadge } from "@/components/ui/status-badge";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { ExternalLink, Target, ChevronRight, ChevronLeft } from "lucide-react";
import { useState, useEffect, useMemo, useRef } from "react";
import { GradingByRubricProps, getScoreColor, getScoreTextColor } from "./types";
import { GradingFeedbackPopover } from "./GradingFeedbackPopover";
import { useGradeProject } from "@/api/grading";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";

const QUICK_SCORES = [0, 25, 50, 75, 100];

export function GradingByRubric({
  projects,
  rubrics,
  selectedRubricId,
  onRubricChange,
  classroomId,
  assignmentId,
}: GradingByRubricProps) {
  const { t } = useTranslation("assignment");
  const selectedRubric = rubrics.find((r) => r.id === selectedRubricId) ?? rubrics[0];

  const sortedProjects = useMemo(() => {
    return [...projects].sort((a, b) => {
      const aGraded = a.gradingResult?.rubricResults?.[selectedRubric?.name ?? ""] !== undefined;
      const bGraded = b.gradingResult?.rubricResults?.[selectedRubric?.name ?? ""] !== undefined;
      if (aGraded !== bGraded) return aGraded ? 1 : -1;
      return a.team.name.localeCompare(b.team.name);
    });
  }, [projects, selectedRubric]);

  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    setCurrentIndex(0);
  }, [selectedRubricId]);

  if (!selectedRubric) {
    return (
      <Card className="border-border/50">
        <CardContent className="p-12 text-center">
          <Target className="w-12 h-12 mx-auto text-muted-foreground/50 mb-4" />
          <p className="text-muted-foreground">{t("grading.dashboard.noRubricsForAssignment")}</p>
        </CardContent>
      </Card>
    );
  }

  const gradedCount = projects.filter(
    (p) => p.gradingResult?.rubricResults?.[selectedRubric.name] !== undefined
  ).length;

  const navigateTo = (index: number) => {
    if (index >= 0 && index < sortedProjects.length) {
      setCurrentIndex(index);
    }
  };

  const goToNextUngraded = () => {
    for (let i = currentIndex + 1; i < sortedProjects.length; i++) {
      const project = sortedProjects[i];
      if (!project.gradingResult?.rubricResults?.[selectedRubric.name]) {
        setCurrentIndex(i);
        return;
      }
    }
    for (let i = 0; i < currentIndex; i++) {
      const project = sortedProjects[i];
      if (!project.gradingResult?.rubricResults?.[selectedRubric.name]) {
        setCurrentIndex(i);
        return;
      }
    }
  };

  return (
    <div className="space-y-4">
      {/* Rubric selector */}
      <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
        <div className="flex items-center gap-3">
          <Select value={selectedRubric.id} onValueChange={onRubricChange}>
            <SelectTrigger className="w-[250px]">
              <Target className="w-4 h-4 mr-2 text-primary" />
              <SelectValue placeholder={t("grading.dashboard.selectRubric")} />
            </SelectTrigger>
            <SelectContent>
              {rubrics.map((rubric) => (
                <SelectItem key={rubric.id} value={rubric.id}>
                  <span>{rubric.name}</span>
                  <span className="ml-2 text-muted-foreground text-xs">
                    (max {rubric.maxScore})
                  </span>
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <div className="text-sm text-muted-foreground">
            <span className="font-mono font-medium text-foreground">{gradedCount}</span>
            <span className="mx-1">/</span>
            <span className="font-mono">{projects.length}</span>
            <span className="ml-1">{t("grading.dashboard.graded").toLowerCase()}</span>
          </div>
        </div>

        <Button variant="outline" size="sm" onClick={goToNextUngraded}>
          {t("grading.dashboard.nextUngraded")}
          <ChevronRight className="ml-1 h-4 w-4" />
        </Button>
      </div>

      {/* Project cards */}
      <div className="grid gap-3">
        {sortedProjects.map((project, index) => (
          <RubricProjectCard
            key={project.id}
            project={project}
            rubric={selectedRubric}
            allRubrics={rubrics}
            classroomId={classroomId}
            assignmentId={assignmentId}
            isActive={index === currentIndex}
            onClick={() => navigateTo(index)}
          />
        ))}
      </div>

      {/* Navigation */}
      {sortedProjects.length > 0 && (
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="icon"
            className="h-8 w-8"
            onClick={() => navigateTo(currentIndex - 1)}
            disabled={currentIndex === 0}
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <span className="text-sm text-muted-foreground px-3">
            <span className="font-mono font-medium text-foreground">{currentIndex + 1}</span>
            <span className="mx-1">/</span>
            <span className="font-mono">{sortedProjects.length}</span>
          </span>
          <Button
            variant="outline"
            size="icon"
            className="h-8 w-8"
            onClick={() => navigateTo(currentIndex + 1)}
            disabled={currentIndex === sortedProjects.length - 1}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      )}
    </div>
  );
}

interface RubricProjectCardProps {
  project: GradingByRubricProps["projects"][0];
  rubric: GradingByRubricProps["rubrics"][0];
  allRubrics: GradingByRubricProps["rubrics"];
  classroomId: string;
  assignmentId: string;
  isActive: boolean;
  onClick: () => void;
}

function RubricProjectCard({
  project,
  rubric,
  allRubrics,
  classroomId,
  assignmentId,
  isActive,
  onClick,
}: RubricProjectCardProps) {
  const { t } = useTranslation("assignment");
  const rubricResult = project.gradingResult?.rubricResults?.[rubric.name];
  const initialScore = rubricResult?.score ?? 0;
  const initialFeedback = rubricResult?.feedback ?? "";

  const [score, setScore] = useState(initialScore);
  const [feedback, setFeedback] = useState(initialFeedback);

  const { mutateAsync, isPending } = useGradeProject(classroomId, assignmentId, project.id);
  const saveTimeoutRef = useRef<ReturnType<typeof setTimeout>>(undefined);

  useEffect(() => {
    setScore(rubricResult?.score ?? 0);
    setFeedback(rubricResult?.feedback ?? "");
  }, [rubricResult]);

  useEffect(() => {
    return () => clearTimeout(saveTimeoutRef.current);
  }, []);

  const saveGrade = async (newScore: number, newFeedback?: string) => {
    // Build grades for ALL rubrics, not just existing ones
    const allGrades = allRubrics.map((r) => {
      const existingResult = project.gradingManualResults?.find((gr) => gr.rubric.id === r.id);
      const existingRubricResult = project.gradingResult?.rubricResults?.[r.name];

      if (r.id === rubric.id) {
        // This is the rubric being updated
        return {
          id: r.id,
          score: newScore,
          feedback: newFeedback ?? feedback,
        };
      } else if (existingResult) {
        // Use existing manual result
        return {
          id: r.id,
          score: existingResult.score,
          feedback: existingResult.feedback ?? "",
        };
      } else if (existingRubricResult) {
        // Use existing rubric result from report
        return {
          id: r.id,
          score: existingRubricResult.score,
          feedback: existingRubricResult.feedback ?? "",
        };
      } else {
        // No existing grade, use 0
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
      toast.success(t("grading.dashboard.gradeSavedShort"));
    } catch {
      toast.error(t("grading.dashboard.gradeSaveFailed"));
      setScore(initialScore);
      setFeedback(initialFeedback);
    }
  };

  const handleQuickScore = (percentage: number) => {
    const newScore = Math.round((rubric.maxScore * percentage) / 100);
    setScore(newScore);
    saveGrade(newScore);
  };

  const handleScoreChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const num = e.target.value === "" ? 0 : Number(e.target.value);
    const newScore = Math.min(Math.max(num, 0), rubric.maxScore);
    setScore(newScore);

    clearTimeout(saveTimeoutRef.current);
    if (newScore !== initialScore) {
      saveTimeoutRef.current = setTimeout(() => saveGrade(newScore), 300);
    }
  };

  const handleScoreBlur = () => {
    clearTimeout(saveTimeoutRef.current);
    if (score !== initialScore) {
      saveGrade(score);
    }
  };

  const handleFeedbackChange = (newFeedback: string) => {
    setFeedback(newFeedback);
    saveGrade(score, newFeedback);
  };

  const isGraded = rubricResult !== undefined;
  const colorClass = getScoreColor(score, rubric.maxScore);

  return (
    <Card
      className={cn(
        "border-border/50 cursor-pointer transition-all",
        isActive && "ring-2 ring-primary/50 border-primary/50",
        isPending && "opacity-60"
      )}
      onClick={onClick}
    >
      <CardContent className="p-4">
        <div className="flex items-center gap-4">
          {/* Team info */}
          <div className="flex items-center gap-3 min-w-[180px]">
            <div className="w-10 h-10 rounded-lg bg-muted flex items-center justify-center shrink-0">
              <span className="text-sm font-mono font-bold text-muted-foreground">
                {project.team.name.substring(0, 2).toUpperCase()}
              </span>
            </div>
            <div className="min-w-0">
              <p className="font-medium truncate">{project.team.name}</p>
              <StatusBadge variant={isGraded ? "success" : "warning"} size="sm" showDot>
                {isGraded ? t("grading.dashboard.graded") : t("grading.dashboard.pending")}
              </StatusBadge>
            </div>
          </div>

          {/* Score section */}
          <div className="flex-1 flex items-center gap-4">
            {/* Quick score buttons */}
            <div className="hidden sm:flex gap-1">
              {QUICK_SCORES.map((pct) => (
                <Button
                  key={pct}
                  variant="outline"
                  size="sm"
                  className={cn(
                    "h-7 px-2 text-xs font-mono",
                    Math.round((rubric.maxScore * pct) / 100) === score && "bg-primary text-primary-foreground"
                  )}
                  onClick={(e) => {
                    e.stopPropagation();
                    handleQuickScore(pct);
                  }}
                  disabled={isPending}
                >
                  {pct}%
                </Button>
              ))}
            </div>

            {/* Score input */}
            <div className={cn("flex items-center gap-2 px-3 py-1.5 rounded-md border", colorClass)}>
              <Input
                type="number"
                min={0}
                max={rubric.maxScore}
                value={score}
                onChange={handleScoreChange}
                onBlur={handleScoreBlur}
                onClick={(e) => e.stopPropagation()}
                className="h-7 w-16 text-center font-mono text-sm border-0 bg-transparent p-0 focus-visible:ring-0"
                disabled={isPending}
              />
              <span className={cn("text-sm", getScoreTextColor(score, rubric.maxScore))}>
                / {rubric.maxScore}
              </span>
            </div>

            {/* Feedback */}
            <div onClick={(e) => e.stopPropagation()}>
              <GradingFeedbackPopover
                feedback={feedback}
                onFeedbackChange={handleFeedbackChange}
                disabled={isPending}
              />
            </div>
          </div>

          {/* External link */}
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 shrink-0"
                asChild
                onClick={(e) => e.stopPropagation()}
              >
                <a href={project.webUrl} target="_blank" rel="noopener noreferrer">
                  <ExternalLink className="h-4 w-4" />
                </a>
              </Button>
            </TooltipTrigger>
            <TooltipContent>{t("grading.dashboard.viewOnGitLab")}</TooltipContent>
          </Tooltip>
        </div>
      </CardContent>
    </Card>
  );
}
