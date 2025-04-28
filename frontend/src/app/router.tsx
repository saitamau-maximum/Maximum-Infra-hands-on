// src/app/router.tsx
import { createBrowserRouter } from "react-router-dom";
import { HomePage } from "../features/home";
import { Layout } from "../features/layout";
import { LoginPage, RegisterPage } from "../features/auth";
import { CreateRoomPage, RoomListPage } from "../features/room";


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
        path: "user",
        children: [
          {
            path: "register",
            element: <RegisterPage />,
          },
          {
            path: "login",
            element: <LoginPage />,
          }
        ]
      },
      {
        path: "room",
        children: [
          {
            index: true,
            element: <RoomListPage />,
          },
          {
            path: "create",
            element: <CreateRoomPage />,
          },
        ]
      },
    ]
  }
]);
