import { Budget, getUserBudgets } from "@/api/budgeting";
import { BudgetForm } from "@/components/budget-form";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { createFileRoute, Link } from "@tanstack/react-router";
import React from "react";

export const Route = createFileRoute("/budgeting/")({
  loader: ({ abortController }) => getUserBudgets(abortController.signal),
  component: Budgeting,
});

function Budgeting() {
  const budgets = Route.useLoaderData();

  return (
    <main className="flex flex-col gap-4 p-8">
      <div className="grid grid-cols-3 gap-2">
        {budgets.map((budget) => (
          <Link to="/budgeting/$budgetId" params={{ budgetId: budget.id }}>
            <BudgetCard
              key={budget.id}
              budget={budget}
              className="hover:scale-[1.01] hover:shadow-lg transition-all duration-300"
            />
          </Link>
        ))}
      </div>
      <Separator />
      <BudgetForm />
    </main>
  );
}

function BudgetCard({
  budget,
  ...props
}: React.ComponentPropsWithoutRef<typeof Card> & { budget: Budget }) {
  return (
    <Card {...props}>
      <CardHeader>
        <CardTitle>{budget.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <Skeleton className="w-full h-64" />
        <span className="text-sm text-muted-foreground">
          Maybe this ^ is a cool graph?
        </span>
      </CardContent>
    </Card>
  );
}
