import sys
from collections import defaultdict

class SparseMatrixSolver:
    def __init__(self):
        # 1. Estructura principal: Fila -> {Columna: Valor}
        # Cumple con el requisito de "nested dictionary structure"
        self.matrix = defaultdict(dict)
        
        # 2. Estructura auxiliar para consultas rápidas: Columna -> Suma Total
        # Esto permite responder queries en O(1) sin recorrer las filas
        self.col_sums = defaultdict(int)

    def update(self, row, col, value):
        """
        Guarda el valor y actualiza la suma de la columna correspondiente.
        Maneja correctamente si ya existía un valor en esa celda (sobreescritura).
        """
        # Obtenemos el valor anterior (0 si no existía)
        previous_value = self.matrix[row].get(col, 0)
        
        # Guardamos el nuevo valor en la estructura principal
        self.matrix[row][col] = value
        
        # Actualizamos la suma de la columna restando el viejo y sumando el nuevo
        # Esto es más seguro que solo sumar, por si hay correcciones de datos
        self.col_sums[col] += (value - previous_value)

    def query(self, col):
        """
        Retorna la suma total de la columna solicitada.
        """
        return self.col_sums[col]

# Ejemplo de uso basado en tu descripción
def main():
    solver = SparseMatrixSolver()

    # Supongamos que esta es la entrada de datos (R, C, Valor)
    # (En un caso real, esto vendría de input() o un archivo)
    inputs = [
        (1, 5000000, 10),  # Fila 1, Columna 5M, Valor 10
        (2, 5000000, 20),  # Fila 2, Columna 5M, Valor 20
        (100, 42, 5),      # Fila 100, Columna 42, Valor 5
        (1, 5000000, 15)   # Actualización: Fila 1, Columna 5M cambia de 10 a 15
    ]

    print("Procesando entradas...")
    for r, c, v in inputs:
        solver.update(r, c, v)

    # Consultas (Q queries)
    queries = [5000000, 42, 999] # 999 no existe, debería dar 0

    print("\nResultados de las consultas:")
    for col_index in queries:
        total = solver.query(col_index)
        print(f"Suma de la columna {col_index}: {total}")

if __name__ == "__main__":
    main()