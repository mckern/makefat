# makefat
A tool for making fat OSX binaries (a portable lipo)

You give it some executables, it makes a fat executable from them. The fat executable will run on any architecture supported by one of the input executables.

```
makefat <output file> <input file 1> <input file 2> ...
```
# extractfat
A tool for extracting fat OSX binaries (a portable lipo)

You give it some executable, it extracts the executables of any architecture in it.

```
extractfat <fat_file>
```
