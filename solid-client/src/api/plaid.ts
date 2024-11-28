interface CreateLinkTokenResponse {
  linkToken: string;
}

export async function getLinkToken(signal?: AbortSignal): Promise<string> {
  const response = await fetch("/api/link/create", { signal });
  if (!response.ok) throw new Error("Error generating link token");
  const linkTokenResponse: CreateLinkTokenResponse = await response.json();
  return linkTokenResponse.linkToken;
}
