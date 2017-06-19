gopush
======
Prevent pushing to restricted branches

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

License:
--------
MIT

Author:
-------
Bao Pham
