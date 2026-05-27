# Plan de Pruebas

### Pruebas de Seguridad

- **Pruebas de Acceso de Usuario:** 

<img width="838" height="373" alt="image" src="https://github.com/user-attachments/assets/6261bbe3-62d1-4081-bf28-58dd2934fc1d" />


```
go test -v ./test/roble_auth -run TestUserLogin
```

- **Pruebas de Registro:** 

<img width="844" height="306" alt="image" src="https://github.com/user-attachments/assets/8c38908c-1296-42ea-9198-fc1b193912ce" />

```
go test -v ./test/roble_auth -run TestUserRegistration
```

### Pruebas de Persistencia

- **Persistencia de Cursos:** 

<img width="1275" height="616" alt="image" src="https://github.com/user-attachments/assets/3e555e75-d5fd-4ff2-a1cd-86fb535df8e1" />
<img width="905" height="437" alt="image" src="https://github.com/user-attachments/assets/fa818258-d6a4-40ea-b2d3-14fd33ae4019" />

```
go test -v ./test/roble_persistence -run TestCourseCRUD
```

- **Persistencia de Exámenes:** 

<img width="1337" height="796" alt="image" src="https://github.com/user-attachments/assets/47f5b0bb-e54e-4ea3-961a-b06a6db68827" />
<img width="1327" height="703" alt="image" src="https://github.com/user-attachments/assets/df14cd2f-04ba-4e6f-8908-0566d97599a7" />
<img width="983" height="772" alt="image" src="https://github.com/user-attachments/assets/ed5537e0-7632-4f81-ba27-965404f4cdac" />
<img width="880" height="411" alt="image" src="https://github.com/user-attachments/assets/76fd934f-5b15-44c3-b58b-5b0e7619f3f5" />

```
go test -v ./test/roble_persistence -run TestExamCRUD
```

- **Persistencia de Revisiones:**

<img width="1475" height="729" alt="image" src="https://github.com/user-attachments/assets/6d214a9b-2d35-4841-878a-ebae3d3e9083" />
<img width="1306" height="662" alt="image" src="https://github.com/user-attachments/assets/6803ca44-c6fc-4c4f-a3c4-d27f0d4f1662" />
<img width="826" height="752" alt="image" src="https://github.com/user-attachments/assets/0f8cc5d6-d9c5-4077-9e87-c9fbecacfbce" />


```
go test -v ./test/roble_persistence -run TestSubmissionCRUD
```

### Pruebas Funcionales

##### Módulo de Autenticación

- **Autenticación de Usuario Estudiante:**

<img width="1476" height="465" alt="image" src="https://github.com/user-attachments/assets/3049fe95-37f5-4dcf-b4cb-c1bf16fb5a12" />

```
go test -v ./test/use_cases/user -run TestStudentAuth
```

- **Autenticación de Usuario Docente:**
  
  <img width="1499" height="572" alt="image" src="https://github.com/user-attachments/assets/21e72dab-a3f8-46db-9967-822a1ee1bdf2" />

```
go test -v ./test/use_cases/user-run TestTeacherAuth
```

##### Módulo de Cursos

- **CRUD de un Curso:**



```
go test -v ./test/use_cases/course -run TestCourseCRUD
```

- **Inscripción de Estudiante:**



```
go test -v ./test/use_cases/course -run TestCourseEnrollment
```

- **Curso desde la vista de Estudiante:**



```
go test -v ./test/use_cases/course -run TestCourseFromStudentView
```

##### Módulo de Exámenes

- **CRUD de un Examen:**

<img width="1132" height="546" alt="image" src="https://github.com/user-attachments/assets/c6634958-a5f1-4576-90b3-6ead6268d635" />
<img width="903" height="439" alt="image" src="https://github.com/user-attachments/assets/bb3c4430-18ac-450b-b47e-6daf928cf23d" />


```
go test -v ./test/use_cases/exam -run TestExamCRUD
```

- **Examen (`public`) desde la vista de Estudiante:**

<img width="985" height="707" alt="image" src="https://github.com/user-attachments/assets/76ff7f21-c240-4c09-ad30-bbdebf3bb8e3" />
<img width="1010" height="642" alt="image" src="https://github.com/user-attachments/assets/cfecf82b-124d-4af6-8035-a07af4dd0cac" />

```
go test -v ./test/use_cases/exam -run TestExamFromStudentView
```

- **Examen (`teachers`) desde la vista de Docente:**

<img width="1000" height="547" alt="image" src="https://github.com/user-attachments/assets/f91033aa-e2da-470c-8b02-bac2f4452575" />
<img width="1157" height="596" alt="image" src="https://github.com/user-attachments/assets/13f8b639-13cc-4829-95f4-ce3cae3c2dbb" />


```
go test -v ./test/use_cases/exam -run TestExamFromTeacherView
```

- **CRUD de un Reto (`Challenge`):**

<img width="1045" height="699" alt="image" src="https://github.com/user-attachments/assets/d2ff6ff8-ca01-4044-af69-93b2a0254772" />


```
go test -v ./test/use_cases/exam -run TestChallengeCRUD
```

- **Estados de un Reto:**

<img width="897" height="740" alt="image" src="https://github.com/user-attachments/assets/26a0207c-246e-480d-a71e-a0839c4d1a37" />
<img width="931" height="788" alt="image" src="https://github.com/user-attachments/assets/be8056c2-33c0-4338-86fc-403e89c5318a" />

```
go test -v ./test/use_cases/exam -run TestChallengeStates
```

- **Fork de un Reto:**

<img width="987" height="504" alt="image" src="https://github.com/user-attachments/assets/7ef3e53b-8db2-47a2-8239-0461b96a78e7" />
<img width="995" height="505" alt="image" src="https://github.com/user-attachments/assets/3fc4e3cd-6607-41e8-9b65-5f245bb7d8f8" />

```
go test -v ./test/use_cases/exam -run TestChallengeFork
```

- **Visibilidad de un Reto para Docentes:**

<img width="886" height="444" alt="image" src="https://github.com/user-attachments/assets/d283153f-dfd0-42b2-b969-a6f11a6765aa" />
<img width="986" height="438" alt="image" src="https://github.com/user-attachments/assets/9ecf6ece-c3cb-4608-87e3-ee1128874d20" />

```
go test -v ./test/use_cases/exam -run TestChallengeFromTeachersView
```

- **CRUD de Casos de uso (`TestCase`):**

<img width="904" height="772" alt="image" src="https://github.com/user-attachments/assets/dadd9b29-b626-4ab9-a0af-733998cc8029" />

```
go test -v ./test/use_cases/exam -run TestTestCaseCRUD
```

- **Casos de Uso desde la vista de Estudiante:**
- 
<img width="830" height="712" alt="image" src="https://github.com/user-attachments/assets/8d2e324e-df9a-4529-ba48-ec0e93221be5" />


```
go test -v ./test/use_cases/exam -run TestTestCaseFromStudentView
```

- **CRUD de puntos de examen (`ExamItem`):**

<img width="866" height="465" alt="image" src="https://github.com/user-attachments/assets/6ccdb156-4eea-4e43-9b66-4d4de6699c0d" />
<img width="896" height="506" alt="image" src="https://github.com/user-attachments/assets/0b709229-6f5c-4521-a746-c4f89af3b36a" />

```
go test -v ./test/use_cases/exam -run TestExamItemCRUD
```

- **Privacidad de Retos para Puntos de Examen:**

<img width="864" height="766" alt="image" src="https://github.com/user-attachments/assets/704847e4-5dfe-4d61-8235-89322ff8ef1a" />

```
go test -v ./test/use_cases/exam -run 
```

- **Fork automático para Punto de Examen:**

<img width="1085" height="788" alt="image" src="https://github.com/user-attachments/assets/f2b3a9cf-78f6-4a11-bd8f-cb53454dd30c" />

```
go test -v ./test/use_cases/exam -run TestExamItemChallengeFork
```

- **Privacidad de Puntos de Examen:**

<img width="955" height="464" alt="image" src="https://github.com/user-attachments/assets/0cba4315-e6f8-412a-81f1-416a8cdf433e" />
<img width="1085" height="529" alt="image" src="https://github.com/user-attachments/assets/dc200c20-3a08-4ea9-82ac-6a4bac0f0c69" />

```
go test -v ./test/use_cases/exam -run TestExamItemPrivacy
```

##### Módulo de Revisiones

- **CRUD de Sesiones (`Session`):**

<img width="853" height="484" alt="image" src="https://github.com/user-attachments/assets/2d1dceba-2704-4c64-8e28-c4d73a07d4a4" />
<img width="884" height="508" alt="image" src="https://github.com/user-attachments/assets/0e25478d-16ed-434d-9414-e50de146be01" />

```
go test -v ./test/use_cases/submission -run TestSessionCRUD
```

- **Heartbeat de Sesiones:**

<img width="897" height="771" alt="image" src="https://github.com/user-attachments/assets/7987d0ef-ed60-4ced-95f9-419922349ca0" />

```
go test -v ./test/use_cases/submission -run TestSessionHeartbeat
```

- **Congelamiento y Bloqueo de Sesiones:**

<img width="1004" height="683" alt="image" src="https://github.com/user-attachments/assets/34f8d74a-af63-4304-ac1a-90ff7f330867" />
<img width="939" height="661" alt="image" src="https://github.com/user-attachments/assets/0987e9cc-d767-49b0-9f2a-c211c9ed83e1" />

```
go test -v ./test/use_cases/submission -run TestSessionFreezeAndBlock
```

- **Creación y obtención de Revisiones (`Submission`):**

<img width="902" height="578" alt="image" src="https://github.com/user-attachments/assets/7772d238-495e-4db3-8551-f7e45620889a" />
<img width="929" height="703" alt="image" src="https://github.com/user-attachments/assets/dc8a703f-4c9f-4d1f-888f-2befeaa55000" />

```
go test -v ./test/use_cases/submission -run TestSubmissionCreateAndRead
```

- **Límite de intentos para examen:**

<img width="894" height="817" alt="image" src="https://github.com/user-attachments/assets/03db95d8-8de1-4eb1-8c38-8589026870a4" />

```
go test -v ./test/use_cases/submission -run TestExamTryLimit
```

- **Límite de intentos para revisiones:**

<img width="860" height="495" alt="image" src="https://github.com/user-attachments/assets/496c1e4e-34af-41ea-af55-63b6123409d1" />
<img width="962" height="642" alt="image" src="https://github.com/user-attachments/assets/7a1c1686-3722-4e9a-aba4-2f016c836d69" />

```
go test -v ./test/use_cases/submission -run TestSubmissionTryLimit
```

- **Revisiones Inválidas:**

<img width="1073" height="837" alt="image" src="https://github.com/user-attachments/assets/69c82efc-e974-479f-9f57-76f0aab98f75" />
<img width="1057" height="842" alt="image" src="https://github.com/user-attachments/assets/02186ce9-4203-45a2-96e3-73da0ffad706" />
<img width="970" height="312" alt="image" src="https://github.com/user-attachments/assets/ec6f93ff-dabf-409e-90c5-e216009db2a9" />

```
go test -v ./test/use_cases/submission -run TestInvalidSubmissions
```

- **Ejecución de Revisiones:**

<img width="857" height="440" alt="image" src="https://github.com/user-attachments/assets/49d2259e-1337-45c9-b377-6a376190337a" />
<img width="963" height="572" alt="image" src="https://github.com/user-attachments/assets/6c046471-b1d3-4d17-8cf7-e3fee355a8d6" />

```
go test -v ./test/use_cases/submission -run TestSubmissionExecution
```

- **Ejecución de Código sin revisión:**

<img width="979" height="486" alt="image" src="https://github.com/user-attachments/assets/f2cc670c-a9d9-457a-aba6-340d4e676d8f" />
<img width="997" height="640" alt="image" src="https://github.com/user-attachments/assets/ac51150c-1f45-423a-acee-480e236b8232" />

```
go test -v ./test/use_cases/submission -run TestSubmissionWithoutScore
```

- **Ejecución de Código con caso personalizado:**

<img width="1018" height="490" alt="image" src="https://github.com/user-attachments/assets/e22c5a89-c89f-4405-8a52-7bfae81100ba" />
<img width="980" height="633" alt="image" src="https://github.com/user-attachments/assets/5272c130-054b-4f22-8b5e-2a3da07a80be" />

```
go test -v ./test/use_cases/submission -run 
```

### Pruebas de Rendimiento

- **Inicios de sesión simultáneos:**

<img width="902" height="487" alt="image" src="https://github.com/user-attachments/assets/a505d85e-338b-4b65-807a-c40709cca5c8" />

```
go test -v ./test/use_cases/performance -run TestLoginPerformance
```

- **Lista de exámenes (`Exam`) públicos:**

<img width="901" height="330" alt="image" src="https://github.com/user-attachments/assets/eefd6799-030f-4db7-ad0d-75f5220993ae" />


```
go test -v ./test/use_cases/performance -run TestExamListPerformance
```

- **Activación y Heartbeat de sesiones (`Session`):**

<img width="901" height="615" alt="image" src="https://github.com/user-attachments/assets/fbc42464-be95-4494-ad8b-70b650209ca8" />

```
go test -v ./test/use_cases/performance -run TestSessionsPerformance
```

- **Revisiones (`Submission`) simultáneas:**

<img width="915" height="679" alt="image" src="https://github.com/user-attachments/assets/cd1f1b7b-6078-4bc7-a674-807c845525be" />

```
go test -v ./test/use_cases/performance -run TestSubmissionPerformance
```
