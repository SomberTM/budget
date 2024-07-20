import { getCurrentUser } from "@/api/users";
import { LoginForm } from "@/components/login-form";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/auth/login")({
  beforeLoad: async () => {
    const user = await getCurrentUser();
    if (user) throw redirect({ to: "/auth/logout" });
  },
  component: () => (
    <main className="flex justify-center p-8">
      <Card className="min-w-96">
        <CardHeader>
          <CardTitle>Login</CardTitle>
          <CardDescription>Enter account credentials below</CardDescription>
        </CardHeader>
        <CardContent>
          <LoginForm />
        </CardContent>
      </Card>
    </main>
  ),
});
