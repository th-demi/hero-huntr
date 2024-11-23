#!/bin/bash
# This script generates a directory tree structure followed by all Go code files into one plain text file.

# Specify the output plain text file
output_file="project_code.txt"

# Clear the output file if it exists
> "$output_file"

# Generate the directory structure in a tree-like format and append it to the output file
echo "Project Directory Structure:" >> "$output_file"
tree -F --noreport >> "$output_file"

# Add a separator line between the directory structure and code
echo -e "\n\n==== Code Files ====\n\n" >> "$output_file"

# Loop through all the code files (e.g., Go files) in the directory
find . -type f -name "*.go" | while read -r file; do
    # Add the directory structure (file path) as a header in the text file
    echo "==== File: $file ====" >> "$output_file"
    
    # Append the content of the Go file to the output file
    cat "$file" >> "$output_file"
    
    # Add a newline for separation between files
    echo -e "\n\n" >> "$output_file"
done

echo "Conversion complete! All files have been appended into '$output_file'."
