import { Assignment, ManualGradingRubric, ProjectResponse, UtilsReportDataItem } from "@/swagger-client";
import { GradeSchema } from "./gradeSchema";

export type ViewMode = "list" | "matrix" | "rubric";

export type StatusFilter = "all" | "graded" | "pending";

export type ZippedProject = ProjectResponse & { gradingResult?: UtilsReportDataItem };

export interface GradeUpdate {
  projectId: string;
  rubricId: string;
  score: number;
  feedback?: string;
}

export interface GradingContextValue {
  classroomId: string;
  assignmentId: string;
  assignment: Assignment;
  rubrics: ManualGradingRubric[];
  projects: ZippedProject[];
  filteredProjects: ZippedProject[];
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  statusFilter: StatusFilter;
  setStatusFilter: (filter: StatusFilter) => void;
  viewMode: ViewMode;
  setViewMode: (mode: ViewMode) => void;
  selectedRubricId: string | null;
  setSelectedRubricId: (id: string | null) => void;
}

export interface GradingMatrixProps {
  projects: ZippedProject[];
  rubrics: ManualGradingRubric[];
  assignment: Assignment;
  classroomId: string;
  assignmentId: string;
  gradeSchema?: GradeSchema;
}

export interface GradingMatrixCellProps {
  project: ZippedProject;
  rubric: ManualGradingRubric;
  allRubrics: ManualGradingRubric[];
  classroomId: string;
  assignmentId: string;
  onScoreChange?: (score: number, feedback?: string) => void;
}

export interface GradingMatrixRowProps {
  project: ZippedProject;
  rubrics: ManualGradingRubric[];
  assignment: Assignment;
  classroomId: string;
  assignmentId: string;
  rowIndex: number;
  gradeSchema?: GradeSchema;
}

export interface GradingByRubricProps {
  projects: ZippedProject[];
  rubrics: ManualGradingRubric[];
  selectedRubricId: string | null;
  onRubricChange: (rubricId: string) => void;
  classroomId: string;
  assignmentId: string;
}

export function isProjectGraded(project: ZippedProject, rubricCount: number): boolean {
  const gradedRubricCount = Object.keys(project.gradingResult?.rubricResults ?? {}).length;
  const hasAutoGrading = project.gradingResult?.autogradingMaxScore === 0 || project.gradingResult?.autogradingScore !== 0;
  return gradedRubricCount === rubricCount && hasAutoGrading;
}

export function getScoreIndicatorColor(score: number, maxScore: number): string {
  if (maxScore === 0) return "bg-muted-foreground/30";
  const percentage = (score / maxScore) * 100;
  if (percentage >= 80) return "bg-emerald-500";
  if (percentage >= 50) return "bg-amber-500";
  if (percentage > 0) return "bg-rose-500";
  return "bg-muted-foreground/30";
}

export function getScoreColor(score: number, maxScore: number): string {
  // No longer used for backgrounds - kept for compatibility
  if (maxScore === 0) return "";
  const percentage = (score / maxScore) * 100;
  if (percentage >= 80) return "text-emerald-400";
  if (percentage >= 50) return "text-amber-400";
  if (percentage > 0) return "text-rose-400";
  return "text-muted-foreground";
}

export function getScoreTextColor(score: number, maxScore: number): string {
  if (maxScore === 0) return "text-muted-foreground";
  const percentage = (score / maxScore) * 100;
  if (percentage >= 80) return "text-emerald-400";
  if (percentage >= 50) return "text-amber-400";
  if (percentage > 0) return "text-rose-400";
  return "text-muted-foreground";
}
