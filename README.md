# GitRel

GitRel is a command-line tool designed to help manage git release branches efficiently. It provides commands to list, create, and checkout release branches, as well as increment version numbers following semantic versioning.

## Configuration

GitRel can be configured using a `.gitrelrc` file. The configuration file can be placed in the current directory, any parent directory, or in the user's home directory. The following options are available:

- `alwaysFetch=true|false`: If set to true, the `--fetch` flag will be presumed for all commands that accept it.
- `remote=<git remote name>`: Specifies the git remote name to use. Defaults to `origin` if not set.

## Installation

To install GitRel, clone the repository and use `go install` to build the tool:

```bash
git clone <repository-url>
cd <repository-directory>
go install .
```

Replace `<repository-url>` and `<repository-directory>` with the actual URL and directory name of your repository.

## Usage

GitRel provides several commands to manage your release branches:

- **list**: List current release branches.
- **new**: Create a new release branch.
  - **<version>**: Create a new release branch with the specified version.
  - **major**: Increment the major version of the latest release.
  - **minor**: Increment the minor version of the latest release.
  - **patch**: Increment the patch version of the latest release.
- **status**: Show the current version and the 5 most recent versions.
- **checkout**: Checkout a release branch.
  - **<version>**: Checkout the release branch matching the specified version prefix.
  - **latest**: Checkout the latest release branch.

### Examples

1. **List release branches**:
   ```bash
   gitrel list
   ```

   To fetch from remote before listing:
   ```bash
   gitrel list --fetch
   ```

2. **Create a new release branch with a specific version**:
   ```bash
   gitrel new 1.0.0
   ```

3. **Increment the major version and create a new release branch**:
   ```bash
   gitrel new major
   ```

4. **Show the current version and recent versions**:
   ```bash
   gitrel status
   ```

   To fetch from remote before showing status:
   ```bash
   gitrel status --fetch
   ```

5. **Checkout a specific release branch**:
   ```bash
   gitrel checkout 1.0.0
   ```

6. **Checkout the latest release branch**:
   ```bash
   gitrel checkout latest
   ```

For more detailed information on each command, you can use the `--help` flag with any command, e.g., `gitrel list --help`.