import { User } from "@/api/users";
import React, { createContext } from "react";

export const CurrentUserContext = createContext<User | undefined>(undefined);

export function UserProvider(
  props: React.PropsWithChildren<{ user: User | undefined }>
) {
  return (
    <CurrentUserContext.Provider value={props.user}>
      {props.children}
    </CurrentUserContext.Provider>
  );
}
