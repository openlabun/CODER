package cache_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/course"
	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	course_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/course"
	cache_infra "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	cache_course "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble/course"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_course_infra "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/course"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
	hasher "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
	test "github.com/openlabun/CODER/apps/api_v2/test"
)

func TestCacheCourseCRUD(t *testing.T) {
	process := test.StartTest(t, "Cache Course Creation and Persistence")
	email := "test@test.com"
	password := "Password123!"

	// [STEP 1] Initialize Roble + MariaDB clients, schema and repositories
	process.StartStep("Inicializando Roble, MariaDB y cache repositories")

	httpClient := &http.Client{Timeout: 15 * time.Second}
	robleClient, err := roble_infrastructure.NewRobleClient(httpClient)
	if err != nil {
		process.Fail("initialize roble client", err)
	}

	mariaClient, err := cache_infra.NewMariaDBClient()
	if err != nil {
		t.Skipf("cache DB unavailable, skipping: %v", err)
	}
	defer mariaClient.Close()

	cacheAdapter := cache_infra.NewCacheAdapter(mariaClient)
	initializer := cache_infra.NewSchemaInitializer(mariaClient, cache_infra.DefaultEntityConfigs())
	if err := initializer.Initialize(); err != nil {
		process.Fail("initialize cache schema", err)
	}

	robleAdapter := roble_infrastructure.NewRobleDatabaseAdapter(robleClient)
	userRepository := roble_user_infrastructure.NewUserRepository(robleAdapter)
	authAdapter := roble_user_infrastructure.NewRobleAuthAdapter(robleAdapter, userRepository)
	courseRepository := cache_course.NewCourseRepository(robleAdapter, cacheAdapter)
	robleDirectCourse := roble_course_infra.NewCourseRepository(robleAdapter)
	process.Log("Clientes y repositories inicializados")
	process.EndStep()

	// [STEP 2] Hash password
	process.StartStep("Hashear password")
	secAdapter := hasher.NewSecurityAdapter()
	hashedPassword, err := secAdapter.Hash(password)
	if err != nil {
		process.Fail("hash password", err)
	}
	process.EndStep()

	// [STEP 3] Login teacher
	process.StartStep("Login docente")
	access, err := authAdapter.LoginUser(email, hashedPassword)
	if err != nil {
		process.Fail("login user", err)
	}
	if access == nil || access.UserData == nil || access.UserData.ID == "" {
		process.Fail("login user", fmt.Errorf("expected logged user data with valid ID"))
	}
	if access.Token == nil || access.Token.AccessToken == "" {
		process.Fail("login user", fmt.Errorf("expected access token in login response"))
	}
	ctx := services.WithAccessToken(context.Background(), access.Token.AccessToken)
	teacherID := access.UserData.ID
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))
	process.EndStep()

	// [STEP 4] Create course (write-through: Roble + cache)
	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("CACHE-ENR-%d", now.UnixNano())

	process.StartStep("Crear curso via cache repository")
	course, err := course_factory.NewCourse(
		"Cache Integration Course",
		"Created by cache integration test",
		consts.CourseVisibilityPublic,
		consts.CourseColourBlue,
		fmt.Sprintf("CC-%d", now.Unix()%100000),
		&course_entities.Period{Year: now.Year(), Semester: consts.AcademicFirstPeriod},
		enrollmentCode,
		"https://example.test/enroll/"+enrollmentCode,
		teacherID,
	)
	if err != nil {
		process.Fail("build course with factory", err)
	}

	createdCourse, err := courseRepository.CreateCourse(ctx, course)
	if err != nil {
		process.Fail("create course", err)
	}
	if createdCourse == nil || createdCourse.ID == "" {
		process.Fail("create course", fmt.Errorf("expected created course with ID"))
	}
	courseID := createdCourse.ID
	process.Log(fmt.Sprintf("Curso creado. courseID=%s", courseID))

	robleCreated, err := robleDirectCourse.GetCourseByID(ctx, courseID)
	if err != nil || robleCreated == nil {
		process.Fail("roble verify create course", fmt.Errorf("expected course persisted in Roble, err=%v", err))
	}
	process.Log(fmt.Sprintf("Roble verify: course persisted id=%s name=%q", robleCreated.ID, robleCreated.Name))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando curso %s", courseID)
		_ = courseRepository.DeleteCourse(ctx, courseID)
	}()

	// [STEP 5] Read from cache (populated on create)
	process.StartStep("Leer curso por ID — cache hit tras creacion")
	cached, err := courseRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		process.Fail("get course by id", err)
	}
	if cached == nil {
		process.Fail("get course by id", fmt.Errorf("expected course from cache"))
	}
	if cached.ID != courseID || cached.ProfessorID != teacherID {
		process.Fail("cache data integrity", fmt.Errorf("cached fields mismatch: id=%s professorID=%s", cached.ID, cached.ProfessorID))
	}
	process.Log(fmt.Sprintf("Cache hit: name=%q professorID=%s", cached.Name, cached.ProfessorID))
	process.EndStep()

	// [STEP 6] Repeated read — consistency check
	process.StartStep("Segunda lectura consecutiva — consistencia de cache")
	cached2, err := courseRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		process.Fail("second get course by id", err)
	}
	if cached2 == nil || cached2.ID != cached.ID || cached2.Name != cached.Name {
		process.Fail("repeated cache hit", fmt.Errorf("expected identical data on repeated read"))
	}
	process.Log("Lecturas repetidas consistentes")
	process.EndStep()

	// [STEP 7] GetCoursesByTeacherID — list caching
	process.StartStep("Listar cursos del docente y verificar presencia en lista")
	teacherCourses, err := courseRepository.GetCoursesByTeacherID(ctx, teacherID)
	if err != nil {
		process.Fail("get courses by teacher id", err)
	}
	found := false
	for _, c := range teacherCourses {
		if c != nil && c.ID == courseID {
			found = true
			break
		}
	}
	if !found {
		process.Fail("course list", fmt.Errorf("expected created course %s in teacher list", courseID))
	}
	process.Log(fmt.Sprintf("Lista docente: %d cursos, curso presente", len(teacherCourses)))
	process.EndStep()

	// [STEP 8] GetCoursesByEnrollmentCode — FindBy cache
	process.StartStep("Buscar por enrollment code (FindBy cache)")
	byCode, err := courseRepository.GetCoursesByEnrollmentCode(ctx, enrollmentCode)
	if err != nil {
		process.Fail("get course by enrollment code", err)
	}
	if byCode == nil || byCode.ID != courseID {
		process.Fail("get course by enrollment code", fmt.Errorf("expected course with matching enrollment code"))
	}
	process.Log(fmt.Sprintf("Curso encontrado por enrollment code: id=%s", byCode.ID))
	process.EndStep()

	// [STEP 9] Update course (write-through update)
	process.StartStep("Actualizar curso — write-through update")
	createdCourse.Name = "Cache Integration Course Updated"
	createdCourse.Description = "Updated by cache integration test"
	updatedCourse, err := courseRepository.UpdateCourse(ctx, createdCourse)
	if err != nil {
		process.Fail("update course", err)
	}
	if updatedCourse == nil || updatedCourse.Name != "Cache Integration Course Updated" {
		process.Fail("update course", fmt.Errorf("expected updated course name"))
	}
	process.Log(fmt.Sprintf("Update aplicado. Nombre=%q", updatedCourse.Name))

	robleUpdated, err := robleDirectCourse.GetCourseByID(ctx, courseID)
	if err != nil || robleUpdated == nil || robleUpdated.Name != "Cache Integration Course Updated" {
		process.Fail("roble verify update course", fmt.Errorf("expected updated course in Roble"))
	}
	process.Log(fmt.Sprintf("Roble verify: updated name=%q", robleUpdated.Name))
	process.EndStep()

	// [STEP 10] Read after update — cache must reflect new value
	process.StartStep("Leer tras actualizacion — cache debe reflejar cambio")
	afterUpdate, err := courseRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		process.Fail("get course after update", err)
	}
	if afterUpdate == nil || afterUpdate.Name != "Cache Integration Course Updated" {
		process.Fail("cache after update", fmt.Errorf("expected updated name in cache, got %q", afterUpdate.Name))
	}
	process.Log(fmt.Sprintf("Cache actualizada. Nombre=%q", afterUpdate.Name))
	process.EndStep()

	// [STEP 11] Delete course — cache invalidation
	process.StartStep("Eliminar curso y verificar invalidacion de cache")
	if err := courseRepository.DeleteCourse(ctx, courseID); err != nil {
		process.Fail("delete course", err)
	}
	deleted, err := courseRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		process.Fail("get course after delete", err)
	}
	if deleted != nil {
		process.Fail("cache invalidation", fmt.Errorf("expected course %s to be deleted from cache and Roble", courseID))
	}
	robleDeleted, err := robleDirectCourse.GetCourseByID(ctx, courseID)
	if err != nil {
		process.Fail("roble verify delete course", err)
	}
	if robleDeleted != nil {
		process.Fail("roble verify delete course", fmt.Errorf("expected course %s deleted from Roble", courseID))
	}
	process.Log("Roble verify: course deleted from Roble")
	process.Log("Curso eliminado — cache y Roble limpios")
	process.EndStep()

	process.End()
}
