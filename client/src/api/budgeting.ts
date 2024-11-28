import { Transaction } from "plaid";

export interface TransactionCategory {
  id: string;
  primary: string;
  detailed: string;
  description: string;
}

export async function getTransactionCategories(
  signal?: AbortSignal
): Promise<TransactionCategory[]> {
  const response = await fetch("/api/budgeting/categories", { signal });
  if (!response.ok) return [];
  return response.json();
}

export interface Budget {
  id: string;
  user_id: string;
  name: string;
  color: string;
}

export async function getUserBudgets(signal?: AbortSignal): Promise<Budget[]> {
  const response = await fetch("/api/ubudgeting/budgets", { signal });
  if (!response.ok) return [];
  return response.json();
}

export async function createBudget(
  budget: Budget,
  signal?: AbortSignal
): Promise<Budget | undefined> {
  const response = await fetch("/api/ubudgeting/budgets", {
    signal,
    method: "POST",
    body: JSON.stringify(budget),
  });
  if (!response.ok) return;
  return response.json();
}

export interface BudgetDefinition {
  id: string;
  user_id: string;
  budget_id: string;
  name: string;
  allocation: number;
}

export interface CreateBudgetDefinitionRequest extends BudgetDefinition {
  transaction_category_ids: string[];
}

export async function createBudgetDefinition(
  definition: CreateBudgetDefinitionRequest,
  signal?: AbortSignal
) {
  const response = await fetch(
    `/api/ubudgeting/budgets/${definition.budget_id}/definitions`,
    {
      signal,
      method: "POST",
      body: JSON.stringify(definition),
    }
  );
  if (!response.ok) return;
  return response.json();
}

export interface BudgetBreakdown {
  budget: Budget;
  budget_definitions: BudgetDefinitionBreakdown[];
}

export interface BudgetDefinitionBreakdown {
  id: string;
  user_id: string;
  budget_id: string;
  name: string;
  allocation: number;
  usage: number;
  categories: TransactionCategory[];
  associated_transactions: Transaction[];
}

export async function getBudgetBreakdown(
  budgetId: string,
  signal?: AbortSignal
): Promise<BudgetBreakdown | undefined> {
  const response = await fetch(
    `/api/ubudgeting/budgets/${budgetId}/breakdown`,
    {
      signal,
    }
  );
  if (!response.ok) return;
  return response.json();
}

export async function getBudgetBreakdowns(
  signal?: AbortSignal
): Promise<BudgetBreakdown[]> {
  const response = await fetch(`/api/ubudgeting/breakdowns`, {
    signal,
  });
  if (!response.ok) return [];
  return response.json();
}

export interface TransactionDataForDateRange {
  label: string;
  total: number;
  count: number;
}

export async function getTransactionChartData(): Promise<
  TransactionDataForDateRange[] | undefined
> {
  const response = await fetch("/api/transactions-chart");
  if (!response.ok) return;
  return response.json();
}
