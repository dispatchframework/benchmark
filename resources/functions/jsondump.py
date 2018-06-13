import random
import json
import string

def seed_generic(generic):
    for key in generic:
        value = ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(100))
        generic[key] = value

def benchmark(data):
    for i in range(20):
        jsonified = json.dumps(data)
        back = json.loads(jsonified)
        seed_generic(data)


def handle(ctx, payload):
    length = 1000
    if payload:
        length = payload.get("Length", length)
    keys = [random.randint(10**6, 4*10**6) for _ in range(length)]
    generic = {key: random.randint(0, 100) for key in keys}
    seed_generic(generic)
    benchmark(generic)
    
