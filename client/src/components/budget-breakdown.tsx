import { type BudgetBreakdown } from "@/api/budgeting";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import { TransactionsTable } from "./transactions-table";
import { Button } from "./ui/button";
import { snakeCaseToTitleCase } from "@/lib/utils";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";
import { Badge } from "./ui/badge";
import { Separator } from "./ui/separator";
import { Progress } from "./ui/progress";

export function BudgetBreakdown(props: { breakdown: BudgetBreakdown }) {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex gap-4 items-center">
        <span className="text-xl text-muted-foreground">Budget Breakdown</span>
        <span className="font-bold text-2xl">
          {props.breakdown.budget.name}
        </span>
      </div>
      <Separator />
      <div className="grid grid-cols-2 gap-2">
        {props.breakdown.budget_definitions.map((definition) => (
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
    </div>
  );
}
