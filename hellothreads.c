#include <stdio.h>
#include <unistd.h>
#include <pthread.h>
#define N 16
static int n[N];
void *func(void *arg);
int main(int argc, char *argv[])
{
    pthread_t th[N];
    int i;
    for (i = 0; i < N; i++)
    {
        int r;
        n[i] = i + 1;
        r = pthread_create(&th[i], NULL, func, (void *)&n[i]);
        if (r != 0)
            perror("new thread");
    }
    for (i = 0; i < N; i++)
        pthread_join(th[i], NULL);
    printf("done.\n");
    return 0;
}
void *func(void *arg)
{
    int n = *(int *)arg;
    sleep(1);
printf("Hello, I am thread-%02d.\n", n);
pthread_exit(NULL);
}