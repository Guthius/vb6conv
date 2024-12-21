# VB6 to .NET C# Project Converter

This tool converts VB6 VBP projects into .NET C# projects. Below is a guide on how to use the tool and its available flags.

## Build Instructions

### Install Go

   Ensure that [Go](https://golang.org/dl/) is installed on your system. The minimum required version of Go is 1.23.

   You can verify the installation by running:
   ```bash
   go version
   ```

### Clone the Repository

   To build the application and generate the executable, follow these steps:

   ```bash
   git clone https://github.com/guthius/vb6conv
   cd vb6conv
   ```

### Build the Application

   ```bash
   go build
   ```

## Usage

```bash
vb6conv --project <path_to_project_file> --output <output_directory> [--namespace <namespace>]
```

## Flags

### Required Flags

- `-p, --project`
  - **Description**: Specifies the path to the VB6 project file to be converted.
  - **Usage**: Provide the full path to the `.vbp` file.
  - **Example**:
    ```bash
    --project "C:/Projects/MyVB6Project.vbp"
    ```

- `-o, --output`
  - **Description**: Specifies the output directory where the converted C# project will be saved.
  - **Usage**: Provide the full path to the desired output directory.
  - **Example**:
    ```bash
    --output "C:/ConvertedProjects/MyCSharpProject"
    ```

### Optional Flags

- `-n, --namespace`
  - **Description**: Specifies the namespace for the converted C# project. If not provided, a default namespace will be used.
  - **Usage**: Provide a valid namespace name.
  - **Example**:
    ```bash
    --namespace "MyNamespace"
    ```

## Examples

### Minimal Example
Convert a VB6 project without specifying a namespace:

```bash
vb6conv --project "C:/Projects/MyVB6Project.vbp" --output "C:/ConvertedProjects/MyCSharpProject"
```

## Notes
- Both the `--project` and `--output` flags are **required** for the tool to run.
- If the `--namespace` flag is not provided, the tool will use the project name as the root namespace for the converted project.
- Ensure that the provided paths are valid and accessible to avoid errors.

## Help
For help or more information, use the `--help` flag:

```bash
vb6conv --help
```

