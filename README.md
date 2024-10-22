# Raccoon's stash

**Raccoon's stash** is a self-hostable service which provides a simple pastebin and file sharing.

## Building

```shell
make        # Build everything
make cli    # Build the CLI
make server # Build the server
make sql    # Regenerate code from SQL
```

### Dependencies

- `sqlc` is used to generate code from SQL.

## Usage

```
raccoonstash-server :8080       # Start the server on port 8080
raccoonstash-cli -regiter rand  # Register a random token
```

### Options

Raccoon's stash provides 2 binaries:

- `raccoonstash-server`: Hosts the service.
- `raccoonstash-cli`: CLI for managing the service.

#### Global options

- `-datadir` `string`: Path to directory with data (default `data/`)

#### Server options

- `-dev`: Enable development mode (use files from source tree instead of EmbedFS)
- `-listen` `string`: Specify the address to listen on (default ":8080")

#### CLI options

- `-register` `string`: Register a token. Use `rand` to generate a random token
- `-unregister` `string`: Unregister a token. Use `all` to unregister all tokens
- `-tokens`: List all tokens
- `-stats`: Show server statistics
