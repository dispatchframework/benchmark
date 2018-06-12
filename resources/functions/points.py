import random

class node:
    def __init__(self, value):
        self.value = value
        self.children = []
        self.power = random.randint(2, 5)
    def add_child(self, child):
        self.children.append(child)
    def __hash__(self):
        return hash(str(self.value**self.power)+str(len(self.children)))
    def stringify(self, prefix):
        result = ""
        result += "Node %s has %s children:"%(str(self.value), str(len(self.children)))
        for child in self.children:
            result += "\n%s%s"%(prefix, child.stringify(prefix+"\t"))
        return result
        
    def __str__(self):
        return self.stringify("\t")
                                             
def graph_gen(size):    
    root = node(random.randint(0, 2*size))
    prev_level = [root]
    this_level = []
    elements = []
    used = set([root.value])
    for i in range(size):        
        advance = random.randint(0, size//4)        
        value = random.randint(0, 2*size)
        while value in used:
            value = random.randint(0, 2*size)
        used.add(value)
        current = node(value)
        elements.append(current.value)
        ind = random.randint(0, len(prev_level)-1)
        prev_level[ind].add_child(current)
        if advance > size//8 or len(this_level) == 0:
            this_level.append(current)
        else:
            prev_level = this_level
            this_level = []
    return root, elements

def print_graph(root):
    queue = [root]
    visited = set()
    print(root)

def bfs(root, dest):
    queue = [root]
    visited = set()
    parents = {root.value: None}
    while queue:
        elem = queue.pop(0)
        if elem in visited:
            continue
        if elem.value == dest:
            return backtrack(parents, dest)
        else:
            for child in elem.children:                
                parents[child.value] = elem.value
                queue.append(child)
            visited.add(elem)
    return None
        
def backtrack(parents, dest):
    current = dest
    parent = parents[dest]
    path = [dest]
    while parent:
        current = parent
        parent = parents[current]
        path.append(current)
    return len(([str(x) for x in path]))

def handle(ctx, payload):
    size = 100000
    search = random.randint(0, size)
    if payload:
        search = payload.get("Search", search)
        size = payload.get("Size", size)
    graph, elements = graph_gen(size)            
    print("Looking for", search)
    result = bfs(graph, search)
    if result:
        return "Found path of length %d ending at %d" % (result, search)
    return "No path found"

print(handle(1, None))
