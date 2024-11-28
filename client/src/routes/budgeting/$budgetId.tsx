import { getBudgetBreakdown } from "@/api/budgeting";
import { BudgetDefinitionForm } from "@/components/budget-definition-form";
import { LoadingDots } from "@/components/loading-dots";
import { TransactionsTable } from "@/components/transactions-table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Progress } from "@/components/ui/progress";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { snakeCaseToTitleCase } from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/budgeting/$budgetId")({
  component: Budget,
});

function Budget() {
  const params = Route.useParams();

  const breakdownQuery = useQuery({
    queryKey: ["breakdown"],
    queryFn: () => getBudgetBreakdown(params.budgetId),
  });

  if (breakdownQuery.isError)
    return <span>Error generating budget report</span>;

  return (
    <main className="flex flex-col gap-4 p-8">
      <div className="flex gap-4 items-center">
        <span className="text-xl text-muted-foreground">Budget Breakdown</span>
        {breakdownQuery.isLoading && <LoadingDots />}
        {breakdownQuery.data && (
          <span className="font-bold text-2xl">
            {breakdownQuery.data.budget.name}
          </span>
        )}
      </div>
      <Separator />
      <div className="grid grid-cols-2 gap-2 min-h-60">
        {breakdownQuery.isLoading &&
          Array.from({ length: 2 }).map(() => (
            <Skeleton className=" h-60 w-full" />
          ))}
        {!!breakdownQuery.data &&
          breakdownQuery.data.budget_definitions.map((definition) => (
            <Card key={definition.id}>
              <CardHeader>
                <CardTitle>{definition.name}</CardTitle>
              </CardHeader>
              <CardContent className="flex flex-col gap-2">
                <span>
                  {definition.usage / 100} / {definition.allocation / 100}
                </span>
                <Progress
                  value={(definition.usage / definition.allocation) * 100}
                />
                <span className="flex gap-1 flex-wrap">
                  {definition.categories.map((category) => (
                    <Tooltip>
                      <TooltipTrigger>
                        <Badge key={category.id} variant="outline">
                          {snakeCaseToTitleCase(category.detailed)}
                        </Badge>
                      </TooltipTrigger>
                      <TooltipContent>{category.description}</TooltipContent>
                    </Tooltip>
                  ))}
                </span>
                {definition.associated_transactions.length > 0 && (
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button variant="outline">View Transactions</Button>
                    </DialogTrigger>
                    <DialogContent className="max-w-[90dvw] max-h-[80dvh] overflow-auto">
                      <DialogHeader>
                        <DialogTitle>
                          <span className="text-muted-foreground">
                            Transactions for
                          </span>{" "}
                          {definition.name}
                        </DialogTitle>
                      </DialogHeader>
                      <TransactionsTable
                        transactions={definition.associated_transactions.sort(
                          (a, b) =>
                            new Date(b.date).getTime() -
                            new Date(a.date).getTime()
                        )}
                      />
                    </DialogContent>
                  </Dialog>
                )}
              </CardContent>
            </Card>
          ))}
      </div>
      <BudgetDefinitionForm budgetId={params.budgetId} />
    </main>
  );
}
