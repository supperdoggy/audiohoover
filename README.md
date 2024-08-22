# AudioHoover

**AudioHoover** is a simple Go script designed to collect all audio files from a directory named `playlists` and move them into a single destination directory, ensuring no duplicates are transferred.

## Features

- Recursively collects all audio files from subdirectories within `playlists`.
- Moves audio files to a specified destination directory.
- Prevents duplicates by checking if the file already exists in the destination.
