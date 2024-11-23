import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Pagination } from "@/components/ui/pagination";

// Utility for truncating text
const truncate = (str, length) => str.length > length ? `${str.substring(0, length)}...` : str;

// Search Input Component
const SearchInput = ({ onSearch, placeholder }) => {
  const [query, setQuery] = useState('');

  const handleSearch = () => {
    onSearch(query.trim());
  };

  return (
    <div className="flex w-full max-w-md space-x-2">
      <Input 
        type="text" 
        placeholder={placeholder}
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        className="flex-grow"
      />
      <Button 
        onClick={handleSearch} 
        className="bg-blue-600 hover:bg-blue-700"
      >
        <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        Search
      </Button>
    </div>
  );
};

// Superhero Card Component
const SuperheroCard = ({ superhero }) => {
  const heroKey = superhero._id?.$oid || `${superhero.name}-${superhero.alignment}`;
  
  return (
    <Card key={heroKey} className="w-full max-w-xs transition-all hover:scale-105 hover:shadow-xl bg-red-100">
      <CardHeader className="p-4 pb-0">
        <img 
          src={superhero.image} // Directly use superhero.image as it is the URL of the image
          alt={superhero.name || 'Unknown Hero'} // Default text if name is undefined
          className="w-full h-64 object-cover rounded-t-lg"
        />
      </CardHeader>
      <CardContent className="p-4">
        <CardTitle className="text-xl font-bold mb-2 text-black">
          {truncate(superhero.name || 'Unknown Hero', 20)}  {/* Default name if undefined */}
        </CardTitle>
        <div className="space-y-2 text-black">
          <Badge variant="secondary">
            ⭐ Power: {superhero.power || "N/A"}
          </Badge>
          <Badge variant="outline">
            Alignment: {superhero.alignment || "N/A"}
          </Badge>
        </div>
      </CardContent>
    </Card>
  );
};

// Movie Card Component
const MovieCard = ({ movie }) => {
  const movieKey = movie.Title ? movie.Title : `${movie.Title || 'Unknown Movie'}-${movie.Year || 'N/A'}`;

  return (
    <Card key={movieKey} className="w-full max-w-xs transition-all hover:scale-105 hover:shadow-xl bg-blue-100">
      <CardHeader className="p-4 pb-0">
        <img 
          src={movie.Poster || '/default-poster.jpg'}  // Accessing the correct Poster property
          alt={movie.Title || 'Unknown Movie'}
          className="w-full h-64 object-cover rounded-t-lg"
        />
      </CardHeader>
      <CardContent className="p-4">
        <CardTitle className="text-xl font-bold mb-2 text-black">
          {truncate(movie.Title || 'Unknown Movie', 20)}
        </CardTitle>
        <div className="space-y-2 text-black">
          <Badge variant="secondary">
            🎬 Year: {movie.Year || "N/A"}
          </Badge>
        </div>
      </CardContent>
    </Card>
  );
};

// Main Search Page Component
const SuperheroSearchPage = () => {
  const [superheroes, setSuperheroes] = useState([]);
  const [movies, setMovies] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  const handleSearch = async (query) => {
    console.log('Search query:', query);
    try {
      // Sending the query to the backend (replace with your backend URL)
      const response = await fetch(`http://localhost:8080/api/search?query=${query}`);
      
      // Check if the response is ok (status code 200-299)
      if (!response.ok) {
        throw new Error("Failed to fetch search results");
      }
  
      // Parse the response as JSON
      const data = await response.json();

      console.log("Fetched Data:", data); 
  
      // Destructure the data to get the superheroes, movies, and pagination info
      const { superheroes, movies, totalPages } = data;
  
      // Update the state with the data from the backend
      setSuperheroes(superheroes);
      setMovies(movies);
      setTotalPages(totalPages);  // You can use actual pagination info from the backend
  
    } catch (err) {
      console.error("Error fetching data:", err);
      
      // In case of an error, reset the results
      setSuperheroes([]);
      setMovies([]);
      setTotalPages(1);  // Set a default value for totalPages (e.g., 1 or 0)
    }
  };

  const combinedData = [...superheroes, ...movies]; // Combine the superhero and movie arrays into one

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-6xl mx-auto">
        <div className="mb-8 flex justify-center items-center">
          <h1 className="text-4xl font-bold text-gray-800">
            Hero Huntr
          </h1>
        </div>

        <div className="mb-8 flex justify-center">
          <SearchInput 
            onSearch={handleSearch} 
            placeholder="Search superheroes or movies..."
          />
        </div>

        {/* Display Cards Continuously */}
        <section className="mb-8">
          <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6">
            {combinedData.map((item) => 
              item.image ? (  // Check if it's a superhero
                <SuperheroCard key={item.id || `${item.name}-${item.alignment}`} superhero={item} />
              ) : (  // If no image, assume it's a movie
                <MovieCard key={item.imdbID || `${item.Title}-${item.Year}`} movie={item} />
              )
            )}
          </div>
        </section>

        {/* Pagination */}
        <div className="mt-8 flex justify-center">
          <Pagination
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={setCurrentPage}
          />
        </div>
      </div>
    </div>
  );
};

export default SuperheroSearchPage;
