import { z } from "zod";
import { Api } from ".";

const userSchema = z.object({
  id: z.string().uuid(),
  user_name: z.string(),
});

export type User = z.infer<typeof userSchema>;

export async function getCurrentUser(signal?: AbortSignal) {
  return Api.get<User>("/me", { signal });
}

export async function login(userName: string, password: string) {
  return Api.post("/auth/login", undefined, {
    headers: { Authorization: `Basic ${btoa(userName + ":" + password)}` },
  });
}

export async function logout() {
  return Api.post("/auth/logout");
}
