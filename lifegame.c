/*
 * The Life Game
 */
#include <stdio.h>
#include <unistd.h>
#include <pthread.h>
#include <stdlib.h>

#define NLOOP (200)
#define N (8192)
#define M (8192)
#define ALIVE (1)
#define DEAD (0)

#define N_THREADS (8192)

typedef int Grid[N + 2][M + 2];
static Grid g[2];
static int cur;


void *computeNextGen(void *arg)
{
	int j;
	int num = *(int *)arg;

	for (j = 1; j <= M; ++j)
	{
		int count = 0;

		/* NW neighbor */
		if (g[cur][num - 1][j - 1] == ALIVE)
			count++;
		/* N neighbor */
		if (g[cur][num - 1][j] == ALIVE)
			count++;
		/* NE neighbor */
		if (g[cur][num - 1][j + 1] == ALIVE)
			count++;
		/* W neighbor */
		if (g[cur][num][j - 1] == ALIVE)
			count++;
		/* E neighbor */
		if (g[cur][num][j + 1] == ALIVE)
			count++;
		/* SW neighbor */
		if (g[cur][num + 1][j - 1] == ALIVE)
			count++;
		/* S neighbor */
		if (g[cur][num + 1][j] == ALIVE)
			count++;
		/* SE neighbor */
		if (g[cur][num + 1][j + 1] == ALIVE)
			count++;

		if (count <= 1 || count >= 4)
			g[(cur + 1) & 1][num][j] = DEAD;
		else if (g[cur][num][j] == ALIVE &&
				 (count == 2 || count == 3))
			g[(cur + 1) & 1][num][j] = ALIVE;
		else if (g[cur][num][j] == DEAD && count == 3)
			g[(cur + 1) & 1][num][j] = ALIVE;
		else
			g[(cur + 1) & 1][num][j] = DEAD;
	}
	pthread_exit(NULL);
}

int main(int argc, char *argv[])
{
	cur = 0;
	int i, j, n, m;

	printf("\033[2J"); /* clear screen */
	for (i = 0; i <= N + 1; ++i)
	{
		for (j = 0; j <= M + 1; ++j)
		{
			g[0][i][j] = random() & 1;
		}
		printf("Initializing g[%6d]...\r", i);
	}

	for (n = 0; n < NLOOP; n++)
	{
		printf("\033[2J"); /* clear screen */
		printf("n = %d\n", n);
		for (i = 1; i <= N >> 8; ++i)
		{
			for (j = 1; j <= M >> 7; ++j)
			{
				printf("%c", g[cur][i][j] == ALIVE ? '@' : '.');
			}
			printf("\n");
		}
		int num[N_THREADS];
		pthread_t th[N_THREADS];
		for (m = 1; m < N; m++)
		{
			num[m] = m;
			int r;
			r = pthread_create(&th[m], NULL, computeNextGen, (void *)&num[m]);
			if (r != 0)
				perror("new thread");
		}
		for (m = 1; m < N; m++)
		{
			pthread_join(th[m], NULL);
		}

		cur = (cur + 1) & 1;
	}

	return 0;
}