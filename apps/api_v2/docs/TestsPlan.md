# Plan de Pruebas

### Pruebas de Seguridad

Actualmente el sistema soporta la autenticación de usuario a través de la plataforma institucional Roble [[1]](https://roble.openlab.uninorte.edu.co/docs). Existen pruebas directas de acceso a través de este servicio.

- **Pruebas de Acceso de Usuario:** Accede al usuario `test@test.com` creado por defecto anteriormente con un hash de la contraseña `Password123!`

```
go test -v ./test/roble_auth -run TestUserLogin
```

- **Pruebas de Registro:** Crea un usuario con el correo `test{random_value}@test.com` y contraseña `Password123!`

```
go test -v ./test/roble_auth -run TestUserRegistration
```

### Pruebas de Persistencia

El sistema de almacenamiento de datos, de igualmanera está soportado por la plataforma institucional Roble. En todos se usa el acceso del usuario `test@test.com`, creado con **permisos de docente**. Se realizó pruebas del **CRUD completo de datos** para cada uno de los módulos de la aplicación:

- **Persistencia de Cursos:** Se crea, actualiza, verifica y se elimina un curso (`Course`).

```
go test -v ./test/roble_persistence -run TestCourseCRUD
```

- **Persistencia de Exámenes:** Inicialmente se crea un curso (`Course`) al que vincular un examen (`Exam`), para posteriormente hacer un reto (`Challenge`) con sus respectivos casos de prueba (`TestCase`), el cual se vinculará al examen a través de un punto (`ExamItem`). En cada caso se hará una prueba de CRUD y verificación que el reto esté vinculado al examen. Para los casos de retos y casos de prueba se deberán crear variables (`IOVariable`).

```
go test -v ./test/roble_persistence -run TestExamCRUD
```

- **Persistencia de Revisiones:** En este proceso se crean un curso (`Course`), examen (`Exam`) y reto (`Challenge`) auxiliares. Posteriormente se procede a crear una sesión (`Session`), se actualizará y verificará su estado y valores. A través de la sesión se creará una revisión (`Submission`) y finalmente los resultados de la revisión (`SubmissionResults`), en ambos casos de validará todo el proceso de CRUD.

```
go test -v ./test/roble_persistence -run TestSubmissionCRUD
```

### Pruebas Funcionales



### Pruebas de Rendimiento
