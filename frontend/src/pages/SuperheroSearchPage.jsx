import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Pagination } from "@/components/ui/pagination";

// Truncate text helper function
const truncate = (str, length) => str.length > length ? `${str.substring(0, length)}...` : str;

const SearchInput = ({ onSearch, placeholder }) => {
  const [query, setQuery] = useState('');

  const handleSearch = () => {
    if (query.trim()) {
      onSearch(query.trim());
    }
  };

  return (
    <div className="flex w-full max-w-md space-x-2">
      <Input 
        type="text" 
        placeholder={placeholder}
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
        className="flex-grow"
      />
      <Button 
        onClick={handleSearch} 
        className="bg-blue-600 hover:bg-blue-700"
        disabled={!query.trim()} // Disable the button if the query is empty
      >
        <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        Search
      </Button>
    </div>
  );
};

const SuperheroCard = ({ superhero }) => {
  const heroKey = superhero._id?.$oid || `${superhero.name}-${superhero.alignment}`;
  
  return (
    <Card key={heroKey} className="w-full max-w-xs transition-all hover:scale-105 hover:shadow-xl bg-red-100">
      <CardHeader className="p-4 pb-0">
        <img 
          src={superhero.image}
          alt={superhero.name || 'Unknown Hero'}
          className="w-full h-64 object-cover rounded-t-lg"
        />
      </CardHeader>
      <CardContent className="p-4">
        <CardTitle className="text-xl font-bold mb-2 text-black">
          {truncate(superhero.name || 'Unknown Hero', 20)}
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

const MovieCard = ({ movie }) => {
  const movieKey = movie.Title ? movie.Title : `${movie.Title || 'Unknown Movie'}-${movie.Year || 'N/A'}`;

  return (
    <Card key={movieKey} className="w-full max-w-xs transition-all hover:scale-105 hover:shadow-xl bg-blue-100">
      <CardHeader className="p-4 pb-0">
        <img 
          src={movie.Poster || '/default-poster.jpg'}
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

const SuperheroSearchPage = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [searchResults, setSearchResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [currentQuery, setCurrentQuery] = useState('');
  const ITEMS_PER_PAGE = 12;

  // Cache for search results
  const [cachedResults, setCachedResults] = useState({});

  useEffect(() => {
    // Reset search results when the query changes
    setSearchResults([]);
    setTotalPages(1);
    setCurrentPage(1);
  }, [currentQuery]);

  const fetchResults = async (query, page) => {
    // Check if data is already in cache
    if (cachedResults[query] && cachedResults[query][page]) {
      // If cached, use the data directly
      setSearchResults(cachedResults[query][page].results);
      setTotalPages(cachedResults[query][page].totalPages);
      return;
    }

    try {
      setLoading(true);

      const response = await fetch(
        `http://localhost:8080/api/search?query=${query}&page=${page}&limit=${ITEMS_PER_PAGE}`
      );
      
      if (!response.ok) {
        throw new Error("Failed to fetch search results");
      }

      const data = await response.json();

      // Combine the paginated results
      const combined = [
        ...(data.superheroes || []).map(hero => ({ ...hero, type: 'superhero' })),
        ...(data.movies || []).map(movie => ({ ...movie, type: 'movie' }))
      ];

      // Update the state with new results
      setSearchResults(combined);
      setTotalPages(data.totalPages);

      // Update the cache with new results for the given query and page
      setCachedResults((prevCache) => ({
        ...prevCache,
        [query]: {
          ...prevCache[query],
          [page]: {
            results: combined,
            totalPages: data.totalPages
          }
        }
      }));

    } catch (err) {
      console.error("Error fetching data:", err);
      setSearchResults([]);
      setTotalPages(1);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async (query) => {
    setCurrentQuery(query);
    setCurrentPage(1); // Reset to first page
    await fetchResults(query, 1);
  };

  const handlePageChange = async (page) => {
    setCurrentPage(page);
    await fetchResults(currentQuery, page);
  };

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-6xl mx-auto">
        <div className="mb-8 flex justify-center items-center">
          <h1 className="text-4xl font-bold text-gray-800">Hero Huntr</h1>
        </div>

        <div className="mb-8 flex justify-center">
          <SearchInput
            onSearch={handleSearch}
            placeholder="Search superheroes or movies..."
          />
        </div>

        {loading ? (
          <div className="flex justify-center">Loading...</div>
        ) : (
          <>
            <section className="mb-8">
              <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6">
                {searchResults.map((item, index) => (
                  item.type === 'superhero' ? (
                    <SuperheroCard 
                      key={`${item._id?.$oid || item.name}-${index}`} 
                      superhero={item} 
                    />
                  ) : (
                    <MovieCard 
                      key={`${item.imdbID || item.Title}-${index}`} 
                      movie={item} 
                    />
                  )
                ))}
              </div>
            </section>

            {totalPages > 1 && searchResults.length > 0 && (
              <div className="mt-8 flex justify-center">
                <Pagination
                  currentPage={currentPage}
                  totalPages={totalPages}
                  onPageChange={handlePageChange}
                />
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default SuperheroSearchPage;
