# example-host

A service host that runs an AI agent managing resiliency. The functionality for this service is contained in the imported modules, so there is only a main.go file for this service.
Main.go has 1 function that handles the following startup responsabilities:

  1. Initialize the runtime environment and access logging.
~~~
// Set runtime environment - defaults to debug
runtime2.SetTestEnvironment()

// Initialize access logging handler and options
//access.SetLogHandler(nil)
access.EnableTestLogHandler()
access.EnableInternalLogging()
~~~


