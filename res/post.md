# Post

Observability in any application is key to monitoring and digesting user flows, errors and deployments. During the
lifecycle of some of our applications at Reddico we have implemented a robust logging platform that stores log entries
to Mongo and sends a message to a chat channel if the entry has an internal error.

This has made debugging our applications extremely rapid; not only that, it allows us to see where and if something went
wrong in our API's during user flows. Inside our applications, we are able to monitor Lambda tasks within the UI which
makes debugging errors a breeze.

Below are three packages we have developed to aid with this approach, check them out!

https://github.com/ainsleyclark/logger
https://github.com/ainsleyclark/errors
https://github.com/ainsleyclark/mogrus
