# 前端静态资源目录说明

## 目录结构

```
client/
├── index.html      # 游戏主页面
└── src/
    ├── js/         # JavaScript 源代码
    │   ├── Controller.js  # 游戏控制器
    │   ├── Game.js        # 游戏核心逻辑
    │   ├── GameStats.js   # 游戏统计
    │   ├── Renderer.js    # 游戏渲染器
    │   ├── Snake.js       # 蛇实体类
    │   └── constants.js   # 常量定义
    ├── components/  # 可复用的UI组件
    └── network/     # 网络通信相关代码
```

## 说明

- `index.html`: 游戏的主入口页面
- `src/js/`: 包含所有游戏相关的JavaScript源代码
  - `Controller.js`: 处理用户输入和游戏控制
  - `Game.js`: 实现游戏的核心逻辑
  - `GameStats.js`: 处理游戏统计数据
  - `Renderer.js`: 负责游戏画面的渲染
  - `Snake.js`: 定义蛇的行为和属性
  - `constants.js`: 存放游戏中使用的常量
- `src/components/`: 存放可复用的UI组件
- `src/network/`: 包含与服务器通信相关的代码