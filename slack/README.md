# Slack client

## 1. Install module

```shell
go get -u github.com/tuhin37/goclient
```

## 2. Import package

```go
import "github.com/tuhin37/goclient/slack"
```

## 3. Create slack client

```go
slackClient := slack.NewSlackClient("https://hooks.slack.com/services/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", "#hypd-invites", "Jhon Doe")
```

## 4. Send plain text

```go
err := slackClient.PushText("this is text body", ":smile:")
if err != nil {
    log.Fatal(err)
}
```

This is how it will look in slack
![loading-ag-305](./doc/2023-05-27-11-24-56-image.png)

## 5. Send json

```go
data := map[string]interface{}{"name": "tuhin", "age": 30}

err := slackClient.Push("header-text", data, ":ghost:")
if err != nil {
    log.Fatal(err)
}
```

This is how the result will look

![](./doc/2023-05-27-11-28-17-image.png)
