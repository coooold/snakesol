import { GRID_SIZE } from './constants.js';

export class Renderer {
    constructor(canvas) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d');
        this.gridSize = GRID_SIZE;
        this.cols = 0;
        this.rows = 0;
        this.scale = 1;
        this.appleAlpha = 1;
        this.lastTime = performance.now();
        this.setupCanvas();
        window.addEventListener('resize', () => this.handleResize());
        this.startAppleAnimation();
    }

    setupCanvas() {
        this.canvas.style.backgroundColor = '#f0f0f0';
        this.ctx.strokeStyle = '#ccc';
        if (this.cols && this.rows) {
            this.updateCanvasSize();
        }
    }

    startAppleAnimation() {
        const animate = (currentTime) => {
            const deltaTime = currentTime - this.lastTime;
            this.lastTime = currentTime;
            
            // 更新苹果的透明度，创建闪烁效果
            this.appleAlpha = 0.5 + Math.sin(currentTime * 0.005) * 0.5;
            
            requestAnimationFrame(animate);
        };
        requestAnimationFrame(animate);
    }

    clear() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    }

    drawSnake(snake) {
        this.ctx.fillStyle = snake.color;
        
        // 绘制蛇身
        snake.body.forEach(segment => {
            this.ctx.fillRect(
                segment.x * this.gridSize,
                segment.y * this.gridSize,
                this.gridSize,
                this.gridSize
            );
        });

        // 绘制蛇头（稍大一些以区分）
        this.ctx.fillRect(
            snake.x * this.gridSize,
            snake.y * this.gridSize,
            this.gridSize + 1,
            this.gridSize + 1
        );

        // 如果是玩家控制的蛇，添加名字标签
        if (!snake.isAI) {
            this.ctx.font = 'bold 12px Arial';
            this.ctx.fillStyle = '#FF0000';
            this.ctx.fillText(
                snake.name,
                snake.x * this.gridSize,
                snake.y * this.gridSize - 2
            );
        }
    }

    drawApple(apple) {
        // 使用当前的透明度值
        this.ctx.fillStyle = `rgba(255, 0, 0, ${this.appleAlpha})`;
        this.ctx.beginPath();
        this.ctx.arc(
            apple.x * this.gridSize + this.gridSize/2,
            apple.y * this.gridSize + this.gridSize/2,
            this.gridSize/2,
            0,
            Math.PI * 2
        );
        this.ctx.fill();
    }

    handleResize() {
        if (this.cols && this.rows) {
            this.updateCanvasSize();
        }
    }

    updateCanvasSize() {
        // 获取窗口可见区域的尺寸，减去一些边距
        const windowWidth = window.innerWidth * 0.95;
        const windowHeight = window.innerHeight * 0.95;
        const gameAspectRatio = this.cols / this.rows;
        const windowAspectRatio = windowWidth / windowHeight;

        if (windowAspectRatio > gameAspectRatio) {
            // 以高度为基准计算缩放比例
            this.scale = windowHeight / (this.rows * this.gridSize);
        } else {
            // 以宽度为基准计算缩放比例
            this.scale = windowWidth / (this.cols * this.gridSize);
        }

        // 应用缩放后的尺寸
        this.canvas.width = this.cols * this.gridSize * this.scale;
        this.canvas.height = this.rows * this.gridSize * this.scale;

        // 设置canvas样式以确保完整显示
        this.canvas.style.width = `${this.canvas.width}px`;
        this.canvas.style.height = `${this.canvas.height}px`;
        this.canvas.style.display = 'block';
        this.canvas.style.margin = 'auto';
    }

    setGameSize(cols, rows) {
        this.cols = cols;
        this.rows = rows;
        this.updateCanvasSize();
    }

    draw(player, snakes, apples) {
        this.clear();
        this.ctx.save();
        this.ctx.scale(this.scale, this.scale);

        // 绘制所有苹果
        apples.forEach(apple => this.drawApple(apple));

        // 绘制所有活着的蛇
        snakes.forEach(snake => {
            if (!snake.dead) {
                this.drawSnake(snake);
            }
        });

        this.ctx.restore();
    }
}