# 🔒 JavaScript 代码混淆工具

基于 **WebAssembly + TinyGo** 实现的高性能 JavaScript 代码混淆器，提供多种混淆策略和优秀的用户体验。

## ✨ 特性

### 🚀 核心功能
- **🔤 标识符混淆** - 将变量名、函数名替换为随机短字符
- **🔐 字符串加密** - 支持 Base64、十六进制、Unicode 多种加密方式
- **🌀 控制流平坦化** - 打乱代码执行流程，增加逆向难度
- **📦 代码压缩** - 移除空格、注释，减小文件体积
- **⚙️ 灵活配置** - 可选择性启用各种混淆策略

### 🎯 技术优势
- **⚡ 极致性能** - TinyGo 编译的 WASM，体积极小，适用于 CloudFlare Free Plan 用户
- **🌐 客户端处理** - 无服务器运算，保护代码隐私
- **☁️ 部署友好** - 静态文件，完美支持 Cloudflare Worker

## 🏗️ 技术架构

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   前端界面      │    │   TinyGo 运行时   │    │  WebAssembly    │
│                 │    │                  │    │                 │
│ HTML5 + CSS3    │◄──►│  wasm_exec.js    │◄──►│  obfuscator.wasm│
│ + JavaScript    │    │  (官方运行时)     │    │   (1.37 MB)     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## 🛠️ 开发环境

### 前置要求
- **Node.js** >= 16.0.0
- **pnpm** >= 8.0.0
- **Go** >= 1.21.0
- **TinyGo** >= 0.30.0

### 安装依赖

```bash
# 安装 Node.js 依赖
pnpm install

# 安装 TinyGo (macOS)
brew install tinygo

# 安装 TinyGo (Linux)
wget https://github.com/tinygo-org/tinygo/releases/download/v0.33.0/tinygo_0.33.0_amd64.deb
sudo dpkg -i tinygo_0.33.0_amd64.deb
```

## 🚀 快速开始

### 本地开发

```bash
# 克隆项目
git clone https://github.com/zaunist/js-obfuscator.git
cd js-obfuscator

# 安装依赖
pnpm install

# 构建项目
pnpm run build

# 启动本地服务器
pnpm run dev
```

### Cloudflare Workers 部署

```bash
# 构建生产版本
pnpm run build

# 部署到 Cloudflare Pages
pnpm run deploy
```

## 📁 项目结构

```
js-obfuscator/
├── frontend/                 # 前端源码
│   ├── app.js               # 主应用逻辑
│   ├── styles.css           # 样式文件
│   └── index.html           # HTML 模板
├── wasm/                    # WebAssembly 源码
│   ├── main.go              # Go 主程序
│   ├── go.mod               # Go 模块配置
├── dist/                    # 构建输出
│   ├── wasm/
│   │   └── obfuscator.wasm  # TinyGo 编译的 WASM
│   ├── bundle.js            # 前端 JS 包
│   ├── styles.css           # 样式文件
│   └── index.html           # 入口页面
├── webpack.config.js        # Webpack 配置
├── wrangler.toml           # Cloudflare 配置
└── package.json            # 项目配置
```

## 🔧 构建脚本

```json
{
  "scripts": {
    "build": "pnpm run build:wasm && pnpm run build:frontend",
    "build:wasm": "mkdir -p dist/wasm && cd wasm && tinygo build -o ../dist/wasm/obfuscator.wasm -target wasm .",
    "build:frontend": "webpack --mode production --output-path dist",
    "dev": "wrangler dev",
    "deploy": "pnpm run deploy"
  }
}
```

## 🎨 功能演示

### 混淆前
```javascript
function calculateSum(numbers) {
    let total = 0;
    for (let i = 0; i < numbers.length; i++) {
        total += numbers[i];
    }
    return total;
}

const result = calculateSum([1, 2, 3, 4, 5]);
console.log("Result:", result);
```

### 混淆后
```javascript
function _0x1a2b(_0x3c4d){let _0x5e6f=0x0;for(let _0x7g8h=0x0;_0x7g8h<_0x3c4d['\x6c\x65\x6e\x67\x74\x68'];_0x7g8h++){_0x5e6f+=_0x3c4d[_0x7g8h];}return _0x5e6f;}const _0x9i0j=_0x1a2b([0x1,0x2,0x3,0x4,0x5]);console['\x6c\x6f\x67']('\x52\x65\x73\x75\x6c\x74\x3a',_0x9i0j);
```

## 🔐 混淆策略详解

### 1. 标识符混淆
- 将变量名、函数名替换为随机生成的短字符
- 保持代码功能不变的同时增加阅读难度
- 支持保留关键字和内置对象

### 2. 字符串加密
- **Base64 编码**: 将字符串转换为 Base64 格式
- **十六进制编码**: 使用 `\x` 转义序列
- **Unicode 编码**: 使用 `\u` 转义序列
- **动态解密**: 运行时自动解密字符串

### 3. 控制流平坦化
- 将线性代码转换为状态机结构
- 使用 switch-case 语句打乱执行顺序
- 增加静态分析和逆向工程难度

### 4. 代码压缩
- 移除所有空白字符和换行符
- 删除注释和无用代码
- 优化代码结构减小文件体积

## 🌐 部署配置

### Cloudflare Worker

```toml
name = "js-obfuscator"
main = "src/worker.js"
compatibility_date = "2025-09-01"

# Assets 配置
[assets]
directory = "./dist"
binding = "ASSETS"
```

## 🤝 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 WTFPL 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🔗 相关链接

- **问题反馈**: [https://github.com/zaunist/js-obfuscator/issues](https://github.com/zaunist/js-obfuscator/issues)
- **个人博客**: [https://ajie.lu](https://ajie.lu)

---

**⭐ 如果这个项目对你有帮助，请给个 Star 支持一下！**