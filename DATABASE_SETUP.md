🚀 SQLite & Bun ORM: Database Setup Guide
This guide is for setting up a fresh SQLite database for your Go application that uses the Bun ORM. Follow these steps to initialize the database, run migrations, and get the application ready to run.

✅ Prerequisites

Before you begin, make sure you have the following tools installed on your system:

Go: The Go programming language. You can check if it's installed by running go version.

SQLite3: The command-line interface for SQLite. You can check for it by running sqlite3 --version. If it's not installed, you can get it from the official SQLite website or via a package manager (like brew install sqlite on macOS or sudo apt-get install sqlite3 on Debian/Ubuntu).

Bun CLI: The command-line tool for managing Bun migrations. Install it with this Go command:

go install github.com/uptrace/bun/buncli@latest

⚙️ Step 1: Create the SQLite Database File

An SQLite database is just a single file on your computer.

Navigate to the root directory of your Go project in your terminal.

Create the database file using the sqlite3 command. It's common to name it something like database.db.

# This command creates the file and opens the SQLite prompt.
sqlite3 database.db

You'll see a prompt like sqlite>. You don't need to do anything here. The file has been created. You can exit by typing .quit and pressing Enter.

sqlite> .quit

You should now have an empty database.db file in your project directory.

🔌 Step 2: Configure Your Go Application

Your Go application needs to know how to connect to this new database. This is done using a Data Source Name (DSN).

Find the Configuration: Locate where your application configures the database connection. This is often in a main.go file or a configuration file (config.yml, .env, etc.).

Set the DSN: The DSN for SQLite is very simple. You need to tell Bun that the SQL dialect is sqlite and provide the path to the database file.

// Example DSN for SQLite in your Go code
dsn := "database.db?_pragma=foreign_keys(1)" // The pragma enables foreign key support

// Connect to the database
sqldb := sql.OpenDB(sqliteshim.NewConnector(dsn, nil))
db := bun.NewDB(sqldb, sqlitedialect.New())

Make sure your application code uses a DSN like this. The _pragma=foreign_keys(1) part is important for ensuring relationships between tables (like your users and short_links) are enforced.

🏃 Step 3: Run the Bun Migrations

Now it's time to create the tables and seed the initial data using the migration file you wrote.

Set the DSN Environment Variable: The buncli tool reads the DSN from an environment variable named BUN_DSN. Set this in your terminal.

# For macOS/Linux
export BUN_DSN="database.db?_pragma=foreign_keys(1)"

# For Windows (Command Prompt)
set BUN_DSN="database.db?_pragma=foreign_keys(1)"

Check Migration Status: See which migrations Bun is aware of and which ones need to be applied.

bun migrate status

Apply Migrations: Run the up command to execute all pending migrations. This will run the Up function in your migration file, creating the users and short_links tables and inserting the admin user.

bun migrate up

You should see output indicating that the migration was successfully applied.

🧐 Step 4: Verify the Setup

Let's make sure everything worked as expected.

Open the Database: Use the SQLite CLI to inspect the database file.

sqlite3 database.db

List Tables: At the sqlite> prompt, use the .tables command to see if your tables were created.

sqlite> .tables
short_links  users

You should see short_links and users.

Check for Admin User: Run a SQL query to see if the admin user was inserted correctly.

sqlite> SELECT email, is_active FROM users WHERE email = 'admin@example.com';

The expected output is:

admin@example.com|1

Exit SQLite: Type .quit to exit.

You're all set! Your database is now initialized and contains the necessary tables and seed data. You can now run your Go application.