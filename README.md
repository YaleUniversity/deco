# deco

`deco` is short for **D**ocker **E**nvironment **Co**ntrol

## Usage

`deco` has three options currently, `validate`, `show`, and `run`.

```text
deco gets your app ready to run when a container
starts.  For example: the filters allow you to specify
individual files to filter and key/value pairs to use when
filtering.  By default, it works from the current directory and
will filter files in place.

Usage:
  deco [command]

Available Commands:
  help        Help about any command
  run         Run executes the tasks in the given control file
  show        Reads and displays a control file on STDOUT
  validate    Validates the control file
  version     Displays version information

Flags:
      --config string   deco config file -- _not_ the control file (default is $HOME/.deco.yaml)
  -d, --dir string      Base directory for filtered files/templates
  -h, --help            help for deco

Use "deco [command] --help" for more information about a command.
```

### Input

`deco` takes a JSON file as input and defaults to `/var/run/secrets/deco.json`.  This allows it to work
out of the box in docker swarm with swarm secrets.  The control file can be base64 encoded (standard encoding)
using the `--encoded` flag

The JSON control file has the format:

```JSON
{
    "filters": {
        "path/to/file1": {
            "replace": "withThis",
            "doThis": "true"
        },
        "path/to/file2": {
            "base64encodedstuff": "IMKvXF8o44OEKV8vwq8="
        }
    }
}
```

It's possible it could do more than just filter in the future.

`deco` can currently source its JSON control file from:

- `/var/run/secrets/deco.json` (default)
- an arbitrary local filesystem with an absolute or relative path
- an http(s) endpoint
- an ssm parameter store endpoint

`deco` also supports passing custom headers for doing things like basic auth to the http(s) endpoint

`deco show http://127.0.0.1:8888/v1/deco.json -H 'Authorization=Basic YWRtaW46cGFzc3dvcmQ='`

## Examples

**Note:** the examples directory has working cases with base64 encoding/decoding.

### Case #1 - basic usage with local filesystem

```
deco run /app/params.json -d /app
[INFO] Using control from file /app/params.json
Filtering /app/configdir/configfile.json
```

**Note:** /app/params.json is the control file.  It contains filters and configuration data inside the filter key:
  - A filter key itself is the relative or full path to the template config file
  - The filter value contains the data to replace within the template
  - -d is an optional basedir.  If you specify full path to the template in the control file, you don't have to specify -d.  If you change dir - _cd_ - to the parent dir above the project, you don't have to specify -d.

#### Control file:

```JSON
{
    "filters": {
        "configdir/configfile.json": {
            "bar": "bar-app",
            "host01": "host01.example.org",
            "path01": "/tmp",
            "baz": "baz-app",
            "host02": "host02.example.org",
            "path02": "/u0/app",
            "log_level": "/debug"
        }
    }
}
```

#### File template:

In our example it exists as /app/configdir/configfile.json.  It will be written over in-place.

/bin/cp /app/configdir/configfile.json.orig /app/configdir/configfile.json

cat /app/configdir/configfile.json

```JSON
{
  "foo": {
    "{{ .bar }}": {
      "host": "{{ .host01 }}",
      "path": "{{ .path01 }}"
    },
    "{{ .baz }}": {
      "host": "{{ .host02 }}",
      "path": "{{ .path02 }}"
    }
  },
  "log_level": "{{ .log_level }}"
}
```

result:

$ cat /app/configdir/configfile.json                                        

```JSON
{ 
  "foo": {
    "bar-app": {
      "host": "host01.example.org",
      "path": "/tmp"
    },
    "baz-app": {
      "host": "host02.example.org",
      "path": "/u0/app"
    }
  },
  "log_level": "debug"
}
```

### Case #2 - use AWS SSM parameter store to get control file

if [ -n "$SSMPATH" ]; then
  echo "Getting config file from SSM Parameter Store (${SSMPATH}) ..."
  deco version
  if [[ $? -ne 0 ]]; then
    echo "ERROR: deco not found!"
    exit 1
  fi
  deco validate -e ssm://${SSMPATH} || exit 1
  deco run -e ssm://${SSMPATH}
else
  echo "ERROR: SSMPATH variable not set!"
  exit 1
fi

## Author

E Camden Fisher <camden.fisher@yale.edu>

## License

The MIT License (MIT)

Copyright (c) 2021 Yale University

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
