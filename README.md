# CSV to Excel Converter

This is a Go application that converts CSV files to Excel (.xlsx) format with a graphical user interface (GUI). It supports flexible parsing options and ensures all data is saved as text to avoid unwanted number formatting. The first row of the CSV is treated as column headers.

## Features
- **Flexible CSV Parsing**:
  - **Delimiter Mode**: Supports multiple delimiters (`;`, `\t`, `,`, space, or custom) with a checkbox-based selection in the GUI.
  - **Fixed Width Mode**: Allows specifying field widths (e.g., `10,20,15`) for fixed-width CSV files (basic support).
- **Graphical Interface**: Drag-and-drop file input or file picker via a button.
- **Text Preservation**: All data is preserved as text in the Excel output (no automatic number conversion).
- **File Naming**: Automatically generates unique filenames (e.g., `file(1).xlsx`) if the output file already exists.
- **Cross-Platform**: Binaries available for Linux and Windows (7, 10, 11). macOS support is excluded due to build issues.
- **Versioned Binaries**: Compiled binaries include the version in their names (e.g., `csv-to-excel-windows-amd64-v1.0.8.exe`).

## Usage
1. **Download**: Grab the latest binary from the [Releases](https://github.com/Fgeeha/csv-to-excel-the-go/releases) page.
   - Example: `csv-to-excel-windows-amd64-v1.0.8.exe` for Windows.
   - Example: `csv-to-excel-linux-amd64-v1.0.8` for Linux.
2. **Run**:
   - **Windows**: Double-click the `.exe` file or run `csv-to-excel-windows-amd64-v<version>.exe` (no console window will appear).
   - **Linux**: Execute `./csv-to-excel-linux-amd64-v<version>` in a terminal.
3. **Convert a File**:
   - Drag a CSV file into the window or click "Выбрать CSV-файл" to select one.
   - Choose parsing method:
     - **"С разделителями"**: Select one or more delimiters (e.g., semicolon, tab) via checkboxes.
     - **"Фиксированная ширина"**: Enter field widths in the text box (e.g., `10,20,15`).
   - The converted `.xlsx` file will be saved in the same directory as the input CSV, with a unique name if needed (e.g., `input(1).xlsx`).

## Example CSV Files
### Delimiter Mode

```text
"Name";"Age";"Phone"
"John Doe";"25";"1234567890"
```

- Save with UTF-8 encoding and use semicolon (`;`) as the delimiter.

### Fixed Width Mode (basic)
```text
Name      Age Phone
John Doe  25  1234567890
```

- Enter widths like `10,3,10` (not fully implemented yet).

## Building from Source
1. **Clone the repository**:
   ```bash
   git clone https://github.com/Fgeeha/csv-to-excel-the-go.git
   cd csv-to-excel-the-go
   ```
2. **Install dependencies**:
  ```bash
  go mod download
  ```
3. **Build**:
  - **Windows** (hides console):
    ```bash
    go build -ldflags "-H=windowsgui -X main.version=1.0.8" -o csv-to-excel.exe main.go
    ```
  - **Linux**:
    ```bash
    go build -ldflags "-X main.version=1.0.8" -o csv-to-excel main.go
    ```
4. **Run**:
   - `./csv-to-excel` (Linux)
   or
   - `csv-to-excel.exe` (Windows).
