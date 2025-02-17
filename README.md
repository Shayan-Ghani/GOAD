# GOAD
Go-based Interface-Driven Task management with Item and Tag Microservices supporting multiple backends in Repository layer. Currently raw SQL is implemented as database backend. You can use the CLI or call the API with a client.  

## Table of Contents
- [Features](#features)
- [Future Plans](#future-plans)
- [Project Structure](#project-structure)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)

## Features
- Command Line Interface
- Full CRUD Microservices in Repository/Service Structure
- REST (http/json) in delivery and microservices communication
- MySQL database backend
- Status tracking (done/pending)
- Timestamp tracking (created/modified)
- Due-date for items

## Future Plans
- Multiple-database backends
- OCR
- Circuit breaker
- Scalability and Reliability testing
- Unit Tests
- Grpc
- Migrate to Event-Driven for services

## Installation

1. Install dependencies:
```bash
go mod download
```
2. run services (you can optionally, change the service address in config/):
```bash
go run cmd/service/item.go
go run cmd/service/tag/tag.go
```

3. build the cli app :
```bash
go build -o cli main.go
```


# Project Structure
```
GOAD/
├── cmd/                    # Executables
│   ├── cli/               # CLI client application
│   └── service/           # Microservice executables

├── config/                # Configuration constants
├── docs/                  # Documentation
├── internal/              # Private application code
│   ├── db/               # Database management
│   │   └── migrations/   # SQL migration files
│   ├── delivery/         # Delivery layer (CLI & HTTP)
│   │   ├── cli/         # CLI implementation
│   │   │   └── requesthandler/
│   │   └── command/     # Command parsing
│   ├── model/           # Domain models
│   ├── repository/      # Data access layer
│   │   └── sql/        # SQL implementation
│   └── service/         # Business logic
│       ├── item/        # Item service
│       └── tag/         # Tag service
└── pkg/                  # Public packages
    ├── formatter/       # Formatting utilities
    ├── request/        # Request/payload definitions
    │   ├── item/
    │   └── tag/
    ├── response/       # Response handling
    └── validation/     # Input validation
```

## Usage 
+ see the help page for options:
```bash
./cli --help
```
<details>
<summary>As of now here are the options:</summary>
<pre>
Available commands:
item : ...
  --help    	see help for flags and options
    Usage: see help for flags and options

  add       	Add a new item
    Usage: item add -n <name> -d <description> [-t tag1,tag2] [-due-date <date string> (e.g '2025-03-05 15:05:10') ] 

  view      	View an item(s), also use -t <tag-name> instead of -i(single)/--all to see items filtered by that tag name.
    Usage: item view [-i <id>] [--done=true] [--all=true] [-t <items-with-these-tags,tag2>] [--format=json/table]

  delete    	Delete an item or its tags with --del-tags
    Usage: item delete -i <id> [-t <tags-to-delete> ] [--del-tags=true]

  update    	Update an item
    Usage: item update -i <id> [-n <name>] [-d <description>] [-t <tag1,tag2>] [-due-date <date string> (e.g '2025-03-05 15:05:10') ]

  done      	update item status to done from pending
    Usage: item done -i <id>


tag : ...
  view      	View tags
    Usage: tag view --all=true

  delete    	Delete a tag or item tags, if item id not provided all refrences of the tag will be removed from items.
    Usage: tag delete -n <name> [-item-id <id> [-t <tags,to,remove> / -all (remove all tags for item)] ]

</pre>
</details>

## Examples

1. getting all the items:
```bash
./cli item view --all=true
```
<details>
<summary> Output:</summary>
<pre>
The Game Begins.
ID   Name             Description               Status    Due_Date              Tags             Created_At
--   ----             -----------               ------    ----                  ------           ------
2    fix table view   fix table view response   Done      Not Set               No tags          2025-02-04 16:41:53
3    fix json view    fix json view response    Pending  2025-05-02 18:32:02    fix, chore       2025-02-04 16:44:10
4    json marshal     adding json response      Pending   Not Set               gocasts, chore   2025-02-04 20:18:04
</pre>
</details>