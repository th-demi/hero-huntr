import React, { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Slider } from "./Slider";  // Assuming this is a custom slider component
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select";

// Initial filters object
const initialFilters = {
  power: [0, 100],
  alignment: ''
};

// RangeSlider component for handling power range
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
      onValueChange={onChange}  // The onChange callback for the slider
      className="w-full"
      defaultValue={[min, max]}
    />
  </div>
);

const AdvancedSearchDialog = ({ onAdvancedSearch, onClose }) => {
  const [filters, setFilters] = useState(initialFilters);
  const [open, setOpen] = useState(true);

  // Handle power range changes
  const handlePowerChange = (value) => {
    setFilters(prevFilters => ({
      ...prevFilters,
      power: value  // Update power range
    }));
  };

  // Handle alignment selection changes
  const handleAlignmentChange = (alignment) => {
    setFilters(prevFilters => ({
      ...prevFilters,
      alignment: alignment  // Update alignment
    }));
  };

  // Handle search submit action
  const handleSearch = () => {
    onAdvancedSearch({ filters });
    setOpen(false); // Close the dialog after search
    onClose(); // Call the provided onClose prop, if any
  };

  // Reset filters to initial values
  const handleReset = () => {
    setFilters(initialFilters);
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

        <div className="space-y-6">
          {/* Power Slider */}
          <RangeSlider
            label="Power"
            min={0}
            max={100}
            value={filters.power}
            onChange={handlePowerChange}  // Use the handlePowerChange function here
          />

          {/* Alignment Select */}
          <div className="flex flex-col space-y-2">
            <label className="text-sm font-medium text-gray-900 dark:text-gray-100">Alignment</label>
            <Select 
              value={filters.alignment}
              onValueChange={handleAlignmentChange}  // Use the handleAlignmentChange function here
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
        </div>

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
