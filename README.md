# Trans-socket

With this library you can transfer live network connections between two processes. I.e. you can update your service 
without lose the established TCP connections.

# How does it work

*man 7 unix*
```
SCM_RIGHTS
    Send or receive a set of  open  file  descriptors  from  another
    process.  The data portion contains an integer array of the file
    descriptors.
    Commonly, this operation is referred to as "passing a  file  de‐
    scriptor" to another process.  However, more accurately, what is
    being passed is a reference to an  open  file  description  (see
    open(2)),  and in the receiving process it is likely that a dif‐
    ferent file descriptor number will be used.  Semantically,  this
    operation  is equivalent to duplicating (dup(2)) a file descrip‐
    tor into the file descriptor table of another process.
```

# Inspiration
- [link 1.](https://github.com/mindreframer/golang-stuff/blob/master/github.com/youtube/vitess/go/umgmt/fdpass.go)
- [link 2.](https://github.com/ftrvxmtrx/fd/)
