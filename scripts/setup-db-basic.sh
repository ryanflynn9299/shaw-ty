#!/bin/bash

TARGET="/Users/ryan.flynn/GolandProjects/URL Shortener/sql"
DB_FILE="/Users/ryan.flynn/GolandProjects/URL Shortener/sql/url_shortener_test.db"

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

echo "SQLite Database Setup Script"
echo "----------------------------"

# Check if sqlite3 command is available
if ! command_exists sqlite3; then
    echo "Error: sqlite3 command-line tool is not installed."
    echo "Please install it to continue (e.g., 'sudo apt-get install sqlite3' or 'brew install sqlite')."
    exit 1
fi

# Check if the database file already exists
if [ -f "$DB_FILE" ]; then
    echo "Database file '$DB_FILE' already exists."
    read -p "Do you want to overwrite it? (y/N): " choice
    case "$choice" in
        y|Y )
            echo "Removing existing database file..."
            rm "$DB_FILE"
            ;;
        * )
            echo "Exiting without creating a new database."
            exit 0
            ;;
    esac
fi

echo "Creating SQLite database file: $DB_FILE"

# Create the database
# The sqlite3 command will create the file if it doesn't exist.
sqlite3 "$TARGET/$DB_FILE"

echo "Created SQLite database file: $TARGET/$DB_FILE"
printf 'DB_TYPE="sqlite3"\n# The DSN is the primary connection string for SQLite\nDB_DSN="file:"$target"/sql/file.db?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000"\n\n# The following are not applicable for SQLite\nDB_HOST=\nDB_PORT=\nDB_USER=\nDB_PASS=\n'