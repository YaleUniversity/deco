# deco

`deco` is short for **D**ocker **E**nvironment **Co**ntrol

## Usage

`deco` has two options currently, `validate` and `run`.

```
./deco --help
deco gets your app ready to run when a container
starts.  For example: the filters allow you to specify
individual files to filter and key/value pairs to use when
filtering.  By default, it works from the current directory and
will filter files in place.

Usage:
  deco [command]

Available Commands:
  help        Help about any command
  run         Run executes the taks in the given control file
  validate    Validates the control file
  version     Displays version information

Flags:
  -d, --dir string      Base directory for filtered files/templates
      --file string     location of control file (default "/var/run/secrets/deco.json")
  -h, --help            help for deco

Use "deco [command] --help" for more information about a command.
```

### Input

`deco` takes a JSON file as input and defaults to `/var/run/secrets/deco.json`.  This allows it to work
out of the box in docker swarm with swarm secrets.

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

## Author

E Camden Fisher <camden.fisher@yale.edu>

## License

The MIT License (MIT)

Copyright (c) 2017 Yale University

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