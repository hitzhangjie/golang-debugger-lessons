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
        ptrace(PTRACE_TRACEME, 0, NULL, NULL);
        execl("/bin/ls", "~", NULL);
    }
    else {
        wait(NULL);
        orig_eax = ptrace(PTRACE_PEEKUSER, child, (void *)(ORIG_EAX*8), (void *)NULL);
        printf("The child made a system call %ld\n", orig_eax);
        ptrace(PTRACE_CONT, child, NULL, NULL);
    }
    return 0;
}
