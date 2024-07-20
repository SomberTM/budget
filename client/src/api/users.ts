import { z } from "zod";

const userSchema = z.object({
  id: z.string().uuid(),
  user_name: z.string(),
});

export type User = z.infer<typeof userSchema>;

export async function getCurrentUser(
  signal?: AbortSignal
): Promise<User | undefined> {
  const response = await fetch("/api/me", { signal });
  if (!response.ok) return undefined;
  const user = await response.json();
  return userSchema.parse(user);
}

export async function login(
  userName: string,
  password: string
): Promise<boolean> {
  const response = await fetch("/api/auth/login", {
    method: "POST",
    headers: {
      Authorization: `Basic ${btoa(userName + ":" + password)}`,
    },
  });
  if (!response.ok) return false;
  return true;
}

export async function logout() {
  await fetch("/api/auth/logout", { method: "POST" });
}
