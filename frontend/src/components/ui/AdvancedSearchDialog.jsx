import React, { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Slider } from "./Slider";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "./tabs";

const languages = [
  "Afrikaans", "Albanian", "Amharic", "Arabic", "Armenian", "Basque", "Bengali", 
  "Bosnian", "Bulgarian", "Catalan", "Cebuano", "Chinese", "Croatian", "Czech", 
  "Danish", "Dutch", "English", "Estonian", "Filipino", "Finnish", "French", 
  "Georgian", "German", "Greek", "Gujarati", "Haitian Creole", "Hebrew", "Hindi", 
  "Hungarian", "Icelandic", "Igbo", "Indonesian", "Irish", "Italian", "Japanese", 
  "Javanese", "Kannada", "Kazakh", "Khmer", "Korean", "Kurdish", "Kyrgyz", 
  "Lao", "Latvian", "Lithuanian", "Macedonian", "Malay", "Malayalam", "Maltese", 
  "Marathi", "Mongolian", "Nepali", "Norwegian", "Pashto", "Persian", "Polish", 
  "Portuguese", "Punjabi", "Romanian", "Russian", "Serbian", "Sesotho", "Shona", 
  "Sindhi", "Sinhala", "Slovak", "Slovenian", "Somali", "Spanish", "Sundanese", 
  "Swahili", "Swedish", "Tagalog", "Tamil", "Telugu", "Thai", "Turkish", "Ukrainian", 
  "Urdu", "Uzbek", "Vietnamese", "Welsh", "Xhosa", "Yiddish", "Yoruba", "Zulu"
];

const initialSuperheroFilters = {
  intelligence: [0, 100],
  strength: [0, 100],
  speed: [0, 100],
  durability: [0, 100],
  power: [0, 100],
  combat: [0, 100],
  alignment: '',
  gender: ''
};

const initialMovieFilters = {
  type: '',
  imdbRating: [0, 10],
  language: '',
  runtime: [0, 300],
  year: [1950, 2025]
};

const RangeSlider = ({ min, max, step = 1, value, onChange, label }) => (
  <div className="flex flex-col space-y-2">
    <div className="flex items-center justify-between">
      <span className="text-sm font-medium text-gray-900 dark:text-gray-100">{label}</span>
      <span className="text-sm text-gray-500 dark:text-gray-400">
        {value[0]} - {value[1]}
      </span>
    </div>
    <Slider
      min={min}
      max={max}
      step={step}
      value={value}
      onValueChange={onChange}
      className="w-full"
      defaultValue={[min, max]}
    />
  </div>
);

const AdvancedSearchDialog = ({ onAdvancedSearch, onClose }) => {
  const [activeTab, setActiveTab] = useState("superheroes");
  const [superheroFilters, setSuperheroFilters] = useState(initialSuperheroFilters);
  const [movieFilters, setMovieFilters] = useState(initialMovieFilters);
  const [open, setOpen] = useState(true); // Add state to manage dialog visibility

  const handleSearch = () => {
    onAdvancedSearch({
      type: activeTab,
      filters: activeTab === "superheroes" ? superheroFilters : movieFilters
    });
    setOpen(false); // Close the dialog after search
    onClose(); // Call the provided onClose prop, if any
  };

  const handleReset = () => {
    if (activeTab === "superheroes") {
      setSuperheroFilters(initialSuperheroFilters);
    } else {
      setMovieFilters(initialMovieFilters);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}> {/* Use open and onOpenChange to control dialog */}
      <DialogTrigger asChild>
        <Button variant="outline" className="inline-flex items-center space-x-2 bg-black text-white">
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

        <Tabs defaultValue="superheroes" className="w-full" onValueChange={setActiveTab}>
          <TabsList className="grid w-full grid-cols-2 mb-4">
            <TabsTrigger value="superheroes" className="px-4 py-2 bg-white text-black border border-black rounded-lg hover:bg-gray-100 hover:text-black focus:outline-none focus:ring-2 focus:ring-primary">Superheroes</TabsTrigger>
            <TabsTrigger value="movies" className="px-4 py-2 bg-white text-black border border-black rounded-lg hover:bg-gray-100 hover:text-black focus:outline-none focus:ring-2 focus:ring-primary">Movies & TV Shows</TabsTrigger>
          </TabsList>

          <TabsContent value="superheroes" className="space-y-6">
            <div className="space-y-4">
              <RangeSlider
                label="Intelligence"
                min={0}
                max={100}
                value={superheroFilters.intelligence}
                onChange={(value) => setSuperheroFilters({...superheroFilters, intelligence: value})}
              />
              <RangeSlider
                label="Strength"
                min={0}
                max={100}
                value={superheroFilters.strength}
                onChange={(value) => setSuperheroFilters({...superheroFilters, strength: value})}
              />
              <RangeSlider
                label="Speed"
                min={0}
                max={100}
                value={superheroFilters.speed}
                onChange={(value) => setSuperheroFilters({...superheroFilters, speed: value})}
              />
              <RangeSlider
                label="Durability"
                min={0}
                max={100}
                value={superheroFilters.durability}
                onChange={(value) => setSuperheroFilters({...superheroFilters, durability: value})}
              />
              <RangeSlider
                label="Power"
                min={0}
                max={100}
                value={superheroFilters.power}
                onChange={(value) => setSuperheroFilters({...superheroFilters, power: value})}
              />
              <RangeSlider
                label="Combat"
                min={0}
                max={100}
                value={superheroFilters.combat}
                onChange={(value) => setSuperheroFilters({...superheroFilters, combat: value})}
              />

              {/* Alignment and Gender selects */}
              <div className="grid grid-cols-2 gap-4">
                <div className="flex flex-col space-y-2">
                  <label className="text-sm font-medium text-gray-900 dark:text-gray-100">Alignment</label>
                  <Select 
                    value={superheroFilters.alignment}
                    onValueChange={(value) => setSuperheroFilters({...superheroFilters, alignment: value})}
                  >
                    <SelectTrigger className="w-full bg-white text-black border-black">
                      <SelectValue placeholder="Select alignment" />
                    </SelectTrigger>
                    <SelectContent className="bg-white text-black border border-gray-300 dark:bg-white dark:text-black dark:border-gray-300">
                      {['bad', 'neutral', 'good'].map(align => (
                        <SelectItem key={align} value={align} className="capitalize">
                          {align}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="flex flex-col space-y-2">
                  <label className="text-sm font-medium text-gray-900 dark:text-gray-100">Gender</label>
                  <Select
                    value={superheroFilters.gender}
                    onValueChange={(value) => setSuperheroFilters({...superheroFilters, gender: value})}
                  >
                    <SelectTrigger className="w-full bg-white text-black border-black">
                      <SelectValue placeholder="Select gender" />
                    </SelectTrigger>
                    <SelectContent className="bg-white text-black border border-gray-300 dark:bg-white dark:text-black dark:border-gray-300">
                      {['Female', 'Male'].map(gender => (
                        <SelectItem key={gender} value={gender}>
                          {gender}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </div>
          </TabsContent>

          <TabsContent value="movies" className="space-y-6">
            <div className="space-y-4">
              {/* IMDB Rating Slider */}
              <RangeSlider
                label="IMDB Rating"
                min={0}
                max={10}
                step={0.1}
                value={movieFilters.imdbRating}
                onChange={(value) => setMovieFilters({ ...movieFilters, imdbRating: value })}
              />

              {/* Runtime Slider */}
              <RangeSlider
                label="Runtime (minutes)"
                min={0}
                max={300}
                value={movieFilters.runtime}
                onChange={(value) => setMovieFilters({ ...movieFilters, runtime: value })}
              />

              {/* Year Slider */}
              <RangeSlider
                label="Year"
                min={1950}
                max={2025}
                value={movieFilters.year}
                onChange={(value) => setMovieFilters({ ...movieFilters, year: value })}
              />

              {/* Type Select */}
              <div className="flex flex-col space-y-2">
                <label className="text-sm font-medium text-gray-900 dark:text-gray-100">Type</label>
                <Select
                  value={movieFilters.type}
                  onValueChange={(value) => setMovieFilters({ ...movieFilters, type: value })}
                >
                  <SelectTrigger className="w-full bg-white text-black border-black">
                    <SelectValue placeholder="Select type" />
                  </SelectTrigger>
                  <SelectContent className="bg-white text-black border border-gray-300 dark:bg-white dark:text-black dark:border-gray-300">
                    {['series', 'movie'].map(type => (
                      <SelectItem key={type} value={type} className="capitalize">
                        {type}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              {/* Language Select */}
              <div className="flex flex-col space-y-2">
                <label className="text-sm font-medium text-gray-900 dark:text-gray-100">Language</label>
                <Select
                  value={movieFilters.language}
                  onValueChange={(value) => setMovieFilters({ ...movieFilters, language: value })}
                >
                  <SelectTrigger className="w-full bg-white text-black border-black">
                    <SelectValue placeholder="Select language" />
                  </SelectTrigger>
                  <SelectContent className="bg-white text-black border border-gray-300 dark:bg-white dark:text-black dark:border-gray-300 max-h-60 overflow-y-auto">
                    {languages.map(lang => (
                      <SelectItem key={lang} value={lang}>
                        {lang}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </TabsContent>
        </Tabs>

        <div className="flex gap-4 mt-6">
          <Button 
            variant="outline" 
            onClick={handleReset} 
            className="w-full hover:bg-gray-100 dark:hover:bg-gray-800 bg-white text-black border-black"
          >
            Reset Filters
          </Button>
          <Button 
            onClick={handleSearch} 
            className="w-full bg-primary hover:bg-primary/90"
          >
            Apply Filters
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default AdvancedSearchDialog;
