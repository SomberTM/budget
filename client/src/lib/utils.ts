import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function snakeCaseToTitleCase(snakeCase: string): string {
  const lower = snakeCase.toLowerCase();
  const words = lower.split("_");

  let result = "";
  for (const word of words) {
    if (result.length > 0) result += " ";
    result += word.charAt(0).toUpperCase() + word.slice(1);
  }

  return result;
}

export function groupBy<T>(
  items: T[],
  keySelector: (item: T) => string
): Record<string, T[]> {
  const groups: Record<string, T[]> = {};
  for (const item of items) {
    const key = keySelector(item);
    if (!groups[key]) groups[key] = [item];
    else groups[key].push(item);
  }
  return groups;
}
