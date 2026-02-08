import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { AutosizeTextarea } from "@/components/ui/autosize-textarea";
import { MessageSquare, X } from "lucide-react";
import { useState } from "react";
import { cn } from "@/lib/utils";

interface GradingFeedbackPopoverProps {
  feedback: string;
  onFeedbackChange: (feedback: string) => void;
  disabled?: boolean;
}

export function GradingFeedbackPopover({
  feedback,
  onFeedbackChange,
  disabled,
}: GradingFeedbackPopoverProps) {
  const [open, setOpen] = useState(false);
  const [localFeedback, setLocalFeedback] = useState(feedback);
  const hasFeedback = feedback.trim().length > 0;

  const handleSave = () => {
    onFeedbackChange(localFeedback);
    setOpen(false);
  };

  const handleOpenChange = (isOpen: boolean) => {
    if (isOpen) {
      setLocalFeedback(feedback);
    }
    setOpen(isOpen);
  };

  return (
    <Popover open={open} onOpenChange={handleOpenChange}>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className={cn(
            "h-6 w-6 shrink-0",
            hasFeedback && "text-primary"
          )}
          disabled={disabled}
        >
          <MessageSquare className="h-3.5 w-3.5" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-72 p-3" align="end">
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium">Feedback</span>
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6"
              onClick={() => setOpen(false)}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          </div>
          <AutosizeTextarea
            value={localFeedback}
            onChange={(e) => setLocalFeedback(e.target.value)}
            placeholder="Add feedback..."
            minHeight={60}
            maxHeight={120}
            className="text-sm resize-none"
            disabled={disabled}
          />
          <div className="flex justify-end gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setOpen(false)}
            >
              Cancel
            </Button>
            <Button
              variant="default"
              size="sm"
              onClick={handleSave}
              disabled={disabled}
            >
              Save
            </Button>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  );
}
