package main

import (
    "fmt"
    "image"
    "os"
    "path/filepath"
    "strings"
    "time"
    
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"
    
    "github.com/corona10/goimagehash"
    "github.com/disintegration/imaging"
    _ "golang.org/x/image/webp"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Использование: image-finder.exe <целевая_картинка> <папка_поиска> [порог]")
        fmt.Println("Пример: image-finder.exe D:\\photo.jpg D:\\images 5")
        fmt.Println("Порог: 0-10 (5 = рекомендуется)")
        return
    }
    
    targetPath := os.Args[1]
    searchRoot := os.Args[2]
    threshold := 5
    
    if len(os.Args) > 3 {
        fmt.Sscanf(os.Args[3], "%d", &threshold)
    }
    
    // Проверка существования
    if _, err := os.Stat(targetPath); os.IsNotExist(err) {
        fmt.Printf("Ошибка: целевой файл не найден: %s\n", targetPath)
        return
    }
    
    if _, err := os.Stat(searchRoot); os.IsNotExist(err) {
        fmt.Printf("Ошибка: папка поиска не найдена: %s\n", searchRoot)
        return
    }
    
    fmt.Printf("Поиск копий: %s\n", targetPath)
    fmt.Printf("В папке: %s\n", searchRoot)
    fmt.Printf("Порог точности: %d\n\n", threshold)
    
    startTime := time.Now()
    
    // Хеш целевого изображения
    targetHash, err := getImageHash(targetPath)
    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
        return
    }
    
    fmt.Println("Обработка файлов...")
    
    found := make([]string, 0)
    total := 0
    
    filepath.Walk(searchRoot, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return nil
        }
        
        ext := strings.ToLower(filepath.Ext(path))
        if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
            return nil
        }
        
        // Пропускаем исходный файл
        if isSameFile(targetPath, path) {
            return nil
        }
        
        total++
        
        if total%50 == 0 {
            fmt.Printf("Обработано: %d, найдено: %d\r", total, len(found))
        }
        
        currHash, err := getImageHash(path)
        if err != nil {
            return nil
        }
        
        distance, _ := targetHash.Distance(currHash)
        if distance <= threshold {
            found = append(found, path)
            fmt.Printf("\n✅ Найдено: %s (различие: %d)\n", path, distance)
        }
        
        return nil
    })
    
    elapsed := time.Since(startTime)
    
    fmt.Printf("\n\n📊 Обработано файлов: %d\n", total)
    fmt.Printf("⏱️  Время: %v\n", elapsed)
    
    if len(found) == 0 {
        fmt.Println("\n❌ Копий не найдено")
    } else {
        fmt.Printf("\n✅ Найдено %d копий:\n", len(found))
        for i, path := range found {
            fmt.Printf("%d. %s\n", i+1, path)
        }
    }
}

func getImageHash(filePath string) (*goimagehash.ImageHash, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    img, _, err := image.Decode(file)
    if err != nil {
        return nil, err
    }
    
    normalized := imaging.Resize(img, 256, 256, imaging.Lanczos)
    
    return goimagehash.PerceptionHash(normalized)
}

func isSameFile(path1, path2 string) bool {
    info1, err1 := os.Stat(path1)
    info2, err2 := os.Stat(path2)
    if err1 != nil || err2 != nil {
        return false
    }
    return os.SameFile(info1, info2)
}
