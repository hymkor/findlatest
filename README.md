findlatest
==========

`findlatest` finds the latest file in the specified files and directories.

```
findlatest {options} {FILES_OR_DIRECTORIES}
```

`findlatest` displays the timestamp and filename each time it finds a file with a newer timestamp.

### Options:

- `-q` : Be quiet. It suppresses intermediate output and prints only the latest timestamp and filename at the end.
- `-a` : Include dotfiles (files and directories starting with a dot) in the search
- `-until <timestamp>` : Exclude files newer than the specified timestamp (format: `2006-01-02 15:04:05`).

## LICENSE

MIT
