![](https://habrastorage.org/web/9f9/84c/5c8/9f984c5c8bdf45b4b25c1802c4c8b843.png)

# Furdarius\Steamprotocol

TODO: More docs!

Example of usage:

```go
rand.Seed(time.Now().Unix())
cm := cmlist.NewCMList(cl)
servAddr, err := cm.GetRandomServer()
if err != nil {
    log.Error(
        "failed to get random server",
        zap.Error(err))

    os.Exit(1)
}

d := net.Dialer{}
conn, err := d.DialContext(ctx, "tcp", servAddr)
if err != nil {
    log.Error(
        "Failed to connect to server",
        zap.String("addr", servAddr),
        zap.Error(err))

    os.Exit(1)
}

log.Debug("Successfully connected to server",
    zap.String("addr", servAddr))

eventManager := steamprotocol.NewEventManager()

steamClient := steamprotocol.NewClient(conn, eventManager)

cryptoModule := crypto.NewModule(steamClient, eventManager)
cryptoModule.Subscribe()

gen := auth.NewTOTPGenerator(cl)
authModule := auth.NewModule(steamClient, eventManager, gen, auth.Details{
    Username:     "myusername",
    Password:     "mypassword",
    SharedSecret: "mysharedsecret",
})
authModule.Subscribe()

multiModule := multi.NewModule(steamClient, eventManager)
multiModule.Subscribe()

heartbeatModule := heartbeat.NewModule(steamClient, eventManager)
heartbeatModule.Subscribe()

socialModule := social.NewModule(steamClient, eventManager)
socialModule.Subscribe()

// TODO: Select with context for graceful shutdown
go func() {
    for err := range heartbeatModule.ErrorChannel() {
        log.Error(
            "failed to heartbeat",
            zap.Error(err))
    }

    log.Debug("Exit from heartbeat error checking goroutine")
}()

err = steamClient.Listen()
if err != nil {
    log.Error(
        "Failed to listen steam server",
        zap.Error(err))

    os.Exit(1)
}
```