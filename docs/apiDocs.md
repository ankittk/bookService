# 📚 Book Service API Documentation

## Overview

The **Book Service** provides an API for searching, retrieving, and managing books. It allows users to search for books based on various filters, retrieve book details, and manage book resources.

This service will expose **gRPC APIs** to interact with books in the system.

## Core Concepts

- **Book**: Represents a book resource containing information like title, author, description, etc.
- **Search**: Allows searching for books based on various query parameters.
- **Filters**: Parameters that help refine search results (e.g., author, genre, language).
- **Pagination**: Helps control the volume of data returned in search results.

---

## API Endpoints

### 1. Search Books

- **Method**: `SearchBooks`
- **Description**: Searches for books based on a query string and optional filters.

#### Input Fields:
- `query` (string, required): The search query (e.g., book title, author, genre).
- `author` (string, optional): Filter by author.
- `genre` (string, optional): Filter by genre.
- `language` (string, optional): Filter by language.
- `page` (int, optional): The page number for pagination (default is 1).
- `perPage` (int, optional): The number of results per page (default is 10).

#### Output Fields:
- `books` (array of Book, required): A list of books matching the search criteria.
- `totalCount` (int, required): The total number of matching books.
- `currentPage` (int, required): The current page number.
- `totalPages` (int, required): The total number of pages available based on the `perPage` value.

---

### 2. Get Book Details

- **Method**: `GetBookDetails`
- **Description**: Retrieves detailed information about a specific book using its unique identifier.

#### Input Fields:
- `bookId` (string, required): The unique identifier of the book.

#### Output Fields:
- `book` (Book, required): The detailed information of the requested book.

---

### 3. Add a Book

- **Method**: `AddBook`
- **Description**: Adds a new book to the system.

#### Input Fields:
- `book` (Book, required): A Book object containing all the necessary information to create a new book. Fields may include:
	- `title` (string, required)
	- `author` (string, required)
	- `description` (string, optional)
	- `isbn` (string, optional)
	- `publishedDate` (string, optional)
	- `language` (string, optional)
	- `genre` (string, optional)

#### Output Fields:
- `status` (string, required): The status of the request (e.g., "success", "error").
- `message` (string, required): A message describing the result (e.g., "Book added successfully").

---

### 4. Update Book

- **Method**: `UpdateBook`
- **Description**: Updates the details of an existing book.

#### Input Fields:
- `bookId` (string, required): The unique identifier of the book.
- `book` (Book, required): A Book object containing updated information for the book.

#### Output Fields:
- `status` (string, required): The status of the request (e.g., "success", "error").
- `message` (string, required): A message describing the result (e.g., "Book updated successfully").

---

### 5. Delete Book

- **Method**: `DeleteBook`
- **Description**: Deletes a book from the system using its unique identifier.

#### Input Fields:
- `bookId` (string, required): The unique identifier of the book to be deleted.

#### Output Fields:
- `status` (string, required): The status of the request (e.g., "success", "error").
- `message` (string, required): A message describing the result (e.g., "Book deleted successfully").

---

## Third-Party Libraries

### 1. `google/api/annotations.proto` and `protoc-gen-validate`
These libraries are used for defining annotations and validation rules for gRPC APIs. They provide a way to specify HTTP mappings, validation constraints, and more for API methods.

- **google/api/annotations.proto**: Used to define HTTP mappings for gRPC methods (e.g., RESTful API mappings).
- **protoc-gen-validate**: Provides a set of validation rules for Proto files, such as checking the length of strings, ensuring fields are not empty, etc.

### 2. `grpc-gateway`
`grpc-gateway` is a plugin for the Protocol Buffers compiler that generates a reverse proxy server, allowing HTTP/1.1 clients to access gRPC services. It translates HTTP RESTful calls into gRPC calls, enabling seamless integration between REST and gRPC APIs.

---

## Generating Code

To generate the Go code for your service, including the necessary gRPC client and server code, as well as the HTTP proxy, you can use the `protoc` compiler with the appropriate plugins. Below is the `protoc` command that you should run:

### Steps to Generate Code

1. **Ensure your `proto` and `third_party` folders are set up correctly:**
	- Your `proto` files should be in the `root/proto` directory.
	- Your third-party libraries (`google/api/annotations.proto` and `protoc-gen-validate/validate.proto`) should be in the `root/proto/third_party` directory.

2. **Run the following command to generate the necessary code:**

   ```bash
	mkdir -p ./proto/gen
   
	protoc -I ./proto \
	-I ./proto/third_party \
	--go_out=./proto/gen --go-grpc_out=./proto \    
	--grpc-gateway_out=logtostderr=true:./proto \    
	--validate_out=lang=go:./proto \    
	--go_opt=paths=source_relative \
	./proto/*.proto
   ```
   
   - `-I .`: Specifies the current directory as an import path.
   - `-I ./proto/third_party`: Specifies the path to the third-party libraries.
   - `--go_out=plugins=grpc:./proto`: Generates Go code with gRPC support.
   - `--go-grpc_out=./proto`: Generates gRPC-specific code.
   - `--grpc-gateway_out=logtostderr=true:./proto`: Generates the HTTP reverse proxy code.
   - `--validate_out=lang=go:./proto`: Generates validation code for the Proto files.
   - `./proto/*.proto`: Specifies the Proto files to be processed.
   
3. **Check the generated code:**
   After running the command, you should see the generated Go files in the `./proto/gen` directory. These files will include the gRPC service definitions, HTTP proxy code, and validation code.
