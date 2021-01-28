import ast


class HyperTransformer(ast.NodeTransformer):
    def __init__(self, hyperparamters):
        self.hyperparamters = hyperparamters

    def isInHP(self, nodeid):
        for hp in self.hyperparamters:
            if hp.name == nodeid:
                return True, hp.value
        return False, 0

    def visit_Assign(self, node):
        if isinstance(node, ast.Assign) and len(node.targets) == 1 and isinstance(node.targets[0], ast.Name):
            found, val = self.isInHP(node.targets[0].id)
            if found:
                if isinstance(node.value, ast.Num):
                    new = ast.Num(
                        n=float(val)
                    )
                elif isinstance(node.value, ast.Str):
                    new = ast.Str(
                        s=val
                    )
                # TODO(gaocegege): Support more cases.
                else:
                    return node
                return ast.copy_location(ast.Assign(
                    targets=[
                        ast.Name(id=node.targets[0].id,
                                 ctx=node.targets[0].ctx),
                    ],
                    value=new,
                ), node)
            else:
                return node
        else:
            return node
