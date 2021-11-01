package main

import (
	"fmt"
)
//Sieve of Eratosthenes
//エラトステネスのふるい
//素数xを定めて、ある数yを用い、x*yが上限値prime至るまでyを大きくしていく。
//→このx*yで求まる数は素数ではない
//x*yが上限に至ったら、xを次の素数にする
//このxは√primeになるまで引き上げる
//これを順々に実行していき、√primeにxが至ったら終了
//これを配列の真偽値で実装して524297の要素が素数であるかを確認
//加えて下から数え上げをすれば素数が何番目なのかがわかる

func main() {
	const prime = 524297
	var primenum [prime + 1]bool
	var primeindex = 0
	for x := 2; x*x <= prime; x++ {
		if !primenum[x] {
			for y := 2; x*y <= prime; y++ {
				primenum[x*y] = true
			}
		}
	}
	for n := 2; n <= prime; n++ {
		if !primenum[n] {
			primeindex++
		}
		if n == prime {
			if primenum[n] {
				fmt.Println(prime,"is not Prime Number")
			} else {
				fmt.Println(prime,"is Prime Number(Number", primeindex, ")")
			}
		}
	}
}
