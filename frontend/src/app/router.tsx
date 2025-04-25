// src/app/router.tsx
import { createBrowserRouter } from "react-router-dom";
import { HomePage } from "../features/home";
import { Layout } from "../features/layout";
import { LoginPage, RegisterPage } from "../features/auth/pages";


export const appRouter = createBrowserRouter([
  {
    path: "/",
    element: <Layout />,
    children: [
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: "register",
        element: <RegisterPage />,
      },
      {
        path: "login",
        element: <LoginPage />,
      }
    ]
  }
]);
