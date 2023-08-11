# this is a justfile (https://github.com/casey/just)

# tests require xmlstarlet (https://xmlstar.sourceforge.net/)
# macOS: brew install xmlstarlet
# *nix: sudo apt-get install xmlstarlet

# default recipe to display help information
default:
  @just --list

# run example 1
s01:
  go run 01-getkev.go

# run example 2
s02:
  go run 02-kev2feed.go

# run example 3
s03:
  go run 03-kev2feedfile.go

# build final project
@build:
  go build -o kev2feed 03-kev2feedfile.go

# test RSS output (see justfile for required dependency)
@test:
  (command -v xmlstarlet > /dev/null) || (echo "Please install xmlstarlet (https://xmlstar.sourceforge.net/)" ; exit 1)
  just build && xmlstarlet val --well-formed kev-rss.xml