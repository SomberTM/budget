import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/auth/logout')({
  component: () => <div>Hello /auth/logout!</div>
})