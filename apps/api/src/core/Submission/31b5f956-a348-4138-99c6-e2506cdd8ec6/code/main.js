const fs = require('fs');

// Leer toda la entrada
const input = fs.readFileSync(0, 'utf-8').trim();
const numbers = input.split(' ');

const a = parseInt(numbers[0]);
const b = parseInt(numbers[1]);

console.log(a + b);