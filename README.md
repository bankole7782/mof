# mof

mof: make one file.

An archive algorithm.

API documentation can be found on [godoc](https://godoc.org/github.com/bankole7782/mof)


## How it Works
First use the `MTF` function to produce a two files and then pass it to `MOF` function 
which would now produce one file.


### Description (MTF)

#### Packing

The first file is the index file which contains the fullname ( parent directory(s) and file name)
of every file in the archive and the size of the file.

The second file is the data file which contains the contents of every file written to it
in the order it was written in the index file (first file).

#### Unpacking

For the first entry (full file name and path) in the index file, you create the file with the name
from the index file and read from the data file the size of the entry and store in the newly created
file.

For other entries you sum up the sizes of all the entries before the entry you want to extract. Then
you seek to this sum, create the full file name from the current entry and then read data with the size
of the data file and write it to the full file name.


### Description (MOF)

#### Packing
First create an MTF archive.

Make a new file and write the size of the index file from above and add a newline. Then write the index file
and then the data file.


#### Unpacking
Loop until you get to the first newline and use that to get the size of the index file.

Read the index file into a temporary file (use the size gotten above) and start from the
size of the size of the index file with the newline.

Then compute the size of the data file subtracting the size of the index file and the size of the size of
the index file. Use this size to write the data file to a new file.

Then pass the gotten index file and data file to the `UndoMTF` function.

# License

Released with MIT