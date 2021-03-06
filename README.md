# Joyread
Joyread is a lightweight self-hosted ebook reader written in Go. Licensed under AGPLv3. UNDER ACTIVE DEVELOPMENT.

### Cross-platform
Joyread runs anywhere Go can compile for: Windows, macOS, Linux, ARM, etc. Choose the one you love!

### Easy installation
Simply run the binary for your platform. Or ship Joyread with Docker.

### Share your ebooks
Share ebooks with your family and friends: Being a multi-user product, this software can be utilized for sharing ebooks public to all users in the platform or share only with selected users. You can also keep your ebooks private.
 
### Full-text search
Full-text search across all your ebooks: Search by metadata, content and tags.

### Nextcloud integration
You might already have your ebook collection on your Nextcloud: Why not just use Nextcloud sync feature in Joyread to grab all the ebooks and read it?

### Source folder sync
Do you find it cumbersome to upload your massive collection of ebooks via HTTP file upload? In Joyread, each user has a separate source folder where the respective user ebooks will be stored. You can SSH or FTP (or whatever way) upload all your ebooks to your source folder in the Joyread server and sync it via the web interface.

# Setup
Joyread is under development. It is not ready for production use.
### Prerequisites
 - PostgreSQL 10
 - Create a new role and database in PostgreSQL. Example shown below
   ```
   CREATE ROLE joyreaduser WITH LOGIN PASSWORD 'jellyfish' VALID UNTIL 'infinity';
 
   CREATE DATABASE joyreaddb WITH ENCODING='UTF8' OWNER=joyreaduser CONNECTION LIMIT=-1;
   ```
 ### Development
  - Clone the repo and put it in an appropriate `GOPATH`. For eg: `$GOPATH/src/github.com/joyread/server`
  - Configure the values in `config/app.yaml`
  - Run `go get -d ./...` or `dep ensure` inside the project folder
  - Then `go run ./cmd/joyread/main.go`. This will run the joyread server on the port mentioned in the `app.yaml` configuration. Default port is `8080`
  - In order to compile SCSS, you can do `gulp` inside the project folder
