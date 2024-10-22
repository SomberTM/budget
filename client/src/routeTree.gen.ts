/* prettier-ignore-start */

/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file is auto-generated by TanStack Router

// Import Routes

import { Route as rootRoute } from './routes/__root'
import { Route as IndexImport } from './routes/index'
import { Route as BudgetingIndexImport } from './routes/budgeting/index'
import { Route as LinkPublicTokenImport } from './routes/link/$publicToken'
import { Route as BudgetingBudgetIdImport } from './routes/budgeting/$budgetId'
import { Route as AuthLogoutImport } from './routes/auth/logout'
import { Route as AuthLoginImport } from './routes/auth/login'
import { Route as AccountsAccountIdImport } from './routes/accounts/$accountId'

// Create/Update Routes

const IndexRoute = IndexImport.update({
  path: '/',
  getParentRoute: () => rootRoute,
} as any)

const BudgetingIndexRoute = BudgetingIndexImport.update({
  path: '/budgeting/',
  getParentRoute: () => rootRoute,
} as any)

const LinkPublicTokenRoute = LinkPublicTokenImport.update({
  path: '/link/$publicToken',
  getParentRoute: () => rootRoute,
} as any)

const BudgetingBudgetIdRoute = BudgetingBudgetIdImport.update({
  path: '/budgeting/$budgetId',
  getParentRoute: () => rootRoute,
} as any)

const AuthLogoutRoute = AuthLogoutImport.update({
  path: '/auth/logout',
  getParentRoute: () => rootRoute,
} as any)

const AuthLoginRoute = AuthLoginImport.update({
  path: '/auth/login',
  getParentRoute: () => rootRoute,
} as any)

const AccountsAccountIdRoute = AccountsAccountIdImport.update({
  path: '/accounts/$accountId',
  getParentRoute: () => rootRoute,
} as any)

// Populate the FileRoutesByPath interface

declare module '@tanstack/react-router' {
  interface FileRoutesByPath {
    '/': {
      id: '/'
      path: '/'
      fullPath: '/'
      preLoaderRoute: typeof IndexImport
      parentRoute: typeof rootRoute
    }
    '/accounts/$accountId': {
      id: '/accounts/$accountId'
      path: '/accounts/$accountId'
      fullPath: '/accounts/$accountId'
      preLoaderRoute: typeof AccountsAccountIdImport
      parentRoute: typeof rootRoute
    }
    '/auth/login': {
      id: '/auth/login'
      path: '/auth/login'
      fullPath: '/auth/login'
      preLoaderRoute: typeof AuthLoginImport
      parentRoute: typeof rootRoute
    }
    '/auth/logout': {
      id: '/auth/logout'
      path: '/auth/logout'
      fullPath: '/auth/logout'
      preLoaderRoute: typeof AuthLogoutImport
      parentRoute: typeof rootRoute
    }
    '/budgeting/$budgetId': {
      id: '/budgeting/$budgetId'
      path: '/budgeting/$budgetId'
      fullPath: '/budgeting/$budgetId'
      preLoaderRoute: typeof BudgetingBudgetIdImport
      parentRoute: typeof rootRoute
    }
    '/link/$publicToken': {
      id: '/link/$publicToken'
      path: '/link/$publicToken'
      fullPath: '/link/$publicToken'
      preLoaderRoute: typeof LinkPublicTokenImport
      parentRoute: typeof rootRoute
    }
    '/budgeting/': {
      id: '/budgeting/'
      path: '/budgeting'
      fullPath: '/budgeting'
      preLoaderRoute: typeof BudgetingIndexImport
      parentRoute: typeof rootRoute
    }
  }
}

// Create and export the route tree

export const routeTree = rootRoute.addChildren({
  IndexRoute,
  AccountsAccountIdRoute,
  AuthLoginRoute,
  AuthLogoutRoute,
  BudgetingBudgetIdRoute,
  LinkPublicTokenRoute,
  BudgetingIndexRoute,
})

/* prettier-ignore-end */

/* ROUTE_MANIFEST_START
{
  "routes": {
    "__root__": {
      "filePath": "__root.tsx",
      "children": [
        "/",
        "/accounts/$accountId",
        "/auth/login",
        "/auth/logout",
        "/budgeting/$budgetId",
        "/link/$publicToken",
        "/budgeting/"
      ]
    },
    "/": {
      "filePath": "index.tsx"
    },
    "/accounts/$accountId": {
      "filePath": "accounts/$accountId.tsx"
    },
    "/auth/login": {
      "filePath": "auth/login.tsx"
    },
    "/auth/logout": {
      "filePath": "auth/logout.tsx"
    },
    "/budgeting/$budgetId": {
      "filePath": "budgeting/$budgetId.tsx"
    },
    "/link/$publicToken": {
      "filePath": "link/$publicToken.tsx"
    },
    "/budgeting/": {
      "filePath": "budgeting/index.tsx"
    }
  }
}
ROUTE_MANIFEST_END */
