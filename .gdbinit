add-auto-load-safe-path /usr/local/go/src/runtime/runtime-gdb.py
set disassembly-flavor intel
set arch i386:intel
file build/kernel-amd64.bin
target remote localhost:1234
set arch i386:x86-64:intel
layout split
break kernel.Kmain