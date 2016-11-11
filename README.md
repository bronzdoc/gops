### gops - a toy port scanner built in GO

#### Install
```go get github.com/bronzdoc/gops```

#### Usage

To scan a  remote server<br />
```gops -host=somehost.com -tcp```

To scan locally<br />
```gops -tcp```

#### Flags
    -host  host to scan
    -tcp   Show only tcp ports open
    -udp   Show only udp ports open
    -start Port to start the scan
    -end   Port to end the scan
