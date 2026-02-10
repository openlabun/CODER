import sys

def solve():
    # Leer la entrada desde stdin
    s = sys.stdin.read().strip()
    if not s:
        return

    stack = []
    current_num = 0
    last_op = '+'

    for i, char in enumerate(s):
        if char.isdigit():
            current_num = int(char)
        
        # Si es un operador o es el último carácter
        if char in '+-*' or i == len(s) - 1:
            if last_op == '+':
                stack.append(current_num)
            elif last_op == '-':
                stack.append(-current_num)
            elif last_op == '*':
                # Precedencia: Multiplicar con el último valor en la pila
                top = stack.pop()
                stack.append(top * current_num)
            
            last_op = char
            current_num = 0

    print(sum(stack))

if __name__ == "__main__":
    solve()
