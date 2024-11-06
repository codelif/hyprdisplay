import random


def get_intersections(A, B):
    get_sorted_by_element(A, 1)
    get_sorted_by_element(A, 0)
    get_sorted_by_element(B, 1)
    get_sorted_by_element(B, 0)
    res = []
    a, b = 0, 0
    if A and B:
        while a < len(A) and b < len(B):
            lo = max(A[a][0], B[b][0])
            hi = min(A[a][1], B[b][1])
            if lo <= hi:
                res.append((lo, hi))

            if A[a][1] < B[b][1]:
                a += 1
            else:
                b += 1
    return res


def get_random_rects(n):
    random.seed(1)
    rects = []
    for i in range(n):
        a = c = random.randint(31, 60)
        b = random.randint(70, 90)
        d = random.randint(70, 110)
        rects.append((a, b, c, d))
    return rects


def get_sorted_by_element(A, i):
    return sorted(A, key=lambda x: x[i])


rects_a = get_random_rects(1)
rects_b = get_random_rects(20)[1:]
x_axis = get_intersections(
    [(a, b) for a, b, _, _ in rects_a], [(a, b) for a, b, _, _ in rects_b]
)
y_axis = get_intersections(
    [(c, d) for _, _, c, d in rects_a], [(c, d) for _, _, c, d in rects_b]
)


print(f"rectangle to compare: {rects_a}")
print(f"other rectangles: {rects_b}")
print(f"x_axis intersections: {x_axis}")
print(f"y_axis intersections: {y_axis}")

x_axis = get_sorted_by_element(x_axis, 0)
y_axis = get_sorted_by_element(y_axis, 1)
print(f"sorted x_axis intersections: {x_axis}")
print(f"sorted y_axis intersections: {y_axis}")

print(x_axis, y_axis)
