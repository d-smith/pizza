pizza - View Pizza order status

This project provides an example of building a status visualization using
the general status tracking implementation from the [status](https://github.com/d-smith/status)
project.

The generation of events and writing them to the status event stream
is currently accomplished in the jupyter notebook in the status 
project. This project serves up two URIs - one to list orders (/orders), and 
one to drill into step progression for specific orders (/orderstatus).

To run, copy setenv-template and populate it with environment variable
values specific to your setup, source it, then `go run main.go`

The process listens on port 8778.


Dependencies:

<pre>
go get github.com/gorilla/mux
go get github.com/aws/aws-sdk-go/...
</pre>
