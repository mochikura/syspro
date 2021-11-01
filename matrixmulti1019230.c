#include <stdio.h>
#include <unistd.h>
#include <pthread.h>
#include <stdlib.h>
#define h 6  //horizontal
#define v 6  //vertical
struct index //並列実行する関数へ渡す要素の構造体
{
    int x;
    int y;
};
int mx[v][h] = {{1, 2, 3, 4, 5, 6}, {2, 3, 4, 5, 6, 7}, {3, 4, 5, 6, 7, 8}, {4, 5, 6, 7, 8, 9}, {5, 6, 7, 8, 9, 0}, {6, 7, 8, 9, 0, 1}};
int my[v][h] = {{9, 8, 7, 6, 5, 4}, {8, 7, 6, 5, 4, 3}, {7, 6, 5, 4, 3, 2}, {6, 5, 4, 3, 2, 1}, {5, 4, 3, 2, 1, 0}, {4, 3, 2, 1, 0, 9}};
//計算対象mx*my
void *multi(void *arg);
//プロトタイプ
int main()
{
    struct index p[v][h];//配列の番地を渡すために配列での構造体
    pthread_t th[v][h]; //スレッドを入れる
    int x, y;           //for用
    void *ans[v][h];    //答えを入れる、ポインタでの受け渡しがやりやすかったので最初からvoidにした
    for (x = 0; x < v; x++)
    {
        for (y = 0; y < h; y++)
        {
            p[x][y].x = x;
            p[x][y].y = y;
            int r;
            r = pthread_create(&th[x][y], NULL, multi, (void *)&p[x][y]); //スレッド作成、構造体を渡す
            if (r != 0)                                             //エラー
                perror("new thread");
        }
    }
    for (x = 0; x < v; x++)
    {
        for (y = 0; y < h; y++)
        {
            //printf(":%d%d\n", x, y);
            pthread_join(th[x][y], &ans[x][y]); //スレッドからansに値を受け取る
        }
    }
    //printf("check\n");
    for (y = 0; y < h; y++)
    {
        for (x = 0; x < v; x++)
        {
            printf("%d ", *(int *)ans[y][x]);
            //ポインタをintのポインタにしてポインタにする？？？？？？
            free(ans[x][y]);
            //メモリ解放
        }
        printf("\n");
    }
    printf("done\n");
}
//いろいろ調べながらやったので中身をすべて説明できる自信はないですが...

void *multi(void *arg)
{
    struct index *data = (struct index *)arg; //構造体を設定
    void *rt = malloc(sizeof(int));           //ポインタで受け渡しするのでメモリ確保でやる
    *(int *)(rt) = 0;                         //初期化
    sleep(1);                                 //？
    int n;
    for (n = 0; n < v; n++)
    {
        *(int *)(rt) += mx[data->x][n] * my[n][data->y];
        //構造体の番号を持ってきて配列を乗算して加算する
        //printf("%d*%d=%d\n", mx[data->x][n], my[n][data->y], *(int *)rt);
    }
    //printf("p%d %d %d\n", data->x, data->y, *(int *)rt);
    pthread_exit(rt);
    //ポインタを渡す
}