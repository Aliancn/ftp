package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// FTPConn 封装FTP客户端的核心功能
type FTPConn struct {
	controlConn net.Conn // 控制连接
	dataConn    net.Conn // 数据连接
	reader      *bufio.Reader
	// downloadOffset map[string]int64
}

// NewFTPConn 初始化FTP客户端
func NewFTPConn() *FTPConn {
	return &FTPConn{}
}

// Connect 连接到FTP服务器
func (ftp *FTPConn) Dial(serverAddr string) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("连接到FTP服务器失败: %v", err)
	}
	ftp.controlConn = conn

	ftp.reader = bufio.NewReader(conn)

	// 读取初始响应
	response, err := ftp.readResponse()
	if err != nil {
		return fmt.Errorf("读取初始响应失败: %v", err)
	}
	// MyLogger.Info("初始响应:", response)

	if !strings.HasPrefix(response, "220") {
		return fmt.Errorf("服务器未准备好: %s", response)
	}

	// MyLogger.Info("dial 成功连接到服务器")
	return nil
}

// readResponse reads a response from the server
func (ftp *FTPConn) readResponse() (string, error) {
	response, err := ftp.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}
	return strings.TrimSpace(response), nil
}

// SendCommand 向服务器发送命令，并接收响应
func (ftp *FTPConn) SendCommand(command string) (string, error) {
	_, err := ftp.controlConn.Write([]byte(command + "\r\n"))
	if err != nil {
		return "", fmt.Errorf("发送命令失败: %v", err)
	}

	// 读取服务器响应
	response, err := ftp.readResponse()
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}
	return strings.TrimSpace(response), nil
}

// Close 关闭控制连接和数据连接
func (ftp *FTPConn) Close() error {
	if ftp.controlConn != nil {
		ftp.controlConn.Close()
		res, _ := ftp.readResponse()
		MyLogger.Info("关闭控制连接:", res)
	}
	if ftp.dataConn != nil {
		ftp.dataConn.Close()
		res, _ := ftp.readResponse()
		MyLogger.Info("关闭数据连接:", res)
	}
	// MyLogger.Info("成功关闭连接")
	return nil
}

// Login 登录FTP服务器
func (ftp *FTPConn) Login(username, password string) error {

	response, err := ftp.SendCommand(fmt.Sprintf("USER %s", username))
	if err != nil {
		return err
	}
	MyLogger.Info("USER 服务器响应:", response)

	response, err = ftp.SendCommand(fmt.Sprintf("PASS %s", password))
	if err != nil {
		return err
	}
	MyLogger.Info("PASS 服务器响应:", response)

	if !strings.HasPrefix(response, "230") {
		return fmt.Errorf("登录失败: %s", response)
	}
	return nil
}

// ListFiles 获取指定路径下的文件列表
func (ftp *FTPConn) ListFiles(path string) ([]string, error) {
	// 建立数据连接
	dataConn, err := ftp.establishDataConn()
	if err != nil {
		return nil, err
	}
	ftp.dataConn = dataConn
	defer ftp.closeDataConn() // 确保数据连接关闭

	// 发送LIST命令，附加路径参数
	command := fmt.Sprintf("LIST %s", path)
	_, err = ftp.SendCommand(command)
	if err != nil {
		return nil, fmt.Errorf("发送LIST命令失败: %v", err)
	}

	// 从数据连接读取文件列表
	var files []string
	scanner := bufio.NewScanner(ftp.dataConn)
	MyLogger.Info(fmt.Sprintf("目录 '%s' 下的文件列表:", path))
	for scanner.Scan() {
		line := scanner.Text()
		// MyLogger.Info(line)
		files = append(files, line)
	}

	// 检查扫描是否出错
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件列表时出错: %v", err)
	}

	return files, nil
}

// parsePASVResponse 解析PASV响应，返回数据连接的地址
func parsePASVResponse(response string) (string, error) {
	start := strings.Index(response, "(")
	end := strings.Index(response, ")")
	if start == -1 || end == -1 {
		MyLogger.Info(response)
		return "", fmt.Errorf("无效的PASV响应: %s", response)
	}

	// 提取地址和端口部分
	parts := strings.Split(response[start+1:end], ",")
	if len(parts) != 6 {
		return "", fmt.Errorf("无效的PASV地址格式: %s", response)
	}

	// 解析IP和端口
	ip := strings.Join(parts[:4], ".")
	portPart1 := parseInt(parts[4])
	portPart2 := parseInt(parts[5])
	port := portPart1*256 + portPart2

	return fmt.Sprintf("%s:%d", ip, port), nil
}

// parseInt 将字符串转换为整数
func parseInt(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

// establishDataConn establishes a data connection using PASV mode
func (ftp *FTPConn) establishDataConn() (net.Conn, error) {
	if ftp.dataConn != nil {
		ftp.dataConn.Close()
		ftp.dataConn = nil
	}
	response, err := ftp.SendCommand("PASV")
	if err != nil {
		return nil, fmt.Errorf("发送PASV命令失败: %v", err)
	}
	MyLogger.Info("PASV 服务器响应:", response)

	dataAddr, err := parsePASVResponse(response)
	if err != nil {
		return nil, fmt.Errorf("解析PASV响应失败: %v", err)
	}

	dataConn, err := net.Dial("tcp", dataAddr)
	if err != nil {
		return nil, fmt.Errorf("建立数据连接失败: %v", err)
	}

	fmt.Println("成功建立数据连接 ", dataAddr)
	return dataConn, nil
}

// closeDataConn
func (ftp *FTPConn) closeDataConn() error {
	if ftp.dataConn != nil {
		ftp.dataConn.Close()
		ftp.dataConn = nil
	}
	res, _ := ftp.readResponse() // 读取服务器响应
	MyLogger.Info("关闭数据连接:", res)
	return nil
}

// UploadFile 上传文件到服务器
func (ftp *FTPConn) STOR(localPath, remotePath string) error {
	// 切换到被动模式
	dataconn, err := ftp.establishDataConn()
	if err != nil {
		return err
	}
	ftp.dataConn = dataconn
	defer ftp.closeDataConn() // 确保数据连接关闭

	// 发送STOR命令
	_, err = ftp.SendCommand(fmt.Sprintf("STOR %s", remotePath))
	if err != nil {
		return err
	}

	// 打开本地文件并发送数据
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(ftp.dataConn, file)
	if err != nil {
		return fmt.Errorf("文件上传失败: %v", err)
	}
	return nil
}

// DownloadFile 从服务器下载文件
func (ftp *FTPConn) RETR(remotePath, localPath string) error {
	// 切换到被动模式
	dataConn, err := ftp.establishDataConn()
	if err != nil {
		return err
	}
	ftp.dataConn = dataConn
	defer ftp.closeDataConn() // 确保数据连接关闭

	// 发送RETR命令
	_, err = ftp.SendCommand(fmt.Sprintf("RETR %s", remotePath))
	if err != nil {
		return err
	}

	// 保存接收到的数据到本地文件
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, ftp.dataConn)
	if err != nil {
		return fmt.Errorf("文件下载失败: %v", err)
	}
	return nil
}

// MakeDirectory 在FTP服务器上创建新文件夹
func (ftp *FTPConn) MakeDir(directoryName string) error {
	// 发送MKD命令到服务器
	response, err := ftp.SendCommand(fmt.Sprintf("MKD %s", directoryName))
	if err != nil {
		return fmt.Errorf("发送MKD命令失败: %v", err)
	}
	MyLogger.Info("MKD 服务器响应:", response)

	// 检查响应是否成功 (通常是以"257"开头的响应)
	if !strings.HasPrefix(response, "257") {
		return fmt.Errorf("创建文件夹失败: %s", response)
	}
	MyLogger.Info("文件夹 '%s' 创建成功\n", directoryName)

	return nil
}

// Remove 删除FTP服务器上的文件或文件夹
func (ftp *FTPConn) Dele(target string, isDirectory bool) error {
	var command string

	// 根据类型选择命令
	if isDirectory {
		command = fmt.Sprintf("RMD %s", target) // 删除文件夹
	} else {
		command = fmt.Sprintf("DELE %s", target) // 删除文件
	}

	// 发送命令到服务器
	response, err := ftp.SendCommand(command)
	if err != nil {
		return fmt.Errorf("发送删除命令失败: %v", err)
	}
	MyLogger.Info("DELE 服务器响应:", response)

	// 检查响应是否成功
	if !strings.HasPrefix(response, "250") { // 250 表示删除成功
		return fmt.Errorf("删除失败: %s", response)
	}

	return nil
}

func (ftp *FTPConn) SetBinaryMode() error {
	response, err := ftp.SendCommand("TYPE I")
	if err != nil {
		return fmt.Errorf("设置二进制模式失败: %v", err)
	}
	MyLogger.Info("TYPE 服务器响应:", response)

	if !strings.HasPrefix(response, "200") {
		return fmt.Errorf("设置二进制模式失败: %s", response)
	}
	return nil
}

type Progress struct {
	FileName   string `json:"fileName"`
	Downloaded int64  `json:"downloaded"`
	TotalSize  int64  `json:"totalSize"`
}

// ResumeDownload 恢复下载文件 TODO
func (ftp *FTPClient) REST_RETR(remoteFile, localFile string, offset int64, c context.Context) error {

	// 打开本地文件，准备写入（从断点开始）
	file, err := os.OpenFile(localFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("无法打开本地文件: %v", err)
	}

	// 设置文件写入的偏移量
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("设置本地文件偏移量失败: %v", err)
	}

	// change to binary mode
	if err := ftp.SetBinaryMode(); err != nil {
		return fmt.Errorf("设置二进制模式失败: %v", err)
	}

	// 发送 REST 命令指定恢复点
	response, err := ftp.SendCommand(fmt.Sprintf("REST %d", offset))
	if err != nil || !strings.HasPrefix(response, "350") {
		return fmt.Errorf("REST 命令失败: %v", response)
	}

	// 发送 RETR 命令开始下载文件
	response, err = ftp.SendCommand(fmt.Sprintf("RETR %s", remoteFile))
	if err != nil || !strings.HasPrefix(response, "150") {
		return fmt.Errorf("RETR 命令失败: %v", response)
	}

	// 建立数据连接
	dataConn, err := ftp.establishDataConn()
	if err != nil {
		return fmt.Errorf("数据连接建立失败: %v", err)
	}
	ftp.dataConn = dataConn
	defer ftp.dataConn.Close() // 确保数据连接关闭

	// 下载文件并反馈进度
	buf := make([]byte, 4*1024*1024) // 每次读取
	downloaded := offset             // 已下载的字节数，从偏移量开始
	runtime.EventsEmit(c, "download-progress", Progress{
		FileName:   remoteFile,
		Downloaded: downloaded,
	})
	lastUpdateTime := time.Now()
	complated := false

	defer func() {
		// 如果任务被取消，清理部分下载的文件
		if ftp.ctx.Err() != nil {
			// change to ascii mode
			if err := ftp.SetAsciiMode(); err != nil {
				fmt.Println("设置ASCII模式失败: ", err)
			}
			err = file.Close()
			if err != nil {
				fmt.Println("关闭文件失败: ", err)
			}
			ftp.dataConn.Close()
			response, err = ftp.readResponse()
			if err != nil {
				fmt.Println("读取服务器响应失败: ", err)
			}
			fmt.Println("下载任务暂停，终止dataConn", response)
		}
	}()

	for {
		select {
		case <-ftp.ctx.Done():
			fmt.Println("下载任务被取消")
			return fmt.Errorf("下载被取消: %v", ftp.ctx.Err())
		default:
			// 从数据连接中读取数据
			n, readErr := ftp.dataConn.Read(buf)
			if n > 0 {
				// 写入本地文件
				if _, writeErr := file.Write(buf[:n]); writeErr != nil {
					return fmt.Errorf("写入文件失败: %v", writeErr)
				}

				// 更新已下载的字节数
				downloaded += int64(n)
				time.Sleep(2 * time.Second)

				// 限制触发进度事件的频率
				if downloaded%1024 == 0 || time.Since(lastUpdateTime) > time.Second {
					runtime.EventsEmit(c, "download-progress", Progress{
						FileName:   remoteFile,
						Downloaded: downloaded,
					})
					lastUpdateTime = time.Now()
				}
			}

			if readErr != nil {
				if readErr == io.EOF {
					fmt.Println("文件读取完成")
					runtime.EventsEmit(c, "download-progress", Progress{
						FileName:   remoteFile,
						Downloaded: downloaded,
					})
					complated = true
					// 检查服务器返回的结束状态码
					response, err = ftp.readResponse()
					if err != nil || !strings.HasPrefix(response, "226") {
						return fmt.Errorf("下载未正确完成: %v", response)
					}
					// change to ascii mode
					if err := ftp.SetAsciiMode(); err != nil {
						return err
					}
					err = file.Close()
					if err != nil {
						return fmt.Errorf("关闭文件失败: %v", err)
					}
					fmt.Println("文件下载完成", remoteFile)
					ftp.validateDataConn()
					break
				}
				return fmt.Errorf("读取数据失败: %v", readErr)
			}
		}
		if complated {
			break
		}
	}

	return nil
}

func GetDownloadedOffset(localFile string) (int64, error) {
	// 检查本地文件是否存在
	fileInfo, err := os.Stat(localFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，说明尚未开始下载
			return 0, nil
		}
		return 0, fmt.Errorf("无法获取本地文件信息: %v", err)
	}

	// 返回文件的大小作为偏移量
	fmt.Println("文件名", fileInfo.Name(), "文件大小", fileInfo.Size())
	return fileInfo.Size(), nil
}

// ResumeUpload 恢复上传文件
// func (ftp *FTPClient) REST_STOR(localFile, remoteFile string) error {
// 	// 打开本地文件，准备读取
// 	file, err := os.Open(localFile)
// 	if err != nil {
// 		return fmt.Errorf("无法打开本地文件: %v", err)
// 	}
// 	defer file.Close()

// 	// 检查服务器上已存在文件的大小
// 	sizeResponse, err := ftp.SendCommand(fmt.Sprintf("SIZE %s", remoteFile))
// 	if err != nil || !strings.HasPrefix(sizeResponse, "213") {
// 		return fmt.Errorf("无法获取远程文件大小: %v", sizeResponse)
// 	}

// 	// 解析服务器返回的文件大小
// 	var offset int64
// 	fmt.Sscanf(sizeResponse, "213 %d", &offset)

// 	// 设置文件读取的偏移量
// 	if _, err := file.Seek(offset, 0); err != nil {
// 		return fmt.Errorf("设置本地文件偏移量失败: %v", err)
// 	}

// 	// 发送 REST 命令指定恢复点
// 	response, err := ftp.SendCommand(fmt.Sprintf("REST %d", offset))
// 	if err != nil || !strings.HasPrefix(response, "350") {
// 		return fmt.Errorf("REST 命令失败: %v", response)
// 	}

// 	// 发送 STOR 命令开始上传文件
// 	response, err = ftp.SendCommand(fmt.Sprintf("STOR %s", remoteFile))
// 	if err != nil || !strings.HasPrefix(response, "150") {
// 		return fmt.Errorf("STOR 命令失败: %v", response)
// 	}

// 	// 建立数据连接
// 	dataConn, err := ftp.establishDataConn()
// 	if err != nil {
// 		return fmt.Errorf("数据连接建立失败: %v", err)
// 	}
// 	ftp.dataConn = dataConn
// 	// defer ftp.closeDataConn() // 确保数据连接关闭
// 	defer ftp.dataConn.Close()

// 	// 从本地文件读取数据并写入数据连接
// 	_, err = io.Copy(dataConn, file)
// 	if err != nil {
// 		return fmt.Errorf("文件上传失败: %v", err)
// 	}

// 	// 检查服务器返回的结束状态码
// 	response, err = ftp.readResponse() // 读取服务器响应
// 	if err != nil || !strings.HasPrefix(response, "226") {
// 		return fmt.Errorf("上传未正确完成: %v", response)
// 	}

// 	fmt.Println("文件上传完成")
// 	return nil
// }

// SetAsciiMode sets the FTP transfer mode to ASCII
func (ftp *FTPClient) SetAsciiMode() error {
	response, err := ftp.SendCommand("TYPE A")
	if err != nil {
		return fmt.Errorf("设置ASCII模式失败: %v", err)
	}
	MyLogger.Info("TYPE A服务器响应:", response)

	if !strings.HasPrefix(response, "200") {
		return fmt.Errorf("设置ASCII模式失败: %s", response)
	}
	return nil
}

func (ftp *FTPClient) validateDataConn() error {
	if ftp.dataConn == nil {
		return fmt.Errorf("数据连接已经关闭")
	}

	// 发送一个简单的 ping 测试数据包
	_, err := ftp.dataConn.Write([]byte{})
	if err != nil {
		return fmt.Errorf("数据连接不可用: %v", err)
	}
	fmt.Println("数据连接正常")
	return nil
}
