class PasswordGenerator {
    constructor() {
        this.checkboxes = ['uppercase', 'lowercase', 'numbers', 'symbols'];
        this.initializeElements();
        this.attachEventListeners();
        this.generatePassword();
        // CSRFトークンを取得（テンプレートから埋め込まれたトークン）
        this.csrfToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content') || '';
    }

    initializeElements() {
        this.elements = {
            symbolsCheckbox: document.getElementById('symbols'),
            symbolsArea: document.getElementById('symbolsCustomArea'),
            generateButton: document.getElementById('generateButton'),
            copyButton: document.getElementById('copyButton'),
            passwordField: document.getElementById('password'),
            customSymbols: document.getElementById('customSymbols')
        };
    }

    attachEventListeners() {
        // 記号カスタム入力欄の表示制御
        this.elements.symbolsCheckbox.addEventListener('change', () => {
            this.elements.symbolsArea.style.display = 
                this.elements.symbolsCheckbox.checked ? 'block' : 'none';
        });

        // チェックボックスのイベントリスナー
        this.checkboxes.forEach(id => {
            document.getElementById(id).addEventListener('change', () => {
                this.validateOptions();
                this.generatePassword();
            });
        });

        // パスワード生成ボタン
        this.elements.generateButton.addEventListener('click', () => {
            this.generatePassword();
        });

        // コピーボタン
        this.elements.copyButton.addEventListener('click', () => {
            this.copyPassword();
        });

        // 文字数選択のイベントリスナー
        document.querySelectorAll('input[name="length"]').forEach(radio => {
            radio.addEventListener('change', () => this.generatePassword());
        });

        // カスタム記号入力のイベントリスナー
        this.elements.customSymbols.addEventListener('change', () => {
            this.generatePassword();
        });

        // ページ読み込み時の表示制御
        document.addEventListener('DOMContentLoaded', () => {
            this.elements.symbolsArea.style.display = 
                this.elements.symbolsCheckbox.checked ? 'block' : 'none';
        });
    }

    validateOptions() {
        const anyChecked = this.checkboxes.some(id => document.getElementById(id).checked);
        this.elements.generateButton.disabled = !anyChecked;
        return anyChecked;
    }

    async generatePassword() {
        if (!this.validateOptions()) return;

        try {
            const params = new URLSearchParams();
            params.append('length', document.querySelector('input[name="length"]:checked').value);
            params.append('uppercase', document.getElementById('uppercase').checked.toString());
            params.append('lowercase', document.getElementById('lowercase').checked.toString());
            params.append('numbers', document.getElementById('numbers').checked.toString());
            params.append('symbols', document.getElementById('symbols').checked.toString());
            params.append('customSymbols', this.elements.customSymbols.value);

            const response = await fetch('/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                    'X-CSRF-Token': this.csrfToken // CSRFトークンを送信
                },
                body: params
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const password = await response.text();
            this.elements.passwordField.value = password;
            this.elements.copyButton.disabled = !password;
        } catch (error) {
            console.error('Error:', error);
            this.elements.passwordField.value = 'エラーが発生しました';
        }
    }

    copyPassword() {
        this.elements.passwordField.select();
        document.execCommand('copy');
        
        const originalText = this.elements.copyButton.textContent;
        this.elements.copyButton.textContent = 'コピーしました';
        setTimeout(() => {
            this.elements.copyButton.textContent = originalText;
        }, 2000);
    }
}

// インスタンス化
new PasswordGenerator();