# Go toolbox

Various packages for Go projects. This is a mixed bag of functions and packages
that are useful with no general theme.

We could name it **utils** but "toolbox" is a better description, this is a library
with all the weird parts you find in a proper toolbox with customized screwdrivers,
hacked-off box keys, welded-together contraptions and objects that make you
think "oh, I can't imagine what this is for".

Some of these are gathered from the now-defunct ExploratoryEngineering organization.

Note that there's no `go.mod` here - your project should include one.

## Referenced libraries

* The Kong library for parameters. This is used throughout the parameter
  structs.
* gRPC
* Logrus for logging
* grpc-middleware for metrics interceptors
* Prometheus (transient, through grpc-middleware)