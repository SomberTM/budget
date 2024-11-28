import { getAccounts } from "@/api/plaid";
import { AccountCard } from "@/components/account-card";
import { useUserMaybe } from "@/hooks/use-user";
import { createFileRoute, Link } from "@tanstack/react-router";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { getBudgetBreakdowns, getTransactionChartData } from "@/api/budgeting";
import { BudgetBreakdown } from "@/components/budget-breakdown";

export const Route = createFileRoute("/")({
  loader: async ({ abortController }) => {
    const accounts = await getAccounts(abortController.signal);
    const chartData = await getTransactionChartData();
    const breakdowns = await getBudgetBreakdowns(abortController.signal);

    return { accounts, chartData, breakdowns };
  },
  component: Index,
});

function Index() {
  const user = useUserMaybe();
  const { accounts, chartData, breakdowns } = Route.useLoaderData();

  return (
    <main className="flex flex-col gap-8 p-8">
      {chartData && (
        <Card>
          <CardHeader>
            <CardTitle>Spend Over Time</CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer
              config={{ outgoing: { label: "Outgoing" } }}
              className="h-48 w-full"
            >
              <AreaChart
                accessibilityLayer
                data={chartData}
                margin={{
                  left: 12,
                  right: 12,
                }}
              >
                <CartesianGrid vertical={false} />
                <YAxis tickFormatter={(tick) => `$${tick}`} />
                <XAxis dataKey="label" tickLine={false} axisLine={false} />
                <ChartTooltip
                  cursor={false}
                  content={
                    <ChartTooltipContent formatter={(value) => `$${value}`} />
                  }
                />
                <Area dataKey="total" type="natural" />
              </AreaChart>
            </ChartContainer>
          </CardContent>
        </Card>
      )}
      {breakdowns.map((breakdown) => (
        <BudgetBreakdown breakdown={breakdown} />
      ))}
      <div className="grid grid-cols-3 gap-4">
        {/* {!!user &&
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
          ))} */}
        {/* {!!user && <BudgetDefinitionForm budgetId="" />} */}
      </div>
    </main>
  );
}
