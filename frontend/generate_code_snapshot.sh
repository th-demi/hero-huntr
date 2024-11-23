#!/bin/bash
# This script generates a directory tree structure followed by all React component files into one plain text file.

# Specify the output plain text file
output_file="project_code.txt"

# Clear the output file if it exists
> "$output_file"

# Generate the directory structure in a tree-like format and append it to the output file, excluding node_modules
echo "Project Directory Structure (excluding node_modules):" >> "$output_file"
tree -F --noreport --prune -I 'node_modules' >> "$output_file"

# Add a separator line between the directory structure and code
echo -e "\n\n==== Code Files ====\n\n" >> "$output_file"

# Loop through all the React-related code files (e.g., .js, .jsx, .ts, .tsx)
find . -type f \( -name "*.js" -o -name "*.jsx" -o -name "*.ts" -o -name "*.tsx" \) -not -path "./node_modules/*" | while read -r file; do
    # Add the directory structure (file path) as a header in the text file
    echo "==== File: $file ====" >> "$output_file"
    
    # Append the content of the React file to the output file
    cat "$file" >> "$output_file"
    
    # Add a newline for separation between files
    echo -e "\n\n" >> "$output_file"
done

echo "Conversion complete! All files have been appended into '$output_file'."
