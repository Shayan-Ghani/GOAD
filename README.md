# GOAD
A Go-based Task Management Interactive CLI implementing CRUD operations for managing todo items using MySQL and Flag.

## Table of Contents
- [Features](#features)
- [Future Plans](#future-plans)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)

## Features
- interactive cli
- Full CRUD operations for todo items
- Tag management system
- Status tracking (done/pending)
- Timestamp tracking (created/modified)
- MySQL database backend

## Future Plans
- REST (http/json) delivery
- Multiple-database backend
- Due-date for tasks
- OCR

## Installation

1. Install dependencies:
```bash
go mod download
```
2. build the app :
```bash
go build -o cli main.go
```

## Usage 
1. see the help page for options:
```bash
./cli --help
```
+ As of now here are the options:
```
Available commands:
item : ...
  view          View an item(s)
    Usage: item view -i <id> [--done=true] [--all=true] [-t <items-with-these-tags,tag2>]

  delete        Delete an item
    Usage: item delete -i <id> [-t <tags-to-delete> ] [--del-tags=true]

  update        Update an item
    Usage: item update -i <id> [-n <name>] [-d <description>] [-t <tag1,tag2>]

  done          update item status to done from pending
    Usage: item done -i <id>

  --help        see help for flags and options
    Usage: see help for flags and options

  add           Add a new item
    Usage: item add -n <name> -d <description> [-t tag1,tag2]


tag : ...
  view          View tags
    Usage: tag view --all=true

  delete        Delete a tag, all refrences of the tag will be removed from items.
    Usage: tag delete -n <name>

```

## Examples

1. getting all the items:
```bash
./cli item view --all=true

output:
The Game Begins.
ID   Name     Description    Status    Tags                       Created_At
--   ----     -----------    ------    ----                       ----
1    done@    done it now    Pending   gemeni, kjg                2025-01-09 14:19:09
2    Daaamn   done it then   Done      No tags                    2025-01-09 14:19:50
3    SOKA     done it ZHEN   Pending   No tags                    2025-01-09 19:53:02
4    SOKA     done it ZHEN   Pending   No tags                    2025-01-09 19:59:07
5    SOKA     done it ZHEN   Pending   gemeni, kjg, gjjg, DAAMN   2025-01-09 20:11:28
6    AAA      Tessing now    Pending   Unit                       2025-01-18 22:19:49

```