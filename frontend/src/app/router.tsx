// src/app/router.tsx
import { createBrowserRouter } from "react-router-dom";
import { HomePage } from "../features/home";
import { Layout } from "../shared/components/Layout/Layout";

export const appRouter = createBrowserRouter([
  {
    path: "/",
    element: <Layout />,
    children: [
      {
        index: true,
        element: <HomePage />,
      }
    ]
  }
]);
