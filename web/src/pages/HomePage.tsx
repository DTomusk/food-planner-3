import { Container, PageTitle } from "@/components";
import SearchBar from "@/components/ui/SearchBar";
import { Page } from "@/layout";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function HomePage() {
  const [searchQuery, setSearchQuery] = useState("");
  const navigate = useNavigate();

  const handleSubmitSearch = () => {
    const trimmedQuery = searchQuery.trim();
    if (trimmedQuery) {
      navigate(`/recipes?q=${encodeURIComponent(trimmedQuery)}`);
    }
  };

  return (
    <Page>
      <PageTitle text="Welcome to FoodSmash" /> 
      <Container size="sm">
      <SearchBar placeholder="Search for recipes..." 
        value={searchQuery}
        onChange={setSearchQuery}
        onSubmit={handleSubmitSearch}
      />  
      </Container>
  </Page>
  )
}