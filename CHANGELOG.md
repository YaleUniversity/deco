# Changelog

## 1.0.0

* default region to us-east-1 for ssm
* compress release artifacts

## 0.5.0

* Add support for base64 encoded deco control files

## 0.4.1

* Support pullinh decofile from AWS SSM Parameters
* migrate to go modules

## 0.3.1

* Add more details to help text

## 0.3.0

* Add `show` subcommand
* [BREAKING] stop using `--file controlFile` and instead use the idiomatic arg format: `deco run /path/to/file`
* Add support for http/https URLs for control file arguments

## 0.2.2

* Migrate to github.com

## 0.2.1

* Remove unused `config` option

## 0.2.0

* Fix base64 decode
* Add support to pass a base directory

## 0.1.0

* Initial Release
