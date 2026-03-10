// German university grading scale
export type GermanGrade = "1.0" | "1.3" | "1.7" | "2.0" | "2.3" | "2.7" | "3.0" | "3.3" | "3.7" | "4.0" | "5.0";

export interface GradeThreshold {
  grade: GermanGrade;
  minPercentage: number;
  label: string;
}

export interface GradeSchema {
  id: string;
  name: string;
  thresholds: GradeThreshold[];
}

// Default German university grading schema
export const DEFAULT_GRADE_SCHEMA: GradeSchema = {
  id: "german-university-default",
  name: "German University Scale",
  thresholds: [
    { grade: "1.0", minPercentage: 95, label: "excellent" },
    { grade: "1.3", minPercentage: 90, label: "excellent" },
    { grade: "1.7", minPercentage: 85, label: "very good" },
    { grade: "2.0", minPercentage: 80, label: "good" },
    { grade: "2.3", minPercentage: 75, label: "good" },
    { grade: "2.7", minPercentage: 70, label: "satisfactory" },
    { grade: "3.0", minPercentage: 65, label: "satisfactory" },
    { grade: "3.3", minPercentage: 60, label: "satisfactory" },
    { grade: "3.7", minPercentage: 55, label: "sufficient" },
    { grade: "4.0", minPercentage: 50, label: "sufficient" },
    { grade: "5.0", minPercentage: 0, label: "failed" },
  ],
};

// Calculate grade from percentage
export function calculateGrade(percentage: number, schema: GradeSchema = DEFAULT_GRADE_SCHEMA): GradeThreshold {
  const sortedThresholds = [...schema.thresholds].sort((a, b) => b.minPercentage - a.minPercentage);

  for (const threshold of sortedThresholds) {
    if (percentage >= threshold.minPercentage) {
      return threshold;
    }
  }

  // Fallback to lowest grade
  return sortedThresholds[sortedThresholds.length - 1];
}

// Calculate grade from score
export function calculateGradeFromScore(score: number, maxScore: number, schema: GradeSchema = DEFAULT_GRADE_SCHEMA): GradeThreshold {
  if (maxScore === 0) {
    return { grade: "5.0", minPercentage: 0, label: "not graded" };
  }
  const percentage = (score / maxScore) * 100;
  return calculateGrade(percentage, schema);
}

// Get color for grade
export function getGradeColor(grade: GermanGrade): string {
  const gradeNum = parseFloat(grade.replace(",", "."));
  if (gradeNum <= 1.5) return "text-emerald-400";
  if (gradeNum <= 2.5) return "text-green-400";
  if (gradeNum <= 3.5) return "text-amber-400";
  if (gradeNum <= 4.0) return "text-orange-400";
  return "text-rose-400";
}

export function getGradeBgColor(grade: GermanGrade): string {
  const gradeNum = parseFloat(grade.replace(",", "."));
  if (gradeNum <= 1.5) return "bg-emerald-500/20";
  if (gradeNum <= 2.5) return "bg-green-500/20";
  if (gradeNum <= 3.5) return "bg-amber-500/20";
  if (gradeNum <= 4.0) return "bg-orange-500/20";
  return "bg-rose-500/20";
}

