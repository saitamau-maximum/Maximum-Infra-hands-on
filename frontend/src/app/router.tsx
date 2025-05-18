import { createBrowserRouter } from "react-router-dom";
import { HomePage } from "../features/home";
import { Layout } from "../features/layout";
import { LoginPage, RegisterPage } from "../features/auth";
import { CreateRoomPage, RoomListPage } from "../features/room";
import { RoomPage } from "../features/chatRoom"; 
import { ImageUploadPage } from "../features/icon/pages";

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
          },
          {
            path: "icon",
            element: <ImageUploadPage />
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
          {
            path: ":roomId",
            element: <RoomPage />,
          }
        ]
      },
    ]
  }
]);
