# 📷 Image Dupes Finder

## 🚀 Summary of Project

Welcome to the **Image Dupes Finder** repository! This project helps you identify and report similar or duplicate images within a specified directory. The tool scans a root directory for images, computes their hashes, and identifies groups of similar images. Finally, it generates an HTML report for easy visualization and review.

## 🛠️ How to Use

### Prerequisites

- Install [Go](https://golang.org/dl/) (version 1.22.2 or later).

### Build

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/harperreed/image-dupes.git
   cd image-dupes
   ```

2. **Install Dependencies**:
   ```sh
   go mod tidy
   ```

3. **Build the Executable**:
   ```sh
   go build -o image-dupes
   ```

### Running the Tool

After building the executable, you can run the Image Dupes Finder using the following command:

```sh
./image-dupes -dir <root_directory>
```

Example:
```sh
./image-dupes -dir /path/to/images -output report.html
```

### Output

The tool generates an HTML report (`report.html` by default) that lists groups of similar images for easy review.

## 🧑‍💻 Tech Info

This project is built with Go and utilizes several external libraries for image processing and terminal display enhancements.

### Directory Structure:

```
image-dupes/
├── go.mod
├── go.sum
├── hash.go
├── main.go
├── progress.go
├── report.go
├── scanner.go
├── similarity.go
└── test.html
```

### Dependencies

- [github.com/briandowns/spinner](https://github.com/briandowns/spinner) for terminal spinner.
- [github.com/corona10/goimagehash](https://github.com/corona10/goimagehash) for image hashing.
- [github.com/fatih/color](https://github.com/fatih/color) for terminal color.
- [github.com/mattn/go-colorable](https://github.com/mattn/go-colorable) and [github.com/mattn/go-isatty](https://github.com/mattn/go-isatty) for cross-platform terminal compatibility.
- [github.com/nfnt/resize](https://github.com/nfnt/resize) for image resizing.
- [golang.org/x/sys](https://golang.org/x/sys) and [golang.org/x/term](https://golang.org/x/term) for system-specific APIs.

### Key Files

- **hash.go**: Contains functions to compute file and perceptual hashes for images.
- **main.go**: Entry point for the application, manages scanning, hashing, finding similar images, and generating the report.
- **progress.go**: Handles progress display in the terminal.
- **progress_test.go**: Contains the test suite for progress.go, ensuring correct functionality of the Progress struct and its methods.
- **report.go**: Contains logic to generate an HTML report from the found image groups.
- **scanner.go**: Recursively scans the directory for images.
- **similarity.go**: Implements algorithms to compare and group similar images.

### Running Tests

To run the test suite for progress.go, use the following command:

```sh
go test -v ./...
```

Thanks for exploring Image Dupes Finder! 🎉 If you encounter any issues or have suggestions, feel free to open an issue or contribute to the project. Happy coding!

## Running Tests

To run the test suite, use the following command in the project root directory:

```sh
go test ./...
```


This will run all tests in the project, including:
- Tests for the `hash.go` file
- Tests for the `scanner.go` file (new)
- Tests for the `progress.go` file

The test suite for `scanner.go` includes comprehensive tests for the `scanDirectoryRecursive` function, covering various scenarios such as:
- Scanning empty directories
- Scanning directories with no image files
- Scanning directories with only image files
- Scanning directories with mixed file types
- Scanning directories with subdirectories containing images
- Handling different image extensions (.jpg, .jpeg, .png)
- Error handling for inaccessible directories or files


This will run all tests in the project, including the newly added tests for the `hash.go` and `report.go` files.

### Test Coverage

The project now includes comprehensive test suites for various components:

- `progress_test.go`: Tests for the progress tracking functionality.
- `hash_test.go`: Tests for the image hashing functionality.
- `report_test.go`: Tests for the HTML report generation functionality.

These test suites cover various scenarios, including successful operations, error handling, and edge cases.

## Contributing

If you're contributing to the project, please ensure that you add or update tests as necessary to maintain code quality and reliability. The existing test suites provide good examples of how to structure and implement tests for new features or modifications.

