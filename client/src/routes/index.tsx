import { getAccounts } from "@/api/plaid";
import { AccountCard } from "@/components/account-card";
import { useUserMaybe } from "@/hooks/use-user";
import { createFileRoute, Link } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  loader: ({ abortController }) => getAccounts(abortController.signal),
  component: Index,
});

function Index() {
  const user = useUserMaybe();
  const accounts = Route.useLoaderData();

  return (
    <main>
      <div className="grid grid-cols-3 p-8 gap-4">
        {!!user &&
          accounts.map((account) => (
            <Link
              key={account.account_id}
              to="/accounts/$accountId"
              params={{ accountId: account.account_id }}
            >
              <AccountCard
                className="hover:scale-[1.01] hover:shadow-lg transition-all duration-300"
                account={account}
              />
            </Link>
          ))}
      </div>
    </main>
  );
}
