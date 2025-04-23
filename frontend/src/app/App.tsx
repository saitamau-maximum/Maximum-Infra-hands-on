// src/app/App.tsx
import { RouterProvider } from "react-router-dom";
import { appRouter } from "./router.tsx";

export default function App() {
  return <RouterProvider router={appRouter} />;
}
