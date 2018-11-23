[![Build Status](https://img.shields.io/travis/joonas-fi/dupfinder.svg?style=for-the-badge)](https://travis-ci.org/joonas-fi/dupfinder)
[![Download](https://img.shields.io/bintray/v/joonas/dupfinder/main.svg?style=for-the-badge&label=Download)](https://bintray.com/joonas/dupfinder/main/_latestVersion#files)

dupfinder finds duplicate files.

Features
--------

- Compares old and new directory for added/removed/unchanged files
- Emphasis on safety. Your files are precious!
- Understands renames (file content equality is based on sha1)
- Comparison results in a human readable log file, which you can re-read and delete the
  unchanged files mentioned in the log, either from the old or the new directory
- Additionally, even though the log-based approach brings safety, we have additional dry run
  in the delete process that explains which files it would delete.
- Can remove empty directories (that result from "remove same files from old OR new directory")


Usage
-----

```
$ dupfinder compare demo/old demo/new > comparison.log
2018/11/23 13:36:58 olddir<demo/old> newdir<demo/new>
2018/11/23 13:36:58 Scanning olddir
2018/11/23 13:36:58 starting initializeMissingMap
2018/11/23 13:36:58 Scanning newdir
2018/11/23 13:36:58 Listing missing files
2018/11/23 13:36:58 Done

$ cat comparison.log
INFO olddir<demo/old> newdir<demo/new>
INFO Scanning olddir
INFO starting initializeMissingMap
INFO Scanning newdir
+ demo/new/subfolder/stoned sloth.jpg
= HgAvAA demo/old/subfolder/tcpdump.PNG|demo/new/subfolder/tcpdump renamed.PNG
= NgA2AA demo/old/pcap-screenshot-cloudflare-duplicate-acks.PNG|demo/new/pcap-screenshot-cloudflare-duplicate-acks.PNG
INFO Listing missing files
- demo/old/subfolder/callingcard.jpg
INFO Done

$ dupfinder removeunchangedfilesfromnew < comparison.log
would remove demo/new/subfolder/tcpdump renamed.PNG
would remove demo/new/pcap-screenshot-cloudflare-duplicate-acks.PNG

$ dupfinder removeunchangedfilesfromnew --really < comparison.log
removing demo/new/subfolder/tcpdump renamed.PNG
removing demo/new/pcap-screenshot-cloudflare-duplicate-acks.PNG

$ dupfinder removeemptydirs demo/new
would remove empty dir: demo/new/emptydir

$ dupfinder removeemptydirs --really demo/new
removing empty dir: demo/new/emptydir

```

Alternatives
------------

- https://github.com/jpillora/dedup
