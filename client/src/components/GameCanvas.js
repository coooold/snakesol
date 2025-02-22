// 游戏画布组件，负责渲染游戏界面
export class GameCanvas {
    constructor(canvasId) {
        this.canvas = document.getElementById(canvasId);
        this.ctx = this.canvas.getContext('2d');
        this.gridSize = 10;
        this.setupCanvas();
    }

    setupCanvas() {
        // 设置画布样式
        this.canvas.style.backgroundColor = '#f0f0f0';
        this.ctx.strokeStyle = '#ccc';
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
            snake.x * this.gridSize - 1,
            snake.y * this.gridSize - 1,
            this.gridSize + 2,
            this.gridSize + 2
        );

        // 如果是玩家控制的蛇，添加名字标签
        if (!snake.isAI) {
            this.ctx.fillStyle = '#000';
            this.ctx.font = '12px Arial';
            this.ctx.fillText(
                snake.name,
                snake.x * this.gridSize,
                snake.y * this.gridSize - 5
            );
        }
    }

    drawApple(apple) {
        this.ctx.fillStyle = '#ff0000';
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

    render(gameState) {
        this.clear();
        
        // 绘制所有苹果
        gameState.apples.forEach(apple => this.drawApple(apple));
        
        // 绘制所有蛇
        gameState.snakes.forEach(snake => this.drawSnake(snake));
    }
}