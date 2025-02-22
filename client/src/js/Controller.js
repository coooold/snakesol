import { DIRECTIONS } from './constants.js';

export class Controller {
    constructor(player) {
        this.player = player;
        this.onDirectionChange = null;
        this.setupKeyboardControls();
    }

    setupKeyboardControls() {
        document.addEventListener('keydown', this.handleKeyPress.bind(this));
    }

    handleKeyPress(event) {
        const key = event.key;
        let newDirection;

        switch (key) {
            case 'ArrowUp':
                newDirection = DIRECTIONS.UP;
                break;
            case 'ArrowDown':
                newDirection = DIRECTIONS.DOWN;
                break;
            case 'ArrowLeft':
                newDirection = DIRECTIONS.LEFT;
                break;
            case 'ArrowRight':
                newDirection = DIRECTIONS.RIGHT;
                break;
            default:
                return;
        }

        // 防止180度转向
        if (this.player.direction.x + newDirection.x === 0 &&
            this.player.direction.y + newDirection.y === 0) {
            return;
        }

        // 更新方向并通知游戏
        this.player.direction = newDirection;
        if (this.onDirectionChange) {
            this.onDirectionChange(newDirection);
        }
    }

    cleanup() {
        document.removeEventListener('keydown', this.handleKeyPress.bind(this));
    }
}