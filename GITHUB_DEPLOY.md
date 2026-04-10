# GitHub + Sealos 自动化部署

## 步骤 1：在 GitHub 创建仓库

1. 打开 https://github.com/new
2. 仓库名称：`location-tracking`
3. 选择 Private（私有）
4. 点击 Create repository

## 步骤 2：推送代码到 GitHub

```bash
cd location-tracking-shortlink
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/你的用户名/location-tracking.git
git push -u origin main
```

## 步骤 3：连接 Sealos

1. 打开 https://cloud.sealos.io
2. 点击「应用管理」→「创建应用」
3. 选择「从 GitHub 导入」
4. 授权 Sealos 访问你的 GitHub
5. 选择 `location-tracking` 仓库
6. Sealos 自动检测 Dockerfile 并部署

## 步骤 4：以后更新代码

修改代码后：
```bash
git add .
git commit -m "Update feature"
git push
```

Sealos 会自动检测到更新并重新部署！

---

## 需要的材料

1. GitHub 账号
2. Sealos 账号（用 GitHub 登录）

---

## 项目文件

- `Dockerfile` - 容器构建配置
- `.github/workflows/deploy.yml` - GitHub Actions 自动部署（可选）
- `config.yaml` - 应用配置
- `templates/` - 前端页面
- `static/` - 静态资源
