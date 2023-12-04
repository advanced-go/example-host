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

  2. Application startup. The messaging module contains an exchange that packages can register a mailbox with. This facilitates messaging and is utilized by the service host to startup
     a package. The configuration needed for a package should be created by the service host. Messaging also supports a ping request for a selected package. The application agent is also
     started. 
~~~
// Run startup where all registered resources/packages will be sent a startup message which may contain
// package configuration information such as authentication, default values...
m := createPackageConfiguration()
status := exchange.Startup[runtime2.Log](time.Second*4, m)
if !status.OK() {
    return r, status
}

// Start application agent
agent.Run(time.Second * 10)
~~~

  3. Request routing. The messaging module contains a multiplexer that routes requests to an HTTP handler. Routing requests to handlers that are not contained in the example-domiain are handled by the ServeMux. 
 ~~~
// Initialize messaging mux for all HTTP handlers in example-domain
mux.Handle(activity.PkgPath, activity.HttpHandler)
mux.Handle(slo.PkgPath, slo.HttpHandler)
mux.Handle(timeseries.PkgPath, timeseries.HttpHandler)
mux.Handle(google.PkgPath, google.HttpHandler)

// Initialize health liveliness handler
r.Handle(healthLivelinessPattern, http.HandlerFunc(healthLivelinessHandler))

// Route all other requests to messaging mux
r.Handle("/", http.HandlerFunc(mux.HttpHandler))
~~~




