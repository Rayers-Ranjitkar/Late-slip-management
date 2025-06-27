import { createBrowserRouter } from "react-router-dom";

import SignUp from "../Pages/SignUp";
import Login from "../Pages/Login"
import AdminDashboard from "../Pages/AdminDashboard";

const router = createBrowserRouter([
  {
    path: "/",
    element: <SignUp />,
  },
  {
    path: "/login",
    element: <Login />,
  },
  {
    path: "/signUp",
    element: <SignUp />,
  },
  {
    path: "/adminDashboard",
    element: <AdminDashboard />
  }
]);
export default router