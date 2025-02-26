#include <sys/ptrace.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>
// see `struct user_regs_struct`, orig_eax field number is 11,
// 11 * 8 is the register offset in memory on arch x86_64.
//
// see: /usr/include/sys/user.h 
#define ORIG_EAX 11

int main()
{   pid_t child;
    long long orig_eax;
    child = fork();
    if(child == 0) {
        // 子进程执行到这里后会请求系统调用，通知tracer跟踪自己
        ptrace(PTRACE_TRACEME, 0, NULL, NULL);
        // 等tracer调用ptrace(PTRACE_CONT,...)时才继续恢复下面的执行
        execl("/bin/ls", "~", NULL);
    }
    else {
        // 等待子进程变成tracee traced状态
        wait(NULL);
        // 读取tracee的数据，因为子进程ptrace执行成功了，所以此时返回值应该是0
        orig_eax = ptrace(PTRACE_PEEKUSER, child, (void *)(ORIG_EAX*8), (void *)NULL);
        printf("The child made a system call %ld\n", orig_eax);
        // 睡眠1s方便观察ptrace(PTRACE_CONT,...)前后的效果
        sleep(1);
        // 通知tracee恢复执行，此时tracee会把ls ~的输出显示出来
        ptrace(PTRACE_CONT, child, NULL, NULL);
    }
    return 0;
}
