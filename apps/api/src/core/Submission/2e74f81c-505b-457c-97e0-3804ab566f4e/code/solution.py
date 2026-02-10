import sys
def main():
    data = sys.stdin.read().strip().split()
    if not data:
        return
    nums = list(map(int, data))
    print(sum(nums))
if __name__ == "__main__":
    main()
