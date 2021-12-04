# exedep

A simple command line tool to find out what dynamically linked libraries are used by your **Go** application built with **MSYS2 MinGW 64-bit**.

## Requirements

- **MSYS2 MinGW 64-bit**;
- **Visual Studio Community 2022** (or any other version) with developer tools for **VC++**.

## Installation

```bash
go get github.com/zergon321/exedep
```

## Usage

It has the following command line options:
- `exe` is a path to the executable file being analyzed for dynamic dependencies. Must be specified by the user;
- `dll` is a path to the directory with the dependencies required by **CGO** (`C:\msys64\mingw64\bin` by default;
- `dumpbin` is a path to the **dumpbin** executable (`C:\Program Files\Microsoft Visual Studio\2022\Community\VC\Tools\MSVC\14.30.30705\bin\Hostx64\x64\dumpbin.exe` by default).

Example:

```bash
exedep --exe=/c/Users/user/go/src/github.com/zergon321/reisen/examples/player/player.exe
```