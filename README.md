gopush
======
Prevent pushing to restricted branches

Table of Content
----------------
* [Usage](#usage)
* [Requirements](#requirements)
* [Install](#install)
* [Autocomplete](#autocomplete)
* [License](#license)
* [Author](#author)

Usage:
------
List your restricted branches in `.gopush_restricted` - one branch per line. The file can be in the current project or parent directory. 
`gopush` will look from the current directory and continue up to `$HOME`.

```
gopush origin master
# Or
gopush master
```


Requirements:
-------------
* Go

Install:
--------
```
go get github.com/baopham/gopush
```

Autocomplete
-----------

To have autocomplete enabled, source [gopush_bash_autocomplete](autocomplete/gopush_bash_autocomplete) or [gopush_zsh_autocomplete](autocomplete/gopush_zsh_autocomplete).
E.g. copy one of these (depending on your shell) to `/local/path` and then add the below to your `.bashrc` or `.zhsrc`:

```bash
source /local/path/gopush_bash_autocomplete
# Or
source /local/path/gopush_zsh_autocomplete
```

License:
--------
MIT

Author:
-------
Bao Pham
