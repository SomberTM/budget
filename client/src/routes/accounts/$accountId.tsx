import { getAccounts, getAccountTransactions } from "@/api/plaid";
import { TransactionsTable } from "@/components/transactions-table";
import { Separator } from "@/components/ui/separator";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/accounts/$accountId")({
  staleTime: 60_000,
  loader: async ({ params, abortController }) => {
    const accounts = await getAccounts(abortController.signal);
    const transactions = await getAccountTransactions(
      params.accountId,
      abortController.signal
    );

    const account = accounts.find(
      (account) => account.account_id === params.accountId
    );
    if (!account) throw redirect({ to: "/" });

    return {
      account,
      transactions,
    };
  },
  component: Account,
});

function Account() {
  const { account, transactions } = Route.useLoaderData();

  return (
    <main className="p-8 flex flex-col gap-4 w-full">
      <div className="flex flex-col gap-2">
        <span>{account.name}</span>
        <span>{account.official_name}</span>
      </div>
      <Separator />
      <TransactionsTable transactions={transactions} />
    </main>
  );
}
