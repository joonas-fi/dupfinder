dupfinder finds duplicate files.

Features
--------

- Compares old and new directory for added/removed/unchanged files
- Comparison results in a human readable log file, which you can re-read and delete the
  unchanged files mentioned in the log, either from the old or the new directory
- Additionally, even though the log-based approach brings safety, we have additional dry run
  in the delete process that explains which files it would delete.
- Can remove empty directories (that result from "remove same files from old OR new directory")


Alternatives
------------

- https://github.com/jpillora/dedup
