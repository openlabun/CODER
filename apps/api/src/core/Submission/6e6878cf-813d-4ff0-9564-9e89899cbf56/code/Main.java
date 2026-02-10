import java.util.Scanner;
import java.util.Stack;

public class Main {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        if (scanner.hasNext()) {
            String s = scanner.next();
            Stack<Integer> stack = new Stack<>();
            int currentNum = 0;
            char lastOp = '+';

            for (int i = 0; i < s.length(); i++) {
                char c = s.charAt(i);

                if (Character.isDigit(c)) {
                    currentNum = c - '0';
                }

                if (c == '+'  c == '-'  c == '' || i == s.length() - 1) {
                    if (lastOp == '+') {
                        stack.push(currentNum);
                    } else if (lastOp == '-') {
                        stack.push(-currentNum);
                    } else if (lastOp == '') {
                        stack.push(stack.pop() * currentNum);
                    }
                    lastOp = c;
                    currentNum = 0;
                }
            }

            int result = 0;
            for (int num : stack) {
                result += num;
            }
            System.out.println(result);
        }
