import re

def analyze_log():
    # 初始化计数器
    total_raydium_calls = 0
    stack_height_counts = {}
    raydium_count_groups = {}  # 新增：按 raydiumCallCount 值分组
    
    # 读取文件并分析
    try:
        with open('./log', 'r') as file:
            for line in file:
                # 查找 raydiumCallCount
                raydium_match = re.search(r'raydiumCallCount:(\d+)', line)
                if raydium_match:
                    count = int(raydium_match.group(1))
                    total_raydium_calls += count
                    
                    # 统计相同 raydiumCallCount 值的出现次数
                    raydium_count_groups[count] = raydium_count_groups.get(count, 0) + 1
                    
                    # 查找对应的 maxStackHeight
                    stack_match = re.search(r'maxStackHeight: (\d+)', line)
                    if stack_match:
                        stack_height = int(stack_match.group(1))
                        stack_height_counts[stack_height] = stack_height_counts.get(stack_height, 0) + count

        # 打印结果
        print(f"总 raydiumCallCount: {total_raydium_calls}")
        
        print("\n各 maxStackHeight 的 raydiumCallCount 统计:")
        for height, count in sorted(stack_height_counts.items()):
            print(f"maxStackHeight {height}: {count}")
            
        print("\n各 raydiumCallCount 值的出现次数:")
        for count_value, frequency in sorted(raydium_count_groups.items()):
            print(f"raydiumCallCount={count_value}: 出现 {frequency} 次")

    except FileNotFoundError:
        print("找不到日志文件，请确保文件存在并且命名为 'log.txt'")
    except Exception as e:
        print(f"处理文件时发生错误: {str(e)}")

if __name__ == "__main__":
    analyze_log()
