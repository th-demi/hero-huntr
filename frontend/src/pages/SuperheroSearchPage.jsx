import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Pagination } from "@/components/ui/pagination";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Slider } from "@/components/ui/Slider";
 // Assuming this is a custom slider component
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select";

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
  const [filters, setFilters] = useState({
    power: [0, 100],
    alignment: '', // 'good', 'neutral', or 'bad'
  });
  const [showAdvancedSearch, setShowAdvancedSearch] = useState(false);
  const ITEMS_PER_PAGE = 12;

  // Cache for search results
  const [cachedResults, setCachedResults] = useState({});

  useEffect(() => {
    // Reset search results when the query or filters change
    setSearchResults([]);
    setTotalPages(1);
    setCurrentPage(1);
  }, [currentQuery, filters]);

  const fetchResults = async (query, page, filters) => {
    if (!query.trim()) {
      // If no query is provided, do not fetch results
      return;
    }
  
    // Destructure filters with default values
    const { power = [0, 100], alignment = '' } = filters;
  
    // Construct query parameters manually
    const params = new URLSearchParams({
      query: query,
      page: page,
      limit: ITEMS_PER_PAGE
    });
  
    // Only add powerMin and powerMax if power range is different from the default [0, 100]
    if (power[0] !== 0 || power[1] !== 100) {
      params.set('powerMin', power[0]);
      params.set('powerMax', power[1]);
    }
  
    // Only add alignment if it's not empty
    if (alignment) {
      params.set('alignment', alignment);
    }
  
    const url = `http://localhost:8080/api/search?${params.toString()}`;
  
    console.log('Search URL:', url);  // Debugging: log the URL to verify it's correct
  
    // Check if data is already in cache
    if (cachedResults[query] && cachedResults[query][page] && cachedResults[query][page].filters === filters) {
      setSearchResults(cachedResults[query][page].results);
      setTotalPages(cachedResults[query][page].totalPages);
      return;
    }
  
    try {
      setLoading(true);
  
      const response = await fetch(url);
  
      if (!response.ok) {
        throw new Error("Failed to fetch search results");
      }
  
      const data = await response.json();
  
      // Combine the paginated results
      const combined = [
        ...(data.superheroes || []).map(hero => ({ ...hero, type: 'superhero' })),
        ...(data.movies || []).map(movie => ({ ...movie, type: 'movie' }))
      ];
  
      setSearchResults(combined);
      setTotalPages(data.totalPages);
  
      setCachedResults((prevCache) => ({
        ...prevCache,
        [query]: {
          ...prevCache[query],
          [page]: {
            results: combined,
            totalPages: data.totalPages,
            filters
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
    setCurrentPage(1);
    await fetchResults(query, 1, filters);
  };

  const handlePageChange = async (page) => {
    setCurrentPage(page);
    await fetchResults(currentQuery, page, filters);
  };

  const handleAdvancedSearch = (newFilters) => {
    setFilters(newFilters);
    setCurrentPage(1);
    setShowAdvancedSearch(false);
    fetchResults(currentQuery, 1, newFilters);
  };

  // Handle power range changes
  const handlePowerChange = (value) => {
    setFilters(prevFilters => ({
      ...prevFilters,
      power: value
    }));
  };

  // Handle alignment selection changes
  const handleAlignmentChange = (alignment) => {
    setFilters(prevFilters => ({
      ...prevFilters,
      alignment: alignment
    }));
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

        {/* Advanced Search Dialog */}
        <div className="mb-8 flex justify-center">
          <Dialog open={showAdvancedSearch} onOpenChange={setShowAdvancedSearch}>
            <DialogTrigger asChild>
              <Button variant="outline" className="inline-flex items-center space-x-2 bg-black text-white" onClick={() => setShowAdvancedSearch(true)}>
                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
                </svg>
                <span>Advanced Search</span>
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl bg-white">
              <DialogHeader>
                <DialogTitle className="text-lg font-semibold text-black">Advanced Search Filters</DialogTitle>
              </DialogHeader>
              <div id="dialog-description" className="text-sm text-gray-600">
                Use the sliders and options to filter your search results based on power level and alignment.
              </div>
              <div className="space-y-6">
                {/* Power Slider */}
                <div className="flex flex-col space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-gray-900">Power</span>
                    <span className="text-sm text-gray-500">{filters.power[0]} - {filters.power[1]}</span>
                  </div>
                  <Slider
                    min={0}
                    max={100}
                    step={1}
                    value={filters.power}
                    onValueChange={handlePowerChange}
                    className="w-full"
                  />
                </div>

                {/* Alignment Select */}
                <div className="flex flex-col space-y-2">
                  <label className="text-sm font-medium text-gray-900">Alignment</label>
                  <Select 
                    value={filters.alignment}
                    onValueChange={handleAlignmentChange}
                  >
                    <SelectTrigger className="w-full bg-white text-black border-black">
                      <SelectValue placeholder="Select alignment" />
                    </SelectTrigger>
                    <SelectContent className="bg-white text-black border border-gray-300">
                      {['bad', 'neutral', 'good'].map(align => (
                        <SelectItem key={align} value={align} className="capitalize">
                          {align}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="flex gap-4 mt-6">
                <Button 
                  variant="outline" 
                  onClick={() => setFilters({ power: [0, 100], alignment: '' })}
                  className="w-full hover:bg-gray-100 bg-white text-black"
                >
                  Reset Filters
                </Button>
                <Button 
                  onClick={() => handleAdvancedSearch(filters)} 
                  className="w-full bg-primary hover:bg-primary/90"
                >
                  Apply Filters
                </Button>
              </div>
            </DialogContent>
          </Dialog>
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
