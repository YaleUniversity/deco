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

### Deco Control File

`deco` takes a JSON file, the control file, as input and defaults to `/var/run/secrets/deco.json`.  This allows it to work
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
- a local filesystem with an absolute or relative path
- an http(s) endpoint
- AWS SSM parameter store

`deco` also supports passing custom headers for doing things like basic auth to the http(s) endpoint

`deco show http://127.0.0.1:8888/v1/deco.json -H 'Authorization=Basic YWRtaW46cGFzc3dvcmQ='`

## Examples

**Note:** the examples directory has several working cases, including base64 encoding/decoding.

### Basic usage with local filesystem

The most basic usage is a control file in the default location and simply running `deco run`.

```bash
deco run
[INFO] Using control from file /var/run/secrets/deco.json
...
```

### Referencing a control file and setting a base directory

`deco` allows you to set a base directory for filtering.  This directory will be prepended to the `filter` paths in the control file.

With the following `params.json` file:

```JSON
{
    "filters": {
        "configdir/configfile.yaml": {
            "bar": "bar-app",
            "host01": "host01.example.org",
            "path01": "/tmp"
        }
    }
}
```

running

```bash
deco run params.json -d /app
```

If the template file `/app/configdir/configfile.yaml` exists, it will be written over in place, replacing `{{ .bar }}` with `bar-app`, `{{ .host01 }}` with `host01.example.org`, and `{{ .path01 }}` with `/tmp`.

### Use AWS SSM parameter store to get control file

```bash
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
```

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
