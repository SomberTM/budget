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
import { useMutation, useQuery } from "@tanstack/react-query";
import {
  createBudgetDefinition,
  CreateBudgetDefinitionRequest,
  getTransactionCategories,
} from "@/api/budgeting";
import { useMemo } from "react";
import { groupBy, snakeCaseToTitleCase } from "@/lib/utils";
import { Checkbox } from "./ui/checkbox";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "./ui/accordion";

const formSchema = z.object({
  budget_id: z.string().uuid(),
  name: z.string().min(1),
  allocation: z.number().min(1),
  transaction_category_ids: z.array(z.string().or(z.undefined())),
});

type FormSchema = z.infer<typeof formSchema>;

export function BudgetDefinitionForm(props: { budgetId: string }) {
  const form = useForm<FormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      budget_id: props.budgetId,
      name: "",
      allocation: 100.0,
      transaction_category_ids: [],
    },
  });

  const transactionCategoriesQuery = useQuery({
    queryKey: ["transaction-categories"],
    queryFn: () => getTransactionCategories(),
  });

  const groupedTransactionCategories = useMemo(() => {
    if (!transactionCategoriesQuery.data) return {};
    return groupBy(transactionCategoriesQuery.data, (item) => item.primary);
  }, [transactionCategoriesQuery.data]);

  const createDefinitionMutation = useMutation({
    mutationKey: ["create-definition"],
    mutationFn: (definition: CreateBudgetDefinitionRequest) =>
      createBudgetDefinition(definition),
  });

  async function onSubmit(values: FormSchema) {
    values.transaction_category_ids = values.transaction_category_ids.filter(
      (id) => !!id
    ) as string[];
    values.allocation = Number(values.allocation.toFixed(2)) * 100;

    const request = values as CreateBudgetDefinitionRequest;
    return createDefinitionMutation.mutateAsync(request);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-2">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor="name">Name</FormLabel>
              <FormControl>
                <Input id="name" {...field} required />
              </FormControl>
              <FormDescription>
                The name of this budgeting category.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="allocation"
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor="limit">Limit</FormLabel>
              <FormControl>
                <Input
                  id="limit"
                  required
                  {...form.register(field.name, {
                    valueAsNumber: true,
                  })}
                />
              </FormControl>
              <FormDescription>
                The max amount you wish to allocate to this category.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        {transactionCategoriesQuery.data && (
          <Accordion type="multiple">
            {Object.entries(groupedTransactionCategories).map(
              ([primary, categories]) => (
                <AccordionItem
                  key={primary}
                  value={primary}
                  className="border-none"
                >
                  <AccordionTrigger className="py-1">
                    {snakeCaseToTitleCase(primary)}
                  </AccordionTrigger>
                  <AccordionContent className="pl-4 flex flex-col gap-1">
                    {categories.map((category) => {
                      const index =
                        transactionCategoriesQuery.data.indexOf(category);

                      return (
                        <FormField
                          key={category.id}
                          control={form.control}
                          name={`transaction_category_ids.${index}`}
                          render={({ field }) => (
                            <FormItem>
                              <span className="flex items-center gap-1">
                                <FormControl>
                                  <Checkbox
                                    checked={!!field.value}
                                    onCheckedChange={(checked) => {
                                      if (checked) field.onChange(category.id);
                                      else field.onChange();
                                    }}
                                  />
                                </FormControl>
                                <FormLabel>
                                  {snakeCaseToTitleCase(category.detailed)}
                                </FormLabel>
                              </span>
                            </FormItem>
                          )}
                        />
                      );
                    })}
                  </AccordionContent>
                </AccordionItem>
              )
            )}
          </Accordion>
        )}
        <Button type="submit">Create</Button>
      </form>
    </Form>
  );
}
