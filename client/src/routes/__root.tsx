import { getCurrentUser } from "@/api/users";
import { LogoutButton } from "@/components/logout-button";
import { PlaidLinkButton } from "@/components/plaid-link-button";
import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";
import { Toaster } from "@/components/ui/sonner";
import { UserProvider } from "@/components/user-provider";
import { createRootRoute, Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";

export const Route = createRootRoute({
  loader: ({ abortController }) => getCurrentUser(abortController.signal),
  component: Root,
});

function Root() {
  const user = Route.useLoaderData();

  return (
    <UserProvider user={user}>
      <nav className="px-4 py-2 flex gap-2 justify-between items-center">
        <span className="flex gap-4 items-center">
          <Link to="/" className="[&.active]:font-bold">
            Home
          </Link>
          {!!user && (
            <Link to="/budgeting" className="[&.active]:font-bold">
              Budgeting
            </Link>
          )}
        </span>
        <span className="flex gap-2 items-center">
          {!user && (
            <Link to="/auth/login" className="[&.active>*]:font-bold">
              <Button variant="outline">Login</Button>
            </Link>
          )}
          {!!user && (
            <p className="text-muted-foreground text-sm">{user.user_name}</p>
          )}
          {!!user && <LogoutButton />}
          {!!user && <PlaidLinkButton />}
          <ThemeToggle />
        </span>
      </nav>
      <hr />
      <Outlet />
      <TanStackRouterDevtools />
      <Toaster />
    </UserProvider>
  );
}
