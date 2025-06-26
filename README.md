# twitter clone

![Build Status](https://github.com/marcusprice/twitter-clone/actions/workflows/test.yaml/badge)

## quickstart
### installation and setup
First clone repo and install dependencies:
```
git clone https://github.com/marcusprice/twitter-clone.git
cd twitter-clone && go mody tidy
```
sqlite is required to initialize and seed the database. To install on macos:
```
brew install sqlite
```
Linux will depend on you distro's package manager, for arch linux:
```
sudo pacman -S sqlite
```

Make a copy of .env-sample and replace environment variable values:
```
cp .env-sample .env
```

Then initiatlize and seed the database:
```
make init-db
make seed-db
```

At this point the core app should be good to go.

### running the app
To run the core app:
```
make run-core
```

To run reply-guy service (Ollama is required, more on that in the following 
section)

```
make run-reply-guy
```

To run all services (requires ollama to be installed and configured):
```
make run-all
```

Logs share the same standard output with service prefix.

### api documentation

In development mode, swagger api documentation is available at 
http://address:port/docs i.e.
```
http://localhost:42069/docs
```

### debugging

Install delve debugger:
```
go install github.com/go-delve/delve/cmd/dlv@latest
```
More info: 

[Delve](https://github.com/go-delve/delve/tree/master)


To debug the services:

```
make debug-core

make debug-reply-guy
```

## reply-guy
reply-guy is an API and job queue for LLM generated responses to posts and
comments.

When a user creates a post or comment and tags an AI account (i.e.
@dalecooper), the core service will send a request to the reply-guy service
with data about the comment. This data includes the post/comment content and
the author as well as context for the original post and comment thread.
reply-guy will add the request to the queue, and processes the requests with a
single worker.

The request is parsed and formatted, and then sent to the ollama REST API to 
generate a response. After ollama responds with the AI generated content, 
reply-guy then makes a POST request to the core serivce's comment endpoint to
create the new comment.

The flow looks like this:
```
CLIENT -> CORE-SERVICE: POST /api/v1/comment/create
author @endlesshappiness
content: @dalecooper, is what OP saying true? has that actually been proven?
200 response on successful comment post, reply-guy call happens concurrently

CORE-SERVICE --> REPLY-GUY: POST /api/v1/@dalecooper REPLY_GUY --> OLLAMA: POST /api/generate
REPLY-GUY <--LLM generated content-- OLLAMA

REPLY-GUY --> CORE-SERVICE: POST /api/v1/comment/create
author @dalecooper
content: LLM generated content
```

### ollama
Ollama is required for the reply-guy, to install on mac:
```
$ brew install ollama
```

Linux:
```
curl -fsSL https://ollama.com/install.sh | sh
```

Then create the agent cooper LLM model:
```
ollama create dalecooper -f ./models/dalecooper.Modelfile
```

To serve ollama:
```
ollama serve
```

## tests:
Run all tests:
```
go test -v ./...
```
Run test for a specific package:
```
go test ./internal/model
```

Run a specific test:

```
go test ./internal/api --run ^TestCreateUser$
```

To run a test in debug mode:
```
make debug-test controller::TestPostNew
```



