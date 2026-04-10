# Sealos 部署指南

## 部署步骤

### 1. 注册账号
打开 https://cloud.sealos.io
- 使用 GitHub 账号登录

### 2. 创建应用
1. 点击「应用管理」→「创建应用」
2. 填写配置：
   - 应用名称：`location-tracking`
   - CPU：0.5核
   - 内存：512MB
   - 端口：8080

### 3. 部署方式选择

**方式A：Docker Hub 镜像**
1. 在本地构建镜像：
```bash
cd location-tracking-shortlink
docker build -t yourusername/location-tracking:latest .
docker push yourusername/location-tracking:latest
```
2. 在 Sealos 填入镜像地址：`yourusername/location-tracking:latest`

**方式B：直接上传代码**
1. 把项目打包 zip 上传
2. Sealos 会自动识别 Dockerfile 部署

### 4. 访问
部署成功后，访问 Sealos 提供的公网地址：
```
http://xxx.sealoshub.com:8080/admin
```

---

## 项目文件说明

- `Dockerfile` - 容器构建文件
- `sealos.yaml` - Sealos 部署配置
- `config.yaml` - 应用配置
- `track` - Linux 可执行文件（已编译）
