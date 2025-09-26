# GHOST 👻

<p align="center">
    <img src="docs/images/ghost.png" alt="drawing" width="600"/>
</p>

## What we do?
You know developing multiple projects/services locally can be a pain in the a**...

You could potentially switch to running a local Kubernetes cluster, sure.
But why bother configuring all of this since you already got your Docker containers running simultaneosly.

And with that, you need to choose your ports for each service...
```
localhost:7000
localhost:7001
...
localhost:7010
```

You can imagine where this is headed... Absolute *LOCAL* chaos!

And ghost is your salvation!!!

What we achieve is simple, but it helps organize your workbench, I mean, your local environment.

We wrap all your servers into a proxy server, to help you access your services using a meaningful nameserver, other than having to memorize which port is which.

So you would access (and it would be proxied for):

```
service-a.local -> localhost:7000
service-b.local -> localhost:7001
...
service-j.local -> localhost:7010
```

Plain and simple.

## How to install

If you already got your Golang installation, just do a:
```bash
> go install github.com/lccmrx/ghost@latest
```

### WIP
Or use HomeBrew:
```bash
brew install ghost
```

## How to use

We built a simple CLI tool to help you through the first config process.

Simply run a:
```bash
> ghost setup
```
And follow the config steps.

After finishing setting up, a:
```bash
> ghost start
```
is necessary to spin up all the necessary containers.

Afterwards,
```bash
> ghost add service-a 7000
```
will add the service to your local proxy registrar.

If you want to remove a name server, run this:
```bash
> ghost remove service-a
```

## Limitations

Today we run the setup under what we call a
> L-TLD - LOCAL TOP-LEVEL DOMAIN

So when configured to a LTLD `*.ghost` you won't be able to switch unless you reset your environment.

> We strongly advise to avoid using the `.local` LTLD, 
> since it's used in Avahi/Bonjour services for the mDNS services
