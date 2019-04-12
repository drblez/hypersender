# hypersender
Send files from directories to «...the, uh, Internets…»

```
Usage:
  hypersender [OPTIONS]

Application Options:
      --debug                  Debug level logging [$DEBUG]
      --console                Output to console [$CONSOLE]
  -p, --path=                  Path to scan (default: .)
  -u, --url=                   URL to send
  -s, --path-substitution      Substitute file name in place of %f in URL
  -q, --substitute-sequence=   Change default sequence '%f' to user sequence (default: %f)
      --log-path=              Path to save log (default: .)
  -f, --fs-parallelism=        Number of workers for file system operations (default: 10)
  -n, --net-parallelism=       Number of worker for network operations (default: 10)
  -t, --content-type=          Content type (default: application/json)
  -E, --panic-on-errors        Panic on error
  -I, --ignore-service-errors  Ignore non-200 status code
  -S, --strip-path             Strip path from substitution
  -P, --file-name-pattern=     Send only file with name matched with pattern
      --dry-run                Do dry run
  -o, --timeout=               Network timeout (default: 30s)

Help Options:
  -h, --help                   Show this help message
```

## Examples:

### Send files *.json from current directory to endpoint

    hypersender -u 'http://localhost:88/upload' -P '*.json'
    
### Send files *.json from current directory to endpoint with file substitution

    hypersender -u 'http://localhost:88/upload?file=%f' -P '*.json' -s
