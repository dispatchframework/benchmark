import random
def handle(ctx, payload):
    list1 = list()
    for i in range(random.randint(4, 15)):
        list1.append(i)
    s = "-"
    s = s.join([str(x) for x  in list1]) 
    return {"result": s}
