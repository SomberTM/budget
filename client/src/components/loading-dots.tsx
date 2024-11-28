import { cn } from "@/lib/utils";
import React from "react";

export function LoadingDots(props: React.ComponentPropsWithoutRef<"div">) {
  return (
    <div className="flex gap-1 justify-center items-center bg-background text-primary">
      <span className="sr-only">Loading...</span>
      <div
        {...props}
        className={cn(
          "h-2 w-2 bg-primary",
          props.className,
          "rounded-full animate-bounce [animation-delay:-0.3s]"
        )}
      ></div>
      <div
        {...props}
        className={cn(
          "h-2 w-2 bg-primary",
          props.className,
          "rounded-full animate-bounce [animation-delay:-0.15s]"
        )}
      ></div>
      <div
        {...props}
        className={cn(
          "h-2 w-2 bg-primary",
          props.className,
          "rounded-full animate-bounce"
        )}
      ></div>
    </div>
  );
}
