# CSV to Excel Converter

This is a simple Go application that converts a CSV file (semicolon-separated, UTF-8 encoded) to an Excel file (.xlsx). The first row of the CSV is treated as column headers, and all data is preserved as text to prevent number formatting issues.

## Features
- Converts CSV files with semicolon (`;`) delimiters to Excel.
- Preserves all data as text (no automatic number conversion).
- Cross-platform: binaries available for Linux, macOS, and Windows (7, 10, 11).

## Usage
1. Download the appropriate binary from the [Releases](https://github.com/username/csv-to-excel/releases) page.
2. Run the program:
   ```bash
   ./csv-to-excel