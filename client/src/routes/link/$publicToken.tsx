import { createFileRoute, Navigate } from "@tanstack/react-router";

async function exchangePublicToken(publicToken: string) {
  return fetch("/api/link/exchange", {
    method: "POST",
    body: JSON.stringify({ publicToken }),
  });
}

export const Route = createFileRoute("/link/$publicToken")({
  loader: ({ params }) => exchangePublicToken(params.publicToken),
  component: () => <Navigate to="/" />,
});
