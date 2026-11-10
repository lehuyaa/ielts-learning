import { GraduationCap, Sprout } from "lucide-react";

import {
  useCourses,
  useSetActiveCourse,
} from "@/features/roadmap/hooks/useCourses";
import { cn } from "@/lib/utils";
import type { CourseItem } from "@/types/course";

export function CourseSwitcher() {
  const coursesQuery = useCourses();
  const setActiveCourse = useSetActiveCourse();
  const courses = coursesQuery.data?.items ?? [];

  if (courses.length < 2) {
    return null;
  }

  return (
    <div
      aria-label="Choose course"
      className="relative z-10 mx-auto mb-6 flex w-full max-w-2xl gap-2 rounded-2xl border border-border bg-white p-1.5 shadow-sm"
      role="tablist"
    >
      {courses.map((course) => (
        <CourseTab
          course={course}
          disabled={setActiveCourse.isPending}
          key={course.id}
          onSelect={() => {
            if (!course.isActive) {
              setActiveCourse.mutate(course.id);
            }
          }}
        />
      ))}
    </div>
  );
}

function CourseTab({
  course,
  disabled,
  onSelect,
}: {
  course: CourseItem;
  disabled: boolean;
  onSelect: () => void;
}) {
  const Icon = course.kind === "CEFR" ? Sprout : GraduationCap;

  return (
    <button
      aria-selected={course.isActive}
      className={cn(
        "flex flex-1 cursor-pointer items-center gap-3 rounded-xl px-4 py-2.5 text-left transition-colors disabled:cursor-wait",
        course.isActive
          ? "bg-primary text-primary-foreground shadow-sm"
          : "text-foreground hover:bg-muted",
      )}
      disabled={disabled}
      onClick={onSelect}
      role="tab"
      type="button"
    >
      <Icon className="size-5 shrink-0" aria-hidden="true" />
      <span className="min-w-0">
        <span className="block truncate text-sm font-bold">{course.title}</span>
        <span
          className={cn(
            "block text-xs",
            course.isActive
              ? "text-primary-foreground/80"
              : "text-muted-foreground",
          )}
        >
          {course.completedLessons}/{course.totalLessons} lessons
        </span>
      </span>
    </button>
  );
}
