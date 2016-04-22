Ignore
=====================

## Description ##
Ignore is a feature that allows users to define patterns and files that they wish to not be synced over the network.
The syntax of a .pfsignore file is similar to .gitignore

## Usage ##

PfsIgnore works by reading patterns from a .pfsignore file in the root of the FUSE file-system. .pfsignore files cannot be ignored and are synced automatically.

## Patterns ##

Currently .pfsignore supports patterns separated by `\n`. Currently .pfsignore supports the following patterns.

1. Exact Patterns:
    Any Pattern that matches the path of a file from the root of the fs to the file

2. Partial Patterns:
    Any Pattern that ignores a directory above files will ignore files in that directory

3. Glob Patterns:
    `*` acts as a wildcard and will match any pattern on that level of the directory

4. Recursive Globbing:
    `**` acts like `*` but will only care about the start and end of the glob

5. Negation
    `!` will allow files that have been ignored by another pattern to be included in syncing.
    (if a parent dir is ignored then negation removed)

## Resyncing ##

Currently files removed from .pfsignore will not sync automatically. These files will either have to be Removed or truncated to size 0
for syncing to occur. If a file is ignored but no changes have happened to it since it was ignored the file will continue to sync as normal.
