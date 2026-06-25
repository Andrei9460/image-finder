# image-finder
поиск копий изображения, 

CentOS 7 на старом железе может не поддерживать современные инструкции CPU. Соберите с явным указанием минимальной архитектуры:
```
set GOOS=linux
set GOARCH=amd64
set GOAMD64=v1
set CGO_ENABLED=0
go build -ldflags="-s -w" -o image-finder main.go
```

# Запустить
./image-finder img_name.png /home/bitrix/www/upload/iblock 5

image-finder.exe "D:\путь\к\вашей\картинке.jpg" "D:\путь\к\папке" 5
