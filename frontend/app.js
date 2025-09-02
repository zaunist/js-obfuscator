// JavaScript 混淆工具主应用
class JSObfuscatorApp {
  constructor() {
    this.wasmModule = null;
    this.isReady = false;
    this.initializeApp();
  }

  async initializeApp() {
    console.log("开始初始化应用...");

    // 初始化 UI 事件
    this.initializeEventListeners();

    // 加载 WASM 模块
    await this.loadWASM();

    // 更新状态
    this.updateWASMStatus();

    // 加载示例代码
    this.loadExampleCode();
  }

  // 初始化事件监听器
  initializeEventListeners() {
    // 混淆按钮
    document.getElementById("obfuscateBtn").addEventListener("click", () => {
      this.obfuscateCode();
    });

    // 验证功能已移除

    // 复制结果按钮
    document.getElementById("copyResult").addEventListener("click", () => {
      this.copyResult();
    });

    // 下载结果按钮
    document.getElementById("downloadResult").addEventListener("click", () => {
      this.downloadResult();
    });

    // 加载示例按钮
    document.getElementById("loadExample").addEventListener("click", () => {
      this.loadExampleCode();
    });

    // 清空输入按钮
    document.getElementById("clearInput").addEventListener("click", () => {
      this.clearInput();
    });

    // 输入代码变化监听
    document.getElementById("inputCode").addEventListener("input", () => {
      this.updateInputStats();
    });

    // 配置变化监听
    const configInputs = document.querySelectorAll(
      '.config-item input[type="checkbox"]'
    );
    configInputs.forEach((input) => {
      input.addEventListener("change", () => {
        this.saveConfig();
      });
    });

    // 加载保存的配置
    this.loadConfig();
  }

  async loadWASM() {
    try {
      console.log("开始加载 WASM...");

      // 使用 TinyGo 官方运行时
      if (!window.Go) {
        throw new Error("TinyGo 运行时未加载");
      }

      const go = new Go();

      console.log("开始获取 WASM 文件...");
      const wasmResponse = await fetch("/wasm/obfuscator.wasm");
      if (!wasmResponse.ok) {
        throw new Error(`WASM 文件获取失败: ${wasmResponse.status}`);
      }

      console.log("WASM 文件获取成功，开始实例化...");
      const wasmBytes = await wasmResponse.arrayBuffer();
      const result = await WebAssembly.instantiate(wasmBytes, go.importObject);

      this.wasmInstance = result.instance;
      console.log("WASM 实例化成功，启动程序...");

      // 启动 Go 程序
      go.run(this.wasmInstance);

      console.log("等待 WASM 函数注册...");
      await this.waitForWASMFunctions();

      this.isReady = true;
      console.log("WASM 模块加载成功！");

      // 测试函数
      if (window.wasmTest) {
        const result = window.wasmTest();
        console.log("测试结果:", result);
      }
    } catch (error) {
      console.error("WASM 加载失败:", error);
      this.showError("WASM 模块加载失败: " + error.message);
    }
  }

  // 创建 WASI 导入对象
  createWASIImports() {
    const textEncoder = new TextEncoder();
    const textDecoder = new TextDecoder();
    let memory;

    return {
      wasi_snapshot_preview1: {
        // 基本的 WASI 函数实现
        proc_exit: (code) => {
          console.log(`WASI proc_exit called with code: ${code}`);
        },
        fd_write: (fd, iovs, iovs_len, nwritten) => {
          // 简单的标准输出实现
          if (!memory) return 0;

          let written = 0;
          const view = new DataView(memory.buffer);

          for (let i = 0; i < iovs_len; i++) {
            const ptr = iovs + i * 8;
            const buf = view.getUint32(ptr, true);
            const bufLen = view.getUint32(ptr + 4, true);

            const data = new Uint8Array(memory.buffer, buf, bufLen);
            const text = textDecoder.decode(data);
            console.log(text);
            written += bufLen;
          }

          if (nwritten !== 0) {
            view.setUint32(nwritten, written, true);
          }

          return 0;
        },
        fd_read: () => 0,
        fd_seek: () => 0,
        fd_close: () => 0,
        path_open: () => 0,
        environ_sizes_get: () => 0,
        environ_get: () => 0,
        args_sizes_get: (argc, argv_buf_size) => {
          if (!memory) return 0;
          const view = new DataView(memory.buffer);
          view.setUint32(argc, 0, true);
          view.setUint32(argv_buf_size, 0, true);
          return 0;
        },
        args_get: () => 0,
        random_get: (buf, buf_len) => {
          if (!memory) return 0;
          const randomBytes = new Uint8Array(memory.buffer, buf, buf_len);
          crypto.getRandomValues(randomBytes);
          return 0;
        },
        clock_time_get: (id, precision, time) => {
          if (!memory) return 0;
          const view = new DataView(memory.buffer);
          const now = BigInt(Date.now() * 1000000);
          view.setBigUint64(time, now, true);
          return 0;
        },
        fd_read: (fd, iovs, iovs_len, nread) => {
          return 0;
        },
        fd_close: (fd) => {
          return 0;
        },
        fd_seek: (fd, offset_low, offset_high, whence, newoffset) => {
          return 0;
        },
        environ_sizes_get: (environ_count, environ_buf_size) => {
          return 0;
        },
        environ_get: (environ, environ_buf) => {
          return 0;
        },
        args_sizes_get: (argc, argv_buf_size) => {
          return 0;
        },
        args_get: (argv, argv_buf) => {
          return 0;
        },
        random_get: (buf, buf_len) => {
          // 填充随机数据
          const memory = new Uint8Array(this.wasmMemory.buffer);
          for (let i = 0; i < buf_len; i++) {
            memory[buf + i] = Math.floor(Math.random() * 256);
          }
          return 0;
        },
        clock_time_get: (id, precision, time) => {
          // 返回当前时间戳
          const now = BigInt(Date.now() * 1000000);
          const memory = new DataView(this.wasmMemory.buffer);
          memory.setBigUint64(time, now, true);
          return 0;
        },
      },
      // TinyGo 需要的 gojs 模块
      gojs: {
        // 基本的 gojs 函数实现
        "runtime.wasmExit": (code) => {
          console.log(`Go program exited with code: ${code}`);
        },
        "runtime.wasmWrite": (fd, ptr, len) => {
          if (!this.wasmInstance || !this.wasmInstance.exports.memory) return 0;

          const memory = new Uint8Array(
            this.wasmInstance.exports.memory.buffer
          );
          const data = memory.slice(ptr, ptr + len);
          const text = new TextDecoder().decode(data);

          if (fd === 1) {
            // stdout
            console.log(text);
          } else if (fd === 2) {
            // stderr
            console.error(text);
          }

          return len;
        },
        "runtime.nanotime": () => {
          return BigInt(Date.now() * 1000000);
        },
        "runtime.walltime": () => {
          const now = Date.now();
          return [Math.floor(now / 1000), (now % 1000) * 1000000];
        },
        "runtime.ticks": () => {
          return BigInt(performance.now() * 1000000);
        },
        "runtime.scheduleTimeoutEvent": (delay) => {
          // 简单的超时事件调度
          setTimeout(() => {
            if (this.wasmInstance && this.wasmInstance.exports.go_scheduler) {
              this.wasmInstance.exports.go_scheduler();
            }
          }, delay / 1000000); // 转换纳秒到毫秒
        },
        "syscall/js.valueGet": () => 0,
        "syscall/js.valueSet": () => {},
        "syscall/js.valueDelete": () => {},
        "syscall/js.valueIndex": () => 0,
        "syscall/js.valueSetIndex": () => {},
        "syscall/js.valueCall": () => 0,
        "syscall/js.valueInvoke": () => 0,
        "syscall/js.valueNew": () => 0,
        "syscall/js.valueLength": () => 0,
        "syscall/js.valuePrepareString": () => 0,
        "syscall/js.valueLoadString": () => {},
        "syscall/js.valueInstanceOf": () => false,
        "syscall/js.copyBytesToGo": () => 0,
        "syscall/js.copyBytesToJS": () => 0,
      },
      // JavaScript 绑定
      js: {
        // 内存导入
        mem: new WebAssembly.Memory({ initial: 256, maximum: 256 }),
        // 全局变量设置函数
        setGlobal: (key, value) => {
          const keyStr = this.getStringFromMemory(key);
          window[keyStr] = value;
        },
        // 字符串获取函数
        getString: (ptr, len) => {
          const memory = new Uint8Array(this.wasmMemory.buffer);
          const bytes = memory.slice(ptr, ptr + len);
          return textDecoder.decode(bytes);
        },
      },
    };
  }

  // 从 WASM 内存中获取字符串
  getStringFromMemory(ptr) {
    if (!this.wasmMemory) return "";

    const memory = new Uint8Array(this.wasmMemory.buffer);
    let len = 0;
    while (memory[ptr + len] !== 0) len++;

    const bytes = memory.slice(ptr, ptr + len);
    return new TextDecoder().decode(bytes);
  }

  // 直接从 WASM 导出设置函数
  setupWASMFunctions() {
    if (!this.wasmInstance || !this.wasmInstance.exports) {
      throw new Error("WASM 实例未正确初始化");
    }

    const exports = this.wasmInstance.exports;

    // 检查并设置混淆函数
    if (exports.obfuscateJS) {
      window.obfuscateJS = exports.obfuscateJS;
      console.log("obfuscateJS 函数已设置");
    } else {
      console.warn("未找到 obfuscateJS 导出函数");
    }

    // 检查并设置测试函数
    if (exports.wasmTest) {
      window.wasmTest = exports.wasmTest;
      console.log("wasmTest 函数已设置");
    }

    console.log("可用的 WASM 导出:", Object.keys(exports));
  }

  async waitForWASMFunctions() {
    const maxAttempts = 50;
    let attempts = 0;

    while (attempts < maxAttempts) {
      console.log(`等待 WASM 函数注册... 尝试 ${attempts + 1}/${maxAttempts}`);

      if (window.wasmReady && window.wasmTest && window.obfuscateJS) {
        console.log("所有 WASM 函数注册成功！");
        return;
      }

      await new Promise((resolve) => setTimeout(resolve, 200));
      attempts++;
    }

    throw new Error("WASM 函数注册超时");
  }

  async loadGoRuntime() {
    return new Promise((resolve, reject) => {
      if (window.Go) {
        resolve();
        return;
      }

      const script = document.createElement("script");
      script.src = "/wasm_exec.js";
      script.onload = () => {
        if (window.Go) {
          resolve();
        } else {
          reject(new Error("Go runtime not available"));
        }
      };
      script.onerror = () => reject(new Error("Failed to load Go runtime"));
      document.head.appendChild(script);
    });
  }

  // 更新 WASM 状态显示
  updateWASMStatus() {
    const statusElement = document.getElementById("wasmStatus");
    if (statusElement) {
      if (this.isReady) {
        statusElement.textContent = "WASM 状态: 已就绪";
        statusElement.style.color = "#28a745";
      } else {
        statusElement.textContent = "WASM 状态: 加载失败";
        statusElement.style.color = "#dc3545";
      }
    }
  }

  // 获取混淆配置
  getObfuscatorConfig() {
    return {
      identifierObfuscation: document.getElementById("identifierObfuscation")
        .checked,
      stringEncryption: document.getElementById("stringEncryption").checked,
      controlFlowFlattening: document.getElementById("controlFlowFlattening")
        .checked,
      // 死代码注入和表达式分解功能已移除
      compactCode: document.getElementById("compactCode").checked,
      preserveComments: false,
    };
  }

  // 执行代码混淆
  async obfuscateCode() {
    if (!this.isReady) {
      this.showStatus("WASM 模块未就绪，请稍后再试", "error");
      return;
    }

    const inputCode = document.getElementById("inputCode").value.trim();
    if (!inputCode) {
      this.showStatus("请输入要混淆的代码", "warning");
      return;
    }

    // 显示加载状态
    this.setLoadingState(true);

    try {
      const config = this.getObfuscatorConfig();
      console.log("开始混淆，配置:", config);

      const result = window.obfuscateJS(inputCode, JSON.stringify(config));
      console.log("混淆结果:", result);

      if (result.success) {
        document.getElementById("outputCode").value = result.code;
        this.updateOutputStats(result.stats);
        this.showStatus("代码混淆完成！", "success");
      } else {
        this.showStatus("混淆失败: " + result.error, "error");
      }
    } catch (error) {
      console.error("混淆过程出错:", error);
      this.showStatus("混淆过程出错: " + error.message, "error");
    } finally {
      this.setLoadingState(false);
    }
  }

  // 验证功能已移除

  // 复制结果到剪贴板
  async copyResult() {
    const outputCode = document.getElementById("outputCode").value;
    if (!outputCode) {
      this.showStatus("没有可复制的内容", "warning");
      return;
    }

    try {
      await navigator.clipboard.writeText(outputCode);
      this.showCopySuccess();
    } catch (error) {
      // 降级方案
      const textarea = document.getElementById("outputCode");
      textarea.select();
      document.execCommand("copy");
      this.showCopySuccess();
    }
  }

  // 下载结果文件
  downloadResult() {
    const outputCode = document.getElementById("outputCode").value;
    if (!outputCode) {
      this.showStatus("没有可下载的内容", "warning");
      return;
    }

    const blob = new Blob([outputCode], { type: "application/javascript" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "obfuscated.js";
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);

    this.showStatus("文件下载已开始", "success");
  }

  // 加载示例代码
  loadExampleCode() {
    const exampleCode = `// 示例 JavaScript 代码
function calculateSum(numbers) {
    let sum = 0;
    for (let i = 0; i < numbers.length; i++) {
        sum += numbers[i];
    }
    return sum;
}

function greetUser(name) {
    const greeting = "Hello, " + name + "!";
    console.log(greeting);
    return greeting;
}

// 主函数
function main() {
    const testNumbers = [1, 2, 3, 4, 5];
    const result = calculateSum(testNumbers);
    
    greetUser("World");
    
    if (result > 10) {
        console.log("Sum is greater than 10: " + result);
    } else {
        console.log("Sum is less than or equal to 10: " + result);
    }
}

// 执行主函数
main();`;

    document.getElementById("inputCode").value = exampleCode;
    this.updateInputStats();
  }

  // 清空输入
  clearInput() {
    document.getElementById("inputCode").value = "";
    document.getElementById("outputCode").value = "";
    this.updateInputStats();
    this.updateOutputStats(null);
    this.hideStatus();
  }

  // 更新输入统计信息
  updateInputStats() {
    const inputCode = document.getElementById("inputCode").value;
    const charCount = inputCode.length;
    const lineCount = inputCode.split("\n").length;

    document.getElementById(
      "inputStats"
    ).textContent = `字符数: ${charCount} | 行数: ${lineCount}`;
  }

  // 更新输出统计信息
  updateOutputStats(stats) {
    const outputCode = document.getElementById("outputCode").value;
    const charCount = outputCode.length;

    let compressionText = "";
    if (stats && stats.originalSize > 0) {
      const compressionRatio = ((1 - stats.compression) * 100).toFixed(1);
      compressionText = ` | 压缩率: ${compressionRatio}%`;
    }

    document.getElementById(
      "outputStats"
    ).textContent = `字符数: ${charCount}${compressionText}`;
  }

  // 显示状态信息
  showStatus(message, type = "success") {
    const statusPanel = document.getElementById("statusPanel");
    const statusMessage = statusPanel.querySelector(".status-message");

    statusPanel.className = `status-panel ${type}`;
    statusMessage.textContent = message;
    statusPanel.style.display = "block";

    // 3秒后自动隐藏
    setTimeout(() => {
      this.hideStatus();
    }, 3000);
  }

  // 隐藏状态信息
  hideStatus() {
    const statusPanel = document.getElementById("statusPanel");
    if (statusPanel) {
      statusPanel.style.display = "none";
    }
  }

  // 错误面板功能已移除

  // 显示复制成功提示
  showCopySuccess() {
    const toast = document.createElement("div");
    toast.className = "copy-success";
    toast.textContent = "已复制到剪贴板";
    toast.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      background: #28a745;
      color: white;
      padding: 10px 20px;
      border-radius: 4px;
      z-index: 1000;
    `;
    document.body.appendChild(toast);

    setTimeout(() => {
      document.body.removeChild(toast);
    }, 2000);
  }

  // 设置加载状态
  setLoadingState(loading) {
    const btn = document.getElementById("obfuscateBtn");
    if (btn) {
      const btnText = btn.querySelector(".btn-text");
      const btnLoading = btn.querySelector(".btn-loading");

      if (loading) {
        btn.disabled = true;
        if (btnText) btnText.style.display = "none";
        if (btnLoading) btnLoading.style.display = "inline-flex";
      } else {
        btn.disabled = false;
        if (btnText) btnText.style.display = "inline";
        if (btnLoading) btnLoading.style.display = "none";
      }
    }
  }

  // 保存配置到本地存储
  saveConfig() {
    const config = this.getObfuscatorConfig();
    localStorage.setItem("obfuscator-config", JSON.stringify(config));
  }

  // 从本地存储加载配置
  loadConfig() {
    try {
      const savedConfig = localStorage.getItem("obfuscator-config");
      if (savedConfig) {
        const config = JSON.parse(savedConfig);

        Object.keys(config).forEach((key) => {
          const element = document.getElementById(key);
          if (element && typeof config[key] === "boolean") {
            element.checked = config[key];
          }
        });
      }
    } catch (error) {
      console.warn("加载配置失败:", error);
    }
  }

  showError(message) {
    console.error(message);
    this.showStatus(message, "error");
  }
}

// 应用启动
document.addEventListener("DOMContentLoaded", () => {
  console.log("DOM 加载完成，启动应用...");
  new JSObfuscatorApp();
});
