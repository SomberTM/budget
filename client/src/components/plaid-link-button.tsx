import { getLinkToken } from "@/api/plaid";
import { useQuery } from "@tanstack/react-query";
import { Button } from "./ui/button";
import { Link } from "lucide-react";
import { useNavigate } from "@tanstack/react-router";
import { usePlaidLink } from "react-plaid-link";

export function PlaidLinkButton() {
  const navigate = useNavigate();
  const linkTokenQuery = useQuery({
    queryKey: ["link-token"],
    queryFn: () => getLinkToken(),
  });

  function onSuccess(publicToken: string) {
    navigate({
      to: "/link/$publicToken",
      params: {
        publicToken,
      },
    });
  }

  const { ready, open } = usePlaidLink({
    token: linkTokenQuery.data ?? "",
    onSuccess,
    env: "sandbox",
  });

  return (
    <Button
      className="aspect-square p-1"
      disabled={!ready}
      onClick={() => open()}
    >
      <Link size={16} />
    </Button>
  );
}
