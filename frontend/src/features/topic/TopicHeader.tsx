import { ArrowLeft } from "lucide-react";
import { useNavigate } from "react-router-dom";

import type { TopicDetailViewModel } from "@/types/topic";

type TopicHeaderProps = {
  topic: TopicDetailViewModel;
};

export function TopicHeader({ topic }: TopicHeaderProps) {
  const navigate = useNavigate();

  return (
    <header className="border-b border-[#e6e6f3] bg-white">
      <div className="mx-auto flex min-h-[72px] max-w-[1200px] items-center gap-5 px-4 md:px-6">
        <button
          aria-label="Go back"
          className="grid size-9 shrink-0 cursor-pointer place-items-center rounded-full text-[#6d7088] transition-colors hover:bg-[#f0f1fb]"
          onClick={() => navigate("/roadmap")}
          type="button"
        >
          <ArrowLeft className="size-5" aria-hidden="true" />
        </button>

        <div className="flex min-w-0 items-center gap-3">
          <span className="text-2xl leading-none" aria-hidden="true">
            {topic.icon}
          </span>
          <div className="min-w-0">
            <h1 className="truncate text-xl font-bold tracking-normal text-[#10111f]">
              {topic.title}
            </h1>
            <p className="mt-1 text-base font-medium text-[#676982]">
              {topic.band}
            </p>
          </div>
        </div>
      </div>
    </header>
  );
}
