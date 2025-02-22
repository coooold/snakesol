import { Renderer } from './Renderer.js';
import { Controller } from './Controller.js';

export class Game {
    constructor() {
        this.canvas = document.getElementById('gameCanvas');
        this.renderer = new Renderer(this.canvas);
        this.player = null;
        this.ws = null;
        this.snakes = new Map();
        this.apples = [];
        this.controller = null;

        // 初始化canvas尺寸
        this.initCanvasSize();
        // 监听窗口大小变化
        window.addEventListener('resize', () => this.initCanvasSize());

        this.init();
    }

    initCanvasSize() {
        // 获取窗口可见区域的尺寸，减去一些边距
        const windowWidth = window.innerWidth * 0.95;
        const windowHeight = window.innerHeight * 0.95;

        // 设置canvas尺寸
        this.canvas.width = windowWidth;
        this.canvas.height = windowHeight;

        // 设置canvas样式以确保完整显示
        this.canvas.style.width = `${windowWidth}px`;
        this.canvas.style.height = `${windowHeight}px`;
        this.canvas.style.display = 'block';
        this.canvas.style.margin = 'auto';
    }

    init() {
        this.connectWebSocket();
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
        const hostname = window.location.hostname;
        const port = window.location.port ? `:${window.location.port}` : ''; // 如果端口不是默认的 80 或 443，则需要包含
         
        // 构建 WebSocket URL
        const socketUrl = `${protocol}${hostname}${port}/ws`; 

        this.ws = new WebSocket(socketUrl);

        this.ws.onopen = () => {
            console.log('已连接到服务器');
        };

        this.ws.onmessage = (event) => {
            const gameState = JSON.parse(event.data);
            this.updateGameState(gameState);
        };

        this.ws.onclose = () => {
            console.log('与服务器断开连接');
            setTimeout(() => this.init(), 1000);
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket错误:', error);
        };
    }

    updateGameState(state) {
        console.log(state)
        // 如果收到场景尺寸信息，更新渲染器配置
        if (state.config && state.config.cols && state.config.rows) {
            this.renderer.setGameSize(state.config.cols, state.config.rows);
        }

        // 更新蛇的状态
        this.snakes.clear();
        for (const [id, snakeData] of Object.entries(state.snakes)) {
            if (!snakeData.isAI && !this.player) {
                this.player = snakeData;
                // 创建控制器并设置方向改变回调
                if (!this.controller) {
                    this.controller = new Controller(this.player);
                    this.controller.onDirectionChange = (direction) => {
                        this.sendDirection(direction);
                    };
                }
            }
            this.snakes.set(id, snakeData);
        }

        // 更新苹果位置
        this.apples = state.apples;

        // 渲染游戏画面
        this.renderer.draw(this.player, Array.from(this.snakes.values()), this.apples);

        // 检查玩家是否死亡
        if (state.deadSnakeId && this.player && state.deadSnakeId === this.player.id && !document.getElementById('restartButton')) {
            const score = this.player.body.length;
            // 创建重新开始按钮
            const restartButton = document.createElement('button');
            restartButton.id = 'restartButton';
            restartButton.textContent = '重新开始游戏';
            restartButton.style.position = 'fixed';
            restartButton.style.top = '20px';
            restartButton.style.left = '50%';
            restartButton.style.transform = 'translateX(-50%)';
            restartButton.style.padding = '10px 20px';
            restartButton.style.fontSize = '16px';
            restartButton.style.backgroundColor = '#4CAF50';
            restartButton.style.color = 'white';
            restartButton.style.border = 'none';
            restartButton.style.borderRadius = '5px';
            restartButton.style.cursor = 'pointer';
            restartButton.style.zIndex = '1000';
            
            // 添加得分显示
            const scoreDisplay = document.createElement('div');
            scoreDisplay.id = 'scoreDisplay';
            scoreDisplay.textContent = `游戏结束！得分：${score}`;
            scoreDisplay.style.position = 'fixed';
            scoreDisplay.style.top = '70px';
            scoreDisplay.style.left = '50%';
            scoreDisplay.style.transform = 'translateX(-50%)';
            scoreDisplay.style.fontSize = '18px';
            scoreDisplay.style.color = '#333';
            scoreDisplay.style.zIndex = '1000';
            
            document.body.appendChild(scoreDisplay);
            document.body.appendChild(restartButton);
            
            restartButton.onclick = () => {
                this.player = null;
                this.ws.close();
                document.body.removeChild(restartButton);
                document.body.removeChild(scoreDisplay);
                this.init();
            };
        }
    }

    sendDirection(direction) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({
                type: 'direction',
                payload: direction
            }));
        }
    }

    cleanup() {
        if (this.ws) {
            this.ws.close();
        }
        if (this.controller) {
            this.controller.cleanup();
        }
    }
}