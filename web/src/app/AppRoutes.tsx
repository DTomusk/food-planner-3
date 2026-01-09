import { Route, Routes } from "react-router-dom";
import HomePage from "../pages/HomePage";
import RecipePage from "../pages/RecipePage";
import CreateRecipePage from "../pages/CreateRecipePage";

export function AppRoutes() {
  return (
    <Routes>
        <Route path="/" element={<HomePage/>}/>
        <Route path="/recipe/create" element={<CreateRecipePage/>}/>
        <Route path="/recipe/:id" element={<RecipePage/>}/>
    </Routes>
  );
}