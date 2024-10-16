# 빌드하는 방법

```bash
fyne package -icon icon.png -name updater
```


## 커스텀 아이콘 사용하는 방법

```bash
fyne bundle -o bundled.go icon.png
```

```go
    icon := fyne.NewStaticResource("icon", resourceIconPng.StaticContent)

	ui.spinnerIcon = widget.NewIcon(icon)
	ui.spinnerIcon.Resize(fyne.NewSize(50, 50))
```

위처럼 사용할 수 있다.


## 의존성 v2.5.2 기준

### fyne.io/systray v1.11.0 
링크: [링크](https://github.com/fyne-io/systray)
추가 여부 : `true`
라이선스: `Apache-2.0`

### github.com/BurntSushi/toml v1.4.0 
링크: [링크](https://github.com/BurntSushi/toml)
추가 여부 : `true`
라이선스: `MIT`

### github.com/davecgh/go-spew v1.1.1 
링크: [링크](https://github.com/davecgh/go-spew)
추가 여부 : `true`
라이선스: `ISC`

### github.com/fredbi/uri v1.1.0 
링크: [링크](https://github.com/fredbi/uri)
추가 여부 : `true`
라이선스: `MIT`

### github.com/fsnotify/fsnotify v1.7.0 
링크: [링크](https://github.com/fsnotify/fsnotify)
추가 여부 : `true`
라이선스: `BSD-3`

### github.com/fyne-io/gl-js v0.0.0-20220119005834-d2da28d9ccfe 
링크: [링크](https://github.com/fyne-io/gl-js)
추가 여부 : `true`
라이선스: `BSD-3`

### github.com/fyne-io/glfw-js v0.0.0-20240101223322-6e1efdc71b7a 
링크: [링크](https://github.com/fyne-io/glfw-js)
추가 여부 : `true`
라이선스: `MIT`

### github.com/fyne-io/image v0.0.0-20220602074514-4956b0afb3d2 
링크: [링크](https://github.com/fyne-io/image)
추가 여부 : `true`
라이선스: `BSD-3`

### github.com/go-gl/gl v0.0.0-20211210172815-726fda9656d6 z
링크: [링크](https://github.com/go-gl/gl)
추가 여부 : `true`
라이선스: `MIT`

### github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240506104042-037f3cc74f2a 
링크: [링크](https://github.com/go-gl/glfw)
추가 여부 : `true`
라이선스: `BSD-3`

### github.com/go-text/render v0.2.0 
링크: [링크](https://github.com/go-text/render)
추가 여부 : `true`
라이선스: `BSD-3`

### github.com/go-text/typesetting v0.2.0 
링크: [링크](https://github.com/go-text/typesetting)
추가 여부 : `true`
라이선스: `BSD-3`

### github.com/godbus/dbus/v5 v5.1.0 
링크: [링크](https://github.com/godbus/dbus)
추가 여부 : `true`
라이선스: `BSD-2`

### github.com/gopherjs/gopherjs v1.17.2 
링크: [링크](https://github.com/gopherjs/gopherjs)
추가 여부 : `true`
라이선스: `BSD-2`

### github.com/jeandeaual/go-locale v0.0.0-20240223122105-ce5225dcaa49 
링크: [링크](https://github.com/jeandeaual/go-locale)
추가 여부 : `true`
라이선스: `MIT`

### github.com/jsummers/gobmp v0.0.0-20151104160322-e2ba15ffa76e 
링크: [링크](https://github.com/jsummers/gobmp)
추가 여부 : `true`
라이선스: `MIT`

### github.com/nicksnyder/go-i18n/v2 v2.4.0 
링크: [링크](https://github.com/nicksnyder/go-i18n)
추가 여부 : `true`
라이선스: `MIT`

### github.com/pmezard/go-difflib v1.0.0 
링크: [링크](https://github.com/pmezard/go-difflib)
추가 여부 : `true`
라이선스: `BSD-3`

### github.com/rymdport/portal v0.2.6 
링크: [링크](https://github.com/rymdport/portal)
추가 여부 : `true`
라이선스: `Apache-2.0`

### github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c 
링크: [링크](https://github.com/srwiley/oksvg)
추가 여부 : `true`
라이선스: `BSD-3`

### github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef 
링크: [링크](https://github.com/srwiley/rasterx)
추가 여부 : `true`
라이선스: `BSD-3`

### github.com/stretchr/testify v1.8.4 
링크: [링크](https://github.com/stretchr/testify)
추가 여부 : `true`
라이선스: `MIT`

### github.com/yuin/goldmark v1.7.1 
링크: [링크](https://github.com/yuin/goldmark)
추가 여부 : `true`
라이선스: `MIT`

### golang.org/x/image v0.18.0 
링크: [링크](https://github.com/golang/image)
추가 여부 : `true`
라이선스: `BSD-3`

### golang.org/x/mobile v0.0.0-20231127183840-76ac6878050a 
링크: [링크](https://github.com/golang/mobile)
추가 여부 : `true`
라이선스: `BSD-3`

### golang.org/x/net v0.25.0 
링크: [링크](https://github.com/golang/net)
추가 여부 : `true`
라이선스: `BSD-3`

### golang.org/x/sys v0.20.0 
링크: [링크](https://github.com/golang/sys)
추가 여부 : `true`
라이선스: `BSD-3`

### golang.org/x/text v0.16.0 
링크: [링크](https://github.com/golang/text)
추가 여부 : `true`
라이선스: `BSD-3`

### gopkg.in/yaml.v3 v3.0.1 
링크: [링크](https://github.com/go-yaml/yaml)
추가 여부 : `true`
라이선스: `MIT`

## go.mod

```go
module fync-updater

go 1.22.2

require (
	fyne.io/fyne/v2 v2.5.1
	golang.org/x/sys v0.20.0
)

require (
	fyne.io/systray v1.11.0 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fredbi/uri v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fyne-io/gl-js v0.0.0-20220119005834-d2da28d9ccfe // indirect
	github.com/fyne-io/glfw-js v0.0.0-20240101223322-6e1efdc71b7a // indirect
	github.com/fyne-io/image v0.0.0-20220602074514-4956b0afb3d2 // indirect
	github.com/go-gl/gl v0.0.0-20211210172815-726fda9656d6 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240506104042-037f3cc74f2a // indirect
	github.com/go-text/render v0.1.1-0.20240418202334-dd62631dae9b // indirect
	github.com/go-text/typesetting v0.1.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jeandeaual/go-locale v0.0.0-20240223122105-ce5225dcaa49 // indirect
	github.com/jsummers/gobmp v0.0.0-20151104160322-e2ba15ffa76e // indirect
	github.com/nicksnyder/go-i18n/v2 v2.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rymdport/portal v0.2.6 // indirect
	github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c // indirect
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/yuin/goldmark v1.7.1 // indirect
	golang.org/x/image v0.18.0 // indirect
	golang.org/x/mobile v0.0.0-20231127183840-76ac6878050a // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

```