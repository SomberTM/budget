import { getBudgetBreakdown } from "@/api/budgeting";
import { BudgetDefinitionForm } from "@/components/budget-definition-form";
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
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { snakeCaseToTitleCase } from "@/lib/utils";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/budgeting/$budgetId")({
  loader: ({ params, abortController }) =>
    getBudgetBreakdown(params.budgetId, abortController.signal),
  component: Budget,
});

function Budget() {
  const breakdown = Route.useLoaderData();

  if (!breakdown) return <span>Error generating budget report</span>;

  return (
    <main className="flex flex-col gap-4 p-8">
      <div className="flex gap-4 items-center">
        <span className="text-xl text-muted-foreground">Budget Breakdown</span>
        <span className="font-bold text-2xl">{breakdown.budget.name}</span>
      </div>
      <Separator />
      <div className="grid grid-cols-2 gap-2">
        {breakdown.budget_definitions.map((definition) => (
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
                  <DialogContent className="max-w-[75dvw]">
                    <DialogHeader>
                      <DialogTitle>
                        <span className="text-muted-foreground">
                          Transactions for
                        </span>{" "}
                        {definition.name}
                      </DialogTitle>
                    </DialogHeader>
                    <TransactionsTable
                      transactions={definition.associated_transactions}
                    />
                  </DialogContent>
                </Dialog>
                // <Accordion type="multiple" className="mt-auto">
                //   <AccordionItem value="transactions" className="border-none">
                //     <AccordionTrigger>Transactions</AccordionTrigger>
                //     <AccordionContent>
                //       <TransactionsTable
                //         transactions={definition.associated_transactions}
                //       />
                //     </AccordionContent>
                //   </AccordionItem>
                // </Accordion>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
      <BudgetDefinitionForm budgetId={breakdown.budget.id} />
    </main>
  );
}
