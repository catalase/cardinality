#/usr/bin/env python3
# -*- coding: utf8 -*-
import random

def solve_quadratic(a, b, c, *, real=True):
    D = b ** 2 - 4 * a * c
    if real:
        if D < 0:
            raise ValueError("quadratic equation does not have real root")

    # 필요한 경우, 제곱근의 계산은 자동으로 복소수 영역에서 이루어진다.
    rD = D ** 0.5

    return (-b + rD) / (2 * a), (-b - rD) / (2 * a)

# Z가 표준정규분포를 따를 때 P(|Z| <= k) * 100 % 신뢰구간
# m 표본 크기
# R 중복된 갯수
# L 알고 있는 개수
def guess(k, m, R, L):
    return (max(solve_quadratic(R, (m - k) * -L, -k * L ** 2)),
            max(solve_quadratic(R, (m + k) * -L, k * L ** 2)))

def test():
    # 실제 갯수
    n = 50000

    # 알고 있는 개수
    L = 100

    # 중복된 갯수
    R = 0

    # 표본 크기
    m = 40000

    # 실험 결과
    # L 도 적당히 커야하지만, m도 그에 못지않게 "많이" 커야한다.
    # 일반적으로 L, m이 둘다 커야 한다.
    #
    # 위에서 언급한 "많이"가 어느정도인지는 나도 모른다. 그렇다고 너무 적자니
    # 추정이 널뛰기하고, 그러지 않자니 정확하게 아는 바가 없고, 하여튼 표본은 커야 한다.
    choose = lambda: random.randint(1, n)
    database = range(1, L + 1)

    for i in range(m):
        num = choose()
        if num in database:
            R += 1

    if R != 0:
        print("중복된 갯수", R)
        print(guess(0.00, m, R, L)[0])
        print(guess(1.96, m, R, L))
        print(guess(2.58, m, R, L))
        print(guess(0.23, m, R, L))
    else:
        print("안타깝게도 중복이 없었음")

# 아래를 실행시켜본다면 신기한 결과를 얻을 수 있다.
# guess(2.58, 1000, 900, 1000)

def main():
    test()

if __name__ == '__main__':
    main()
