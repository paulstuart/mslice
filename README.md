# mslice
Generic mmap'd slices

This package facilitates the creation of a mmap'd back slice of type T.
It uses the unsafe package so those caveats apply.

Note: this will break if appending to the slice exceeds the initial capacity.

Disclaimer: not heavily tested so discretion is advised.
