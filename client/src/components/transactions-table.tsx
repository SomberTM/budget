import type { Transaction } from "plaid";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./ui/table";
import { cn, snakeCaseToTitleCase } from "@/lib/utils";

export function TransactionsTable(props: { transactions: Transaction[] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Merchant Name</TableHead>
          <TableHead>Category</TableHead>
          <TableHead>Date</TableHead>
          <TableHead>Amount</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {props.transactions.map((transaction) => (
          <TableRow key={transaction.transaction_id}>
            <TableCell>{transaction.name}</TableCell>
            <TableCell>{transaction.merchant_name}</TableCell>
            <TableCell>
              {transaction.personal_finance_category &&
                snakeCaseToTitleCase(
                  transaction.personal_finance_category.detailed
                )}
            </TableCell>
            <TableCell>
              {transaction.datetime ?? transaction.date}
            </TableCell>
            <TableCell>
              <span
                className={cn(
                  transaction.amount < 0 ? "text-green-500" : "text-destructive"
                )}
              >
                {transaction.amount < 0 ? "+" : "-"}
                {Math.abs(transaction.amount)}
              </span>
              <span className="text-xs text-muted-foreground">
                {transaction.iso_currency_code}
              </span>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
