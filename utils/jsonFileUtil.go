package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

func ReadJSONFromFile(filepath string, data interface{}) ([]byte, error) {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("file read failed: %w", err)
	}
	if err := json.Unmarshal(fileContent, data); err != nil {
		return fileContent, fmt.Errorf("JSON parse failed: %w", err)
	}
	return fileContent, nil
}

func CopyDirSrc(src string, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("无法获取源文件夹信息: %v", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("源路径不是一个文件夹: %v", src)
	}

	var tempDest string
	var needRollback bool

	// 处理目标文件夹已存在的情况
	if _, err := os.Stat(dest); err == nil {
		tempDest = dest + ".tmp"

		// 移除可能存在的旧临时文件夹
		if err := os.RemoveAll(tempDest); err != nil {
			return fmt.Errorf("清理旧临时文件夹失败: %v", err)
		}

		// 将现有目标文件夹重命名为临时文件夹
		if err := os.Rename(dest, tempDest); err != nil {
			return fmt.Errorf("重命名目标文件夹失败: %v", err)
		}
		needRollback = true
	}

	// 创建目标文件夹
	if err := os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		if needRollback {
			// 恢复原始文件夹
			if rerr := os.Rename(tempDest, dest); rerr != nil {
				return fmt.Errorf("创建目录失败且回滚失败: %v (原错误: %v)", rerr, err)
			}
		}
		return fmt.Errorf("创建目标文件夹失败: %v", err)
	}

	// 执行复制操作
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath)
	})

	// 处理复制结果
	if err != nil {
		// 删除新建的目标文件夹
		if rerr := os.RemoveAll(dest); rerr != nil {
			fmt.Printf("警告：清理目标文件夹失败: %v\n", rerr)
		}

		// 执行回滚
		if needRollback {
			if rerr := os.Rename(tempDest, dest); rerr != nil {
				return fmt.Errorf("复制失败且回滚失败: %v (原错误: %v)", rerr, err)
			}
		}
		return fmt.Errorf("复制过程失败: %v", err)
	}

	// 删除临时备份文件夹
	if needRollback {
		if err := os.RemoveAll(tempDest); err != nil {
			return fmt.Errorf("复制成功但清理临时文件夹失败: %v", err)
		}
	}

	return nil
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("无法打开源文件: %v", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("无法创建目标文件: %v", err)
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("复制文件内容失败: %v", err)
	}

	// 同步文件内容到磁盘
	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("文件同步失败: %v", err)
	}

	return nil
}

func SaveJSONToFile(data interface{}, filepath string) error {
	// 打开文件
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// 如果文件不存在，创建父目录
		parentDir := path.Dir(filepath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return err
		}
	}
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 设置缩进
	return encoder.Encode(data)
}

func JoinPaths(parts ...string) string {
	return path.Join(parts...)
}

func IsDirExist(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func RemoveFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("删除文件失败: %v\n", err)
		return
	}
	fmt.Printf("文件已删除: %s\n", filePath)
}
