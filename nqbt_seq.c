/*
 * 並列分散処理 第5回課題
 * Nov. 30th, 2021
 */
/*
Nクイーン問題
クイーンとは、マスの縦横斜めに位置する別のコマを取れる
この問題はどのクイーンにも取られないクイーンの配置を探し、その位置をN*Nのマスの中にいくらおけるかというパターン数を求める問題である

このコードはそれを実現している
　horizontalの部分
  縦に対してクイーンが存在しないか
  obliquely upwardの部分
  右上から左下のかけての部分でクイーンが存在しないか
  diagonally downwardの部分
  左上から右下にかけての部分でクイーンが存在しないか
  各ポイントに１を置く→その上で別の、条件を満たす場所に１を置く→その上で別の、条件を満たす場所に１を置く→１...を繰り返してパターン数を数えている。
  結果的に条件を満たす1を全て置けたパターンが見つかったら１をreturnしている。
  各ポイントにクイーンを置いてみて、その場合にパターンが存在するかをまた点を置いて調べているので、少なくとも入力した数値の４乗以上は実行を繰り返している。
  また、以下の並列処理をこなった結果、明らかに実行速度が速くなったことがわかった。
*/

#include <stdio.h>
#include <stdlib.h>

#define MAX_N (29)

int check_and_set(unsigned int mat[], int n, int row, int col)
{
	int i, j, c = 0;
	//printf("  %d %d %d\n", n, row, col);

	/* horizontal */
	for (j = col - 1; j >= 0; j--) {
		if (mat[j] & (1U << row)) {
			//printf(" re1 %d %d %d\n",row, col, mat[j]);
			return c;
		}
	}
	//col:上下
	//row:左右

	//右上から左下の線
	/* obliquely upward */
	for (i = row - 1, j = col - 1; j >= 0 && i >= 0; j--, i--) {
		if (mat[j] & (1U << i)) {
			//printf(" re2 %d %d %d\n",i, col, mat[j]);
			return c;
		}
	}

	//左上から右下の線
	/* diagonally downward */
	for (i = row + 1, j = col - 1; j >= 0 && i < n; j--, i++) {
		if (mat[j] & (1U << i)) {
			//printf(" re3 %d %d %d\n",i, col, mat[j]);
			return c;
		}
	}

	/* set */
	mat[col] = 1 << row;
	//printf("col %d %d\n", col, 1<<row);


	if (col == n - 1) {
		/* completed */
		c++;
	} else {
		/* set remain columns */
		for (i = n - 1; i >= 0; i--)
			c += check_and_set(mat, n, i, col + 1);
	}
	mat[col] = 0U;

	return c;
}

int main(int argc, char *argv[])
{
	unsigned int mat[MAX_N], n;
	int i, count = 0;

	/* check argument */
	if (argc < 2) {
		fprintf(stderr, "Usage: %s number\n", argv[0]);
		return -1;
	}

	/* obtain n */
	n = atoi(argv[1]);
	if ((n < 2) || (n > MAX_N)) {
		fprintf(stderr, "You should specify a number "
			"between 2 and %d.\n", MAX_N);
		return -1;
	}

	/* initialize */
	for (i = 0; i < n; i++)
		mat[i] = 0U;

	/* query */
	//このfor文を並列化
	#pragma omp parallel for private(mat) reduction(+:count)
	for (i = n - 1; i >= 0; i--){
		count += check_and_set(mat,n, i, 0);
  }
	printf("Total answer for %d x %d = %d\n", n, n, count);

	return 0;
}
