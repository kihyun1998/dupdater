# 빌드하는 방법

```bash
fyne package -icon icon.png -name updater
```


## 커스텀 아이콘 사용하는 방법

```bash
fyne bundle .\icon.png > icon_resources.go
```

하고 icon_resources.go를 utf-8로 인코딩 돼있는지 확인해봐야함.

```go
icon := fyne.NewStaticResource(resourceIconPng.StaticName, resourceIconPng.StaticContent)
ui.spinnerIcon = widget.NewIcon(icon)
...
```

위처럼 사용할 수 있다.