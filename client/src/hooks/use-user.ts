import { User } from "@/api/users";
import { CurrentUserContext } from "@/components/user-provider";
import { useContext } from "react";

export function useUserMaybe(): User | undefined {
  return useContext(CurrentUserContext);
}

export function useUser(): User {
  const user = useContext(CurrentUserContext);
  if (!user) throw new Error("There was supposed to be a valid user here :(");
  return user;
}
