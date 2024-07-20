import React from "react";
import { Button } from "./ui/button";
import { logout } from "@/api/users";
import { useNavigate } from "@tanstack/react-router";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

export function LogoutButton(
  props: React.ComponentPropsWithoutRef<typeof Button>
) {
  const navigate = useNavigate();
  const logoutMutation = useMutation({
    mutationKey: ["logout"],
    mutationFn: logout,
    onSuccess() {
      navigate({ to: "/" });
    },
    onError() {
      toast("Failed to logout.");
    },
  });

  return (
    <Button
      variant="outline"
      disabled={logoutMutation.isPending}
      onClick={() => logoutMutation.mutate()}
      {...props}
    >
      Logout
    </Button>
  );
}
