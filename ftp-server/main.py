from asyncio import protocols
from pyftpdlib.authorizers import DummyAuthorizer
from pyftpdlib.handlers import FTPHandler, TLS_FTPHandler
from pyftpdlib.servers import FTPServer
import logging
import os 
# 用户权限配置

default_date_readwrite = "./date_rw"
default_date_readonly = "./date_r"
default_date_anonymous = "./date_anonymous"


def ensure_directory_exists(directory_path):
    """
    检查文件夹是否存在，如果不存在则创建它
    :param directory_path: 文件夹路径
    """
    if not os.path.exists(directory_path):
        os.makedirs(directory_path)
        print(f"文件夹 {directory_path} 已创建")
    else:
        print(f"文件夹 {directory_path} 已存在")
        
        

def setup_authorizer():
    authorizer = DummyAuthorizer()
    # 添加用户（用户名，密码，目录，权限）
    date_rw = default_date_readwrite
    date_r = default_date_readonly
    date_anonymous = default_date_anonymous
    # 文件夹不存在则创建
    ensure_directory_exists(date_rw)
    ensure_directory_exists(date_r)
    ensure_directory_exists(date_anonymous)
    # 创建用户
    authorizer.add_user("rw", "123",
                        date_rw, perm="elradfmw")
    authorizer.add_user("readonly", "password123", date_r, perm="elr")
    authorizer.add_anonymous(date_anonymous, perm="elr")  # 匿名用户
    return authorizer

# 日志配置


def setup_logging():
    logging.basicConfig(
        filename="ftp_server.log",
        level=logging.DEBUG,
        format="%(asctime)s - %(levelname)s - %(message)s",
    )

# 自定义FTP处理器


class CustomFTPHandler(FTPHandler):
    def on_connect(self):
        logging.info(f"Client connected: {self.remote_ip}:{self.remote_port}")

    def on_disconnect(self):
        logging.info(
            f"Client disconnected: {self.remote_ip}:{self.remote_port}")

    def on_login(self, username):
        logging.info(f"User logged in: {username}")

    def on_logout(self, username):
        logging.info(f"User logged out: {username}")

    def on_file_sent(self, file):
        logging.info(f"File sent: {file}")

    def on_file_received(self, file):
        logging.info(f"File received: {file}")

    def on_incomplete_file_sent(self, file):
        logging.warning(f"Incomplete file sent: {file}")

    def on_incomplete_file_received(self, file):
        logging.warning(f"Incomplete file received: {file}")

# 支持TLS的FTP处理器


class CustomTLSFTPHandler(TLS_FTPHandler):
    def on_connect(self):
        logging.info(
            f"TLS Client connected: {self.remote_ip}:{self.remote_port}")

    def on_disconnect(self):
        logging.info(
            f"TLS Client disconnected: {self.remote_ip}:{self.remote_port}")

    def on_login(self, username):
        logging.info(f"TLS User logged in: {username}")

    def on_logout(self, username):
        logging.info(f"TLS User logged out: {username}")

    def on_file_sent(self, file):
        logging.info(f"TLS File sent: {file}")

    def on_file_received(self, file):
        logging.info(f"TLS File received: {file}")

# 配置FTP服务器


def run_ftp_server(use_tls=False):
    authorizer = setup_authorizer()
    setup_logging()

    if use_tls:
        handler = CustomTLSFTPHandler
        handler.certfile = "cert.pem"  # 替换为您的TLS证书路径
        handler.keyfile = "key.pem"    # 替换为您的TLS密钥路径
    else:
        handler = CustomFTPHandler

    handler.authorizer = authorizer
    handler.passive_ports = range(60000, 65535)  # 配置被动模式端口范围
    handler.banner = "Welcome to the secure FTP server!"  # 欢迎消息

    # 启动服务器
    server = FTPServer(("0.0.0.0", 2121), handler)
    server.max_cons = 256  # 最大连接数
    server.max_cons_per_ip = 5  # 每IP最大连接数
    logging.info("Starting FTP server...")
    print("Starting FTP server...")
    server.serve_forever()


if __name__ == "__main__":
    # run_ftp_server(use_tls=True)  # 启用TLS
    run_ftp_server()
