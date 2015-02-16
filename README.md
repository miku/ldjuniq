ldjuniq
=======

Like `uniq`, but for LDJ files.

Example
-------

Iterate over all lines in `file.ldj` and drop duplicate
lines with the same value in `document.URL`.

    $ ldjuniq -key 'document.URL' file.ldj

