// src/app/router.tsx
import { createBrowserRouter } from "react-router-dom";
import { HomePage } from "../features/home";

export const appRouter = createBrowserRouter([
  {
    path: "/",
    element: <HomePage />,
  }
]);
