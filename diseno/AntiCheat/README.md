## Anti-cheat de código

Para detectar plagio entre envíos y uso de código generado por IA, se utilizarán las siguientes técnicas:

### Detección de plagio entre estudiantes

- **Tokenización y n-grams**: El código se tokeniza (keywords, operadores, identificadores) y se generan ventanas de n-grams (secuencias de ~5-10 tkens). Se compara el porcentaje de n-grams compartidos entre pares de envíos del mismo reto. Si supera un umbral configurable, se marca para revisión humana. Se inspira en el enfoque de Moss/JPlag.
- **Normalización previa**: Antes de tokenizar, se normaliza el código (quitar comentarios, formateo, nombres de variables genéricos) para reducir falsos positivos por cambios cosméticos.
- **Comparación estructural (AST)**: Se parsea el código a árbl de sintaxis abstracta (AST) y se compara la estructura entre envíos. Permite detectar copias con variables o nombres distintos. Se usará tree-sitter o parsers por lenguaje (Python: ast, JS: acorn, etc.).

### Detección de código generado por IA

- **Clasificador humano vs IA**: Se utiliza un modelo preentrenado (p. ej. Hugging Face o similar) que clasifica el código como humano o generado por IA según features como longitud, patrones típicos y estadísticas del texto. Los envíos marcados como sospechosos pasan a revisión humana, no a sanción automática.
- **Umbrales configurables**: El prfesor puede definir umbrales de similitud (plagio) y de confianza (IA) para decidir cuándo generar alertas.

### Alcance práctico del proyecto

La implementación priorizará tokenización + n-grams y normalización para plagio, y un clasificador existente para IA. La comparación por AST se añadirá en una segunda fase si el tiemp lo permite.
