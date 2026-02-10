const readline = require('readline');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

rl.on('line', (line) => {
    const s = line.trim();
    if (!s) return;
    
    const stack = [];
    let currentNum = 0;
    let lastOp = '+';
    
    for (let i = 0; i < s.length; i++) {
        const char = s[i];
        
        if (/\d/.test(char)) {
            currentNum = parseInt(char);
        }
        
        if (['+', '-', '*'].includes(char) || i === s.length - 1) {
            if (lastOp === '+') {
                stack.push(currentNum);
            } else if (lastOp === '-') {
                stack.push(-currentNum);
            } else if (lastOp === '*') {
                stack.push(stack.pop() * currentNum);
            }
            lastOp = char;
            currentNum = 0;
        }
    }
    
    // Sumar todo el stack
    const result = stack.reduce((a, b) => a + b, 0);
    console.log(result);
});
