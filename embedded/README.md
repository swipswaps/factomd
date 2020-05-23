# factomd.so

Build shared library to use a factom follower embedded as a shared library.

Here's an example for python.
```
class Factomd(threading.Thread):

    def __init__(self, *args, **kwargs):
        self.factomd = cdll.LoadLibrary("/path/to/factomd.so")
        self.factomd.Serve.argtypes = []
        self.factomd.Shutdown.argtypes = []
        super(Factomd, self).__init__(*args, **kwargs)

    def run(self):
        self.factomd.Serve()

    def join(self, *args, **kwargs):
        self.factomd.Shutdown()
        super(Factomd, self).join(*args, **kwargs)
```

## Motivation

Deploy a Dapp factom follower using a library.

Planned: use to support python Dapp framework for Raspberry pi.

### Convenience

 As a developer if you wish to write a DAPP on factom
you have to reference an external instance of factomd.

### Deployment

With this setup it's possible to coodinate the restart
of the embedded folower in the context of the parent app.


### Versioning

You can version the factomd follower using your app's package management tools.
