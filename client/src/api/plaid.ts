import { AccountBase, Transaction } from "plaid";

interface CreateLinkTokenResponse {
  linkToken: string;
}

export async function getLinkToken(signal?: AbortSignal): Promise<string> {
  const response = await fetch("/api/link/create", { signal });
  if (!response.ok) throw new Error("Error generating link token");
  const linkTokenResponse: CreateLinkTokenResponse = await response.json();
  return linkTokenResponse.linkToken;
}

export async function getAccounts(
  signal?: AbortSignal
): Promise<AccountBase[]> {
  const response = await fetch("/api/user-accounts", { signal, method: "GET" });
  if (!response.ok) return [];
  return await response.json();
}

export async function getTransactions(
  signal?: AbortSignal
): Promise<Transaction[]> {
  const response = await fetch("/api/transactions", { signal, method: "GET" });
  if (!response.ok) return [];
  return await response.json();
}
