#include <stdio.h>
#include <pthread.h>
#include <stdlib.h>
static long num_rects = 2 * 1000 * 1000 * 1000;
int main(int argc, char *argv[])
{
    int i;
    double mid, height, width, sum = 0.0;
    double area;
    width = 1.0 / (double)num_rects;
    for (i = 0; i < num_rects; ++i)
    {
        mid = (i + 0.5) * width;
        height = 4.0 / (1.0 + mid * mid);
        sum += height;
    }
    area = width * sum;
    printf("Computed pi = %f\n", area);
    return 0;
}