#include <stdio.h>
#include <pthread.h>
#include <unistd.h>
#include <stdlib.h>
static long num_rects = 2 * 1000 * 1000 * 1000;
#define N_THREADS (4)//スレッド数
static double ans[N_THREADS];//各区分の計算結果を入れる
//並列化した場所について
//積分区間をN_THREADSの数だけ分割して、それぞれの区分において並列処理で計算
//各区分で計算した値を全て加算することで円周率を求めた
void *comp(void *arg);
int main(int argc, char *argv[])
{
    int i;//for用
    pthread_t th[N_THREADS];//スレッド格納
    double area;//最終的な回答
    int num[N_THREADS];//スレッドに渡す区分番号
    for (i = 0; i < N_THREADS; i++)
    {
        num[i] = i;//区分番号を割り当て
        int r;
        r = pthread_create(&th[i], NULL, comp, (void *)&num[i]);
        //スレッド生成、区分番号をcompに渡す
        if (r != 0)
            perror("new thread");
    }
    for (i = 0; i < N_THREADS; i++)
    {
        pthread_join(th[i], NULL);
        //スレッド
    }
    double width;
    width = 1.0 / (double)num_rects;//area計算で必要
    for (i = 0; i < N_THREADS; i++)
    {
        area += width * ans[i];//最終的な円周率の結果を計算
        //区分計算しているので、compでは高さを全て計算しきっている。
        //ここでまとめて幅を乗算することで、区分内の面積を出している
        //printf("%f\n", area);
    }
    printf("Computed pi = %f\n", area);//出力
    return 0;
}

void *comp(void *arg)
{
    int n = *(int *)arg;//受け取る引数
    int range = num_rects / N_THREADS;//計算上の細かい区分
    int i;
    double mid, width, height;
    ans[n] = 0.0;//ここの区分の初期化
    width = 1.0 / (double)num_rects;//計算上の区分
    sleep(1);
    for (i = n * range; i < (n + 1) * range; ++i)
    {
        mid = (i + 0.5) * width;//区分内の中央値
        height = 4.0 / (1.0 + mid * mid);//上の区分内の中央値の座標を出す
        ans[n] += height;//座標値を加算
    }
    //printf("%d : %f\n", n, ans[n]);
    pthread_exit(NULL);//おわり
}