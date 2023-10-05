# xray-auc
Simple cmd utility to control xray users from a console

Supported operations: add user, remove user, get user traffic statistics (in development)

It was written in a hurry to quickly solve the problem of managing xray users, since xray does not provide a simple api for this.

Xray has grpc api, but it's difficult to use it from linux console by some automation scripts.

Some peace of code related to grpc api flow was borrowed from https://github.com/FranzKafkaYu/x-ui
