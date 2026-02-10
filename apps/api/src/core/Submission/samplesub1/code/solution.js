// Lee entrada estándar y suma los números
process.stdin.on("data", (data) => {
    const nums = data.toString().trim().split(/\s+/).map(Number);
    const sum = nums.reduce((a, b) => a + b, 0);
    console.log(sum);
});
