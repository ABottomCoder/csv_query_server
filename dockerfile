# 使用适当的基础镜像
FROM golang:latest

# 创建目录
RUN mkdir -p /app/data

# 修改文件权限
RUN chmod +x /app/data

# 设置工作目录
WORKDIR /app

# 将后端服务的源代码复制到镜像中
COPY / .

# 构建后端服务
RUN go build -o csv_query_server .

# 暴露端口
EXPOSE 9527
EXPOSE 7259

# 设置CSV文件路径
ENV CSV_FILE_PATH /app/data/data.csv

# 运行后端服务
CMD ["./csv_query_server", "-f", "$CSV_FILE_PATH"]
