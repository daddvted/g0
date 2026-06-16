package main

import (
    "fmt"
    "io"
    "os"

    "github.com/pkg/sftp"
    "golang.org/x/crypto/ssh"
)

func main() {
    // 公钥认证配置
    config := &ssh.ClientConfig{
        User:            "psoss_uat",                              // 服务器用户名
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),               // 生产环境建议使用固定HostKey验证
        Auth: []ssh.AuthMethod{publicKeyAuth("id_rsa")}, // 私钥路径，按需修改
    }

    // 建立 SSH 连接
    conn, err := ssh.Dial("tcp", "192.168.6.102:22", config)
    if err != nil {
        panic(fmt.Errorf("SSH 连接失败: %v", err))
    }
    defer conn.Close()

    // 创建 SFTP 客户端
    client, err := sftp.NewClient(conn)
    if err != nil {
        panic(fmt.Errorf("创建 SFTP 客户端失败: %v", err))
    }
    defer client.Close()

    // 示例：上传文件
    srcFile, err := os.Open("tt.txt")
    if err != nil {
        panic(fmt.Errorf("打开本地文件失败: %v", err))
    }
    defer srcFile.Close()

    // dstFile, err := client.Create("/Z2011551000016/tt.txt")
    // // dstFile, err := client.Create("tt.txt")
    // // dstFile, err := client.Create("/tt.txt")
    // if err != nil {
    //     panic(fmt.Errorf("在远端创建文件失败: %v", err))
    // }
    // defer dstFile.Close()

    // _, err = io.Copy(dstFile, srcFile)
    // if err != nil {
    //     panic(fmt.Errorf("复制文件内容失败: %v", err))
    // }

    // fmt.Println("文件上传成功！")
    // 下载文件（指定为用户请求的 abc）
    remotePath := "/Z2011551000016/dst/20260610/abc"   // 远程文件路径
    localPath := "./abc"       // 本地保存路径

	if err := downloadFile(client, remotePath, localPath); err != nil {
		fmt.Printf("下载失败: %v", err)
	}

	fmt.Printf("文件下载成功: %s -> %s\n", remotePath, localPath)
}

// ====================== 公钥认证辅助函数 ======================

// 无密码私钥
func publicKeyAuth(keyPath string) ssh.AuthMethod {
    key, err := os.ReadFile(keyPath)
    if err != nil {
        panic(fmt.Sprintf("无法读取私钥文件: %v", err))
    }

    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        panic(fmt.Sprintf("解析私钥失败: %v", err))
    }

    return ssh.PublicKeys(signer)
}

// 有密码私钥（加密的私钥）
func publicKeyAuthWithPassphrase(keyPath string, passphrase string) ssh.AuthMethod {
    key, err := os.ReadFile(keyPath)
    if err != nil {
        panic(fmt.Sprintf("无法读取私钥文件: %v", err))
    }

    signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
    if err != nil {
        panic(fmt.Sprintf("解析加密私钥失败: %v", err))
    }

    return ssh.PublicKeys(signer)
}

// 下载文件函数
func downloadFile(client *sftp.Client, remotePath, localPath string) error {
	// 打开远程文件
	srcFile, err := client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("打开远程文件失败: %w", err)
	}
	defer srcFile.Close()

	// 创建本地文件
	dstFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %w", err)
	}
	defer dstFile.Close()

	// 复制内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("复制文件内容失败: %w", err)
	}

	return nil
}