import type { AccountBase } from "plaid";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";
import React from "react";

export function AccountCard({
  account,
  ...props
}: { account: AccountBase } & React.ComponentPropsWithoutRef<typeof Card>) {
  return (
    <Card {...props}>
      <CardHeader>
        <CardTitle className="flex justify-between">
          {account.name}
          <span>
            {account.balances.current}
            <span className="text-muted-foreground text-sm">
              {account.balances.iso_currency_code}
            </span>
          </span>
        </CardTitle>
        <CardDescription>{account.official_name}</CardDescription>
      </CardHeader>
      <CardContent></CardContent>
    </Card>
  );
}
