import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "./ui/form";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Budget, createBudget } from "@/api/budgeting";

const formSchema = z.object({
  name: z.string().min(1),
  color: z.string(),
});

type FormSchema = z.infer<typeof formSchema>;

export function BudgetForm() {
  const form = useForm<FormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      color: "#FFFFFF",
    },
  });

  const queryClient = useQueryClient();
  const createBudgetMutation = useMutation({
    mutationKey: ["create-budget"],
    mutationFn: async (budget: Budget) => {
      await createBudget(budget);
      queryClient.invalidateQueries({ queryKey: ["budgets"] });
    },
  });

  async function onSubmit(values: FormSchema) {
    const budget = values as Budget;
    return createBudgetMutation.mutateAsync(budget);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-2">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input {...field} required />
              </FormControl>
              <FormDescription>
                The name of your new budget. You will use this budget later on
                to create budgeting categories to automatically associate
                transactions with this budget.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        {/* <FormField
          control={form.control}
          name="color"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Color</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        /> */}
        <Button type="submit" disabled={createBudgetMutation.isPending}>
          Create
        </Button>
      </form>
    </Form>
  );
}
